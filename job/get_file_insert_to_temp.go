package job

import (
	"bufio"
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
	sftp *sftp.Client,
	InsertKBANKDailyTempFunc InsertKBANKDailyTempFunc,
) error {
	layout := "20060102 15:04:05:000"
	var list []DailyKBankReconcile
	if cfg.EnableS3 {
		file, err := s3Client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(cfg.S3Config.BucketName),
			Key:    aws.String(cfg.S3Config.Key),
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
					RequestDateTime:      ToTime(layout, rows[2]),
					FundTransferDateTime: ToTime(layout, rows[3]),
					TransType:            rows[4],
					ProxyValue:           rows[5],
					Amount:               ToDecimal(rows[6]),
					Ref1:                 rows[7],
					Ref2:                 rows[8],
					CreatedDate:          time.Now(),
					CreatedBy:            "SYSTEM",
				}
				list = append(list, temp)
				//TODO
			}

		}

	} else {

	}

	err := InsertKBANKDailyTempFunc(ctx, logger, list)
	return err

	return nil
}

func ToDecimal(str string) decimal.Decimal {
	v, _ := decimal.NewFromString(str)
	return v
}

func ToTime(layout, str string) time.Time {
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
	return time.Now()
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
