package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(LambdaHandler)
	//config.InitTimeZone()
	//_, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//
	//cfg, err := config.InitConfig()
	//if err != nil {
	//	log.Fatal(errors.Wrap(err, "Unable to initial config."))
	//}
	//
	//logz.Init(cfg.Log.Level, cfg.Server.Name)
	//defer logz.Drop()
	//
	//ctx := context.Background()
	//ctx, cancel = context.WithCancel(ctx)
	//defer cancel()
	//logger := zap.L()

	//err = secret.ConfigCommonSecret(cfg)
	//if err != nil {
	//	log.Fatal(errors.Wrap(err, "Unable to initial common secret"))
	//}
	//err = secret.ConfigRDSSecret(cfg)
	//if err != nil {
	//	log.Fatal(errors.Wrap(err, "Unable to initial rds secret"))
	//}

	//dbPool, err := db.Open(ctx, cfg.DBConfig)
	//if err != nil {
	//	logger.Fatal("server connect to db", zap.Error(err))
	//}
	//defer dbPool.Close()
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
	//TODO
}

func LambdaHandler(ctx context.Context) {
	fmt.Print("hello ")
}
