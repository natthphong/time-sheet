package job

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/config"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/internal/logz"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"
)

func BackUpHisPricing(
	GetDataHisPricingFunc GetDataHisPricingFunc,
	PushToS3Func PushToS3Func,
	DetachPartitionHistoryFunc DetachPartitionHistoryFunc,
) error {

	logger := logz.NewLogger()
	ctx := context.Background()
	_ = ctx
	startPartition := os.Getenv("startPartition")
	//startPartition := "2024-04-01"
	if startPartition == "" {
		currentDate := time.Now()
		tempCurrentDate := currentDate.AddDate(0, -1, 0)
		partitions := tempCurrentDate.Format("_y2006m01")
		zipFile, err := GetDataHisPricingFunc(ctx, logger, partitions)
		if err != nil {
			logger.Error("Error GetDataHisPricingFunc", zap.Any("", err.Error()))
			return err
		}
		err = PushToS3Func(ctx, logger, zipFile, partitions)
		if err != nil {
			logger.Error("Error PushToS3Func", zap.Any("", err.Error()))
			return err
		}
		err = DetachPartitionHistoryFunc(ctx, logger, "")

		if err != nil {
			logger.Error("Error DetachPartitionHistoryFunc", zap.Any("", err.Error()))
			return err
		}
	} else {
		numStr := os.Getenv("numOfMonth")
		//numStr := "1"
		num, err := strconv.Atoi(numStr)
		if err != nil {
			logger.Error("Error parsing date:", zap.Error(err))
		}
		layout := "2006-01-02"
		currentDate, err := time.Parse(layout, startPartition)
		if err != nil {
			logger.Error("Error parsing date:", zap.Error(err))
		}
		for i := 0; i < num; i++ {
			tempCurrentDate := currentDate.AddDate(0, i, 0)
			partitions := tempCurrentDate.Format("_y2006m01")
			zipFile, err := GetDataHisPricingFunc(ctx, logger, partitions)
			if err != nil {
				logger.Error("Error GetDataHisPricingFunc", zap.Any("", err.Error()))
				return err
			}
			err = PushToS3Func(ctx, logger, zipFile, partitions)
			if err != nil {
				logger.Error("Error PushToS3Func", zap.Any("", err.Error()))
				return err
			}

		}
		err = DetachPartitionHistoryFunc(ctx, logger, "")

		if err != nil {
			logger.Error("Error DetachPartitionHistoryFunc", zap.Any("", err.Error()))
			return err
		}

	}

	return nil
}

type DetachPartitionHistoryFunc func(ctx context.Context, logger *zap.Logger, partition string) error

func DetachPartitionHistory(db *pgxpool.Pool) DetachPartitionHistoryFunc {
	return func(ctx context.Context, logger *zap.Logger, partition string) error {
		sql := `truncate table his_pricing`
		//sql := `ALTER TABLE his_pricing DETACH PARTITION his_pricing%s;`
		//sql = fmt.Sprintf(sql, partition)
		_, err := db.Exec(ctx, sql)
		if err != nil {
			return err
		}
		return nil
	}
}

type PushToS3Func func(ctx context.Context, logger *zap.Logger, zipFile bytes.Buffer, partitions string) error

func PushToS3(svc *s3.S3, cfg *config.Config) PushToS3Func {
	return func(ctx context.Context, logger *zap.Logger, zipFile bytes.Buffer, partitions string) error {
		key := fmt.Sprintf(cfg.S3Config.Key, partitions)
		_, err := svc.PutObject(&s3.PutObjectInput{
			Bucket: &cfg.S3Config.BucketName,
			Key:    &key,
			Body:   bytes.NewReader(zipFile.Bytes()),
		})
		if err != nil {
			return err
		}
		defer func() {
			zipFile.Reset()
		}()

		return nil
	}
}

type GetDataHisPricingFunc func(ctx context.Context, logger *zap.Logger, partitions string) (bytes.Buffer, error)

func GetDataHisPricing(db *pgxpool.Pool) GetDataHisPricingFunc {
	return func(ctx context.Context, logger *zap.Logger, partitions string) (bytes.Buffer, error) {
		tx, err := db.Begin(ctx)
		if err != nil {
			return bytes.Buffer{}, err
		}
		defer func(tx pgx.Tx) {
			_ = tx.Rollback(ctx)
		}(tx)
		sql := `
						select 'unix_created_time' || ',' ||
				'created_date' || ',' ||
				'request_ref' || ',' ||
				'buy_price' || ',' ||
				'sell_price' || ',' ||
				'request_time'
				union  all
				(select unix_created_time || ',' ||
				created_date || ',' ||
				request_ref || ',' ||
				buy_price || ',' ||
				sell_price || ',' ||
				request_time
				from his_pricing%s hp )
			`

		sql = fmt.Sprintf(sql, partitions)

		rows, err := tx.Query(ctx, sql)
		defer rows.Close()

		i := 0
		start := time.Now()
		logger.Info("start ", zap.Time("time", start))
		csvBuffer := &bytes.Buffer{}
		csvWriter := csv.NewWriter(csvBuffer)
		for rows.Next() {
			//startProcess := time.Now()
			var temp string
			err := rows.Scan(&temp)
			if err != nil {
				logger.Error("Error scanning  row", zap.Any("", err.Error()))
			}
			fields := strings.Split(temp, ",")
			_ = csvWriter.Write(fields)
			i++
			//if i%200000 == 0 {
			//	duration := time.Since(startProcess)
			//	logger.Debug("time use : ", zap.Any(strconv.Itoa(i), duration.Seconds()))
			//	//fmt.Print(i, " time use:", duration.Seconds(), "s  ")
			//}
		}

		csvWriter.Flush()
		zipBuffer := &bytes.Buffer{}
		zipWriter := zip.NewWriter(zipBuffer)

		csvFileName := fmt.Sprintf("his_pricing%s.csv", partitions)
		csvFile, err := zipWriter.Create(csvFileName)
		if err != nil {
			logger.Error("Error creating CSV file in zip archive:", zap.Error(err))
			return bytes.Buffer{}, err
		}
		_, err = csvFile.Write(csvBuffer.Bytes())
		if err != nil {
			logger.Error("Error writing CSV data to zip archive:", zap.Error(err))
			return bytes.Buffer{}, err
		}

		err = zipWriter.Close()
		if err != nil {
			logger.Error("Error closing zip archive:", zap.Error(err))
			return bytes.Buffer{}, err
		}
		duration := time.Since(start)

		defer func() {
			csvBuffer.Reset()
		}()
		logger.Info(fmt.Sprintf("his_pricing%s time to use : %s s", partitions, duration.Seconds()))
		return *zipBuffer, nil
	}

}
