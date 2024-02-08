package job

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/config"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/internal/sftp"
	"go.uber.org/zap"
	"regexp"
	"strings"
	"time"
)

type DailyKBankReconcile struct {
	TransactionBankID    string          `db:"transaction_bank_id"`
	RSTransID            string          `db:"rs_trans_id"`
	RequestDateTime      time.Time       `db:"request_date_time"`
	FundTransferDateTime time.Time       `db:"fund_transfer_date_time"`
	TransType            string          `db:"trans_type"`
	ProxyValue           string          `db:"proxy_value"`
	Amount               decimal.Decimal `db:"amount"`
	Ref1                 string          `db:"ref1"`
	Ref2                 string          `db:"ref2"`
	CreatedBy            string          `db:"created_by"`
	CreatedDate          time.Time       `db:"created_date"`
}

func GetFileInsertToTblTemp(
	cfg config.Config,
	ctx context.Context,
	logger *zap.Logger,
	s3Client *s3.S3,
	downLoadFileAndPushToS3Func DownLoadFileAndPushToS3Func,
	InsertKBANKDailyTempFunc InsertKBANKDailyTempFunc,
) error {

	var list []DailyKBankReconcile
	formattedDate := time.Now().Format("20060102")

	fileName := fmt.Sprintf(cfg.SFTPConfig.Destination.PrefixFileName, formattedDate)
	err := downLoadFileAndPushToS3Func(ctx, logger, fileName, cfg.SFTPConfig.Directory)
	if err != nil {
		return err
	}

	file, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(cfg.S3Config.BucketName),
		Key:    aws.String(cfg.SFTPConfig.Directory + "/" + fileName),
	})
	if err != nil {
		panic(err)
	}
	defer file.Body.Close()
	scanner := bufio.NewScanner(file.Body)

	for scanner.Scan() {
		line := scanner.Text()
		rows := strings.Split(line, "|")
		if len(rows) >= 9 {
			temp := DailyKBankReconcile{
				RSTransID:            rows[0],
				TransactionBankID:    rows[1],
				RequestDateTime:      ToTime(rows[2]),
				FundTransferDateTime: ToTime(rows[3]),
				TransType:            rows[4],
				ProxyValue:           rows[5],
				Amount:               ToDecimal(rows[6]),
				Ref1:                 rows[7],
				Ref2:                 rows[8],
				CreatedDate:          time.Now(),
				CreatedBy:            "SYSTEM",
			}
			list = append(list, temp)
		}

	}

	err = InsertKBANKDailyTempFunc(ctx, logger, list)
	return err

	return nil
}

func ToDecimal(str string) decimal.Decimal {
	v, _ := decimal.NewFromString(str)
	return v
}

func ToTime(str string) time.Time {
	//fmt.Println("str", str)
	//v, err := time.Parse(layout, str)
	//if err != nil {
	//	fmt.Println("parse date", err.Error())
	//}

	//str = "20220809 17:17:47:013"
	//layout = "20060102 15:04:05:000" // Specify the layout matching the string format
	//
	//v, err := time.Parse(layout, str)
	//if err != nil {
	//	fmt.Println("parse date", err.Error())
	//}
	layout := "20060102 15:04:05.000"
	re := regexp.MustCompile(`^(.*?:.*?:.*?):(.*?)$`)

	output := re.ReplaceAllString(str, "$1.$2")
	v, err := time.Parse(layout, output)
	if err != nil {
		fmt.Println("parse date", err.Error())
	}
	return v
}

type InsertKBANKDailyTempFunc func(ctx context.Context, logger *zap.Logger, list []DailyKBankReconcile) error

func InsertKBANKDailyTemp(db *pgxpool.Pool) InsertKBANKDailyTempFunc {
	return func(ctx context.Context, logger *zap.Logger, list []DailyKBankReconcile) error {
		rows := make([][]interface{}, len(list))
		ix := 0
		for _, l := range list {
			rows[ix] = []interface{}{
				l.TransactionBankID, l.RSTransID,
				l.RequestDateTime, l.FundTransferDateTime,
				l.TransType, l.ProxyValue,
				l.Amount, l.Ref1,
				l.Ref2, l.CreatedBy,
				l.CreatedDate,
			}
			ix++
		}

		i, err := db.CopyFrom(ctx, pgx.Identifier{"tbl_daily_kbank_reconcile"}, []string{"transaction_bank_id", "rs_trans_id", "request_date_time", "fund_transfer_date_time", "trans_type", "proxy_value", "amount", "ref1", "ref2", "created_by", "created_date"}, pgx.CopyFromRows(rows))
		fmt.Println("row insert ", i)
		if err != nil {
			return err
		}

		return nil
	}
}

type DownLoadFileAndPushToS3Func func(ctx context.Context, logger *zap.Logger, fileName, folder string) error

func DownLoadFileAndPushToS3(
	sftp *sftp.Client,
	svc *s3.S3,
	putFileToS3 PutFileToS3Func,
) DownLoadFileAndPushToS3Func {
	return func(ctx context.Context, logger *zap.Logger, fileName, folder string) error {
		key := fmt.Sprintf("%s/%s", folder, fileName)

		//listFile, err := sftp.ListFiles(folder)
		//if err != nil {
		//	logger.Error(err.Error())
		//}
		//logger.Info("smft", zap.Any("file", listFile))
		byteFile, err := sftp.Download(key)

		fileReader := bytes.NewReader(byteFile)
		err = putFileToS3(fileName, folder, fileReader, svc)
		if err != nil {
			return err
		}
		return nil
	}

}

type PutFileToS3Func func(filename, key string, file *bytes.Reader, svc *s3.S3) error

func PutFileToS3(bucket string) PutFileToS3Func {
	return func(filename, key string, file *bytes.Reader, svc *s3.S3) error {
		_, err := svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fmt.Sprintf("%s/%s", key, filename)),
			Body:   file,
		})
		if err != nil {
			return err
		}
		return nil
	}
}
