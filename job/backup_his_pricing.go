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
	"strconv"
	"strings"
	"time"
)

func BackUpHisPricing(
	GetDataHisPricingFunc GetDataHisPricingFunc,
	PushToS3Func PushToS3Func,
) {
	logger := logz.NewLogger()
	ctx := context.Background()
	zipFile, err := GetDataHisPricingFunc(ctx, logger)
	if err != nil {
		logger.Error("Error GetDataHisPricingFunc", zap.Any("", err.Error()))
	}
	err = PushToS3Func(ctx, logger, zipFile)
	if err != nil {
		logger.Error("Error PushToS3Func", zap.Any("", err.Error()))
	}

}

type PushToS3Func func(ctx context.Context, logger *zap.Logger, zipFile bytes.Buffer) error

func PushToS3(svc *s3.S3, cfg *config.Config) PushToS3Func {
	return func(ctx context.Context, logger *zap.Logger, zipFile bytes.Buffer) error {

		formattedDate := time.Now().Format("20060102")
		key := fmt.Sprintf(cfg.S3Config.Key, formattedDate)
		_, err := svc.PutObject(&s3.PutObjectInput{
			Bucket: &cfg.S3Config.BucketName,
			Key:    &key,
			Body:   bytes.NewReader(zipFile.Bytes()),
		})
		if err != nil {
			return err
		}

		return nil
	}
}

type GetDataHisPricingFunc func(ctx context.Context, logger *zap.Logger) (bytes.Buffer, error)

func GetDataHisPricing(db *pgxpool.Pool) GetDataHisPricingFunc {
	return func(ctx context.Context, logger *zap.Logger) (bytes.Buffer, error) {
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
				from his_pricing hp )
			`

		rows, err := tx.Query(ctx, sql)
		defer rows.Close()

		i := 0
		start := time.Now()
		logger.Info("start ", zap.Time("time", start))
		csvBuffer := &bytes.Buffer{}
		csvWriter := csv.NewWriter(csvBuffer)
		for rows.Next() {
			startProcess := time.Now()
			var temp string
			err := rows.Scan(&temp)
			if err != nil {
				logger.Error("Error scanning  row", zap.Any("", err.Error()))
			}
			fields := strings.Split(temp, ",")
			_ = csvWriter.Write(fields)
			i++
			if i%200000 == 0 {
				duration := time.Since(startProcess)
				logger.Debug("time use : ", zap.Any(strconv.Itoa(i), duration.Seconds()))
				//fmt.Print(i, " time use:", duration.Seconds(), "s  ")
			}
		}

		csvWriter.Flush()

		zipBuffer := &bytes.Buffer{}

		zipWriter := zip.NewWriter(zipBuffer)

		csvFile, err := zipWriter.Create("his_pricing.csv")
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
		fmt.Print("\n end time use:", duration.Seconds(), "s  ")
		return *zipBuffer, nil
	}

}
