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
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/internal/db"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/internal/logz"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/job"
	"go.uber.org/zap"
	"log"
)

func main() {

	lambda.Start(LambdaHandler)
	//LambdaHandler()
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

	//if cfg.Env != "" {
	//	err = secret.ConfigCommonSecret(cfg)
	//	if err != nil {
	//		log.Fatal(errors.Wrap(err, "Unable to initial common secret"))
	//	}
	//	err = secret.ConfigRDSSecret(cfg)
	//	if err != nil {
	//		log.Fatal(errors.Wrap(err, "Unable to initial rds secret"))
	//	}
	//}

	jsonCfg, err := json.Marshal(cfg)
	_ = jsonCfg
	fmt.Println("after cfg : ", string(jsonCfg))
	dbPool, err := db.Open(ctx, cfg.DBConfig)
	if err != nil {
		logger.Fatal("server connect to db", zap.Error(err))
	}
	defer dbPool.Close()

	//redisClient, err := cache.Initialize(ctx, cfg.RedisConfig)
	//if err != nil {
	//	log.Fatal(errors.Wrap(err, "Error cannot connect redis."))
	//}

	//defer redisClient.Close()
	//redisCmd := redisClient.CMD()

	//
	//internalProducer, err := scramkafka.NewSyncProducer(cfg.Kafka.Internal)
	//if err != nil {
	//	logger.Fatal("Fail Create NewSyncProducer", zap.Error(err))
	//}
	//defer func() {
	//	if err = internalProducer.Close(); err != nil {
	//		logger.Fatal("Fail Close SyncProducer", zap.Error(err))
	//	}
	//}()
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	})
	if err != nil {
		log.Fatal(errors.Wrap(err, "Unable to initial config."))
		return
	}
	svc := s3.New(sess)
	_ = svc

	job.BackUpHisPricing(
		job.GetDataHisPricing(dbPool),
		job.PushToS3(svc, cfg),
	)
	fmt.Print("end ")
}
