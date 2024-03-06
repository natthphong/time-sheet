package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/config"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/eft"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/internal/cache"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/internal/db"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/internal/httputil"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/internal/kafka"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/internal/logz"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/internal/scramkafka"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/internal/secret"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/internal/sftp"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/job"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {

	lambda.Start(LambdaHandler)
}

func LambdaHandler() {
	ctx := context.Background()

	config.InitTimeZone()
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Unable to initial config."))
	}

	logz.Init(cfg.Log.Level, cfg.Server.Name)
	//defer logz.Drop()

	ctx, cancel = context.WithCancel(ctx)
	defer cancel()
	logger := zap.L()

	err = secret.ConfigCommonSecret(cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Unable to initial common secret"))
	}
	err = secret.ConfigRDSSecret(cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Unable to initial rds secret"))
	}

	jsonCfg, err := json.Marshal(cfg)
	fmt.Println("after cfg : ", string(jsonCfg))
	dbPool, err := db.Open(ctx, cfg.DBConfig)
	if err != nil {
		logger.Fatal("server connect to db", zap.Error(err))
	}
	defer dbPool.Close()

	redisClient, err := cache.Initialize(ctx, cfg.RedisConfig)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error cannot connect redis."))
	}

	defer redisClient.Close()
	redisCmd := redisClient.CMD()

	var httpClient *http.Client
	logger.Debug("cert", zap.Any("", string(cfg.HTTP.CertFile)))
	if strings.Contains(cfg.Env, "UAT") {
		httpClient, err = httputil.InitHttpClientWithCertAndKey(
			cfg.HTTP.TimeOut,
			cfg.HTTP.MaxIdleConn,
			cfg.HTTP.MaxIdleConnPerHost,
			cfg.HTTP.MaxConnPerHost,
			cfg.HTTP.CertFile,
			cfg.HTTP.KeyFile,
		)
	} else {
		httpClient, err = httputil.InitHttpClientWithCert(
			cfg.HTTP.TimeOut,
			cfg.HTTP.MaxIdleConn,
			cfg.HTTP.MaxIdleConnPerHost,
			cfg.HTTP.MaxConnPerHost,
			cfg.HTTP.CertFile,
		)
	}

	if err != nil {
		log.Fatal(errors.Wrap(err, "Error in init httpClient."))
	}

	_ = httpClient

	_ = redisCmd
	internalProducer, err := scramkafka.NewSyncProducer(cfg.Kafka.Internal)
	if err != nil {
		logger.Fatal("Fail Create NewSyncProducer", zap.Error(err))
	}
	defer func() {
		if err = internalProducer.Close(); err != nil {
			logger.Fatal("Fail Close SyncProducer", zap.Error(err))
		}
	}()
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	})
	if err != nil {
		log.Fatal(errors.Wrap(err, "Unable to initial config."))
		return
	}
	svc := s3.New(sess)
	_ = svc
	sftpConfig := sftp.Config{
		Username: cfg.SFTPConfig.Username,
		Password: cfg.SFTPConfig.Password,
		Server:   cfg.SFTPConfig.Server,
		Timeout:  time.Second * 30,
	}

	sftpClient, err := sftp.New(sftpConfig)
	if err != nil {
		logger.Fatal("Error on sftp Connection", zap.Error(err))
	}
	defer sftpClient.Close()

	err = job.GetFileInsertToTblTemp(
		*cfg,
		ctx,
		logger,
		svc,
		job.DownLoadFileAndPushToS3(
			sftpClient,
			svc,
			job.PutFileToS3(cfg.S3Config.BucketName),
		),
		job.InsertKBANKDailyTemp(dbPool),
	)
	if err != nil {
		logger.Error("GetFileInsertToTblTemp", zap.Any("err ", err))
	}

	err = job.StageCheckFunc(
		ctx,
		logger,
		job.InsertUnMatedHeader(dbPool),
		job.GetListResult(dbPool),
		job.InsertUnMatedDetail(
			cfg.Producer.FinalTxnListener,
			cfg.Producer.GoldListener,
			cfg.FundTransferConfig,
			cfg.Exception,
			cache.GetRedis(redisCmd),
			eft.HTTPOauthFundTransferHttp(
				httpClient,
				cfg.FundTransferConfig.OauthFundTransferUrl,
				cfg.Toggle.OauthFundTransfer,
				cfg.FundTransferConfig.OauthRetry,
			),
			eft.HTTPInquiryStatusFundTransfer(
				httpClient,
				cfg.FundTransferConfig.InquiryStatusFundTransferUrl,
				cfg.Toggle.InquiryStatusFundTransfer,
				cfg.FundTransferConfig.InquiryStatusRetry,
			),
			dbPool,
			kafka.NewSendMessageSyncWithTopic(internalProducer),
			job.InsertRevertBank(),
		),
		job.UpdateUnMatedHeader(dbPool),
	)

	if err != nil {
		logger.Error("StageCheckFunc", zap.Any("err ", err))
	}
	fmt.Print("end ")
}
