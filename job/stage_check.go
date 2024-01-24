package job

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"time"
)

type ResultStruct struct {
	TransactionResultId decimal.Decimal `json:"transactionResultId"`
	AspBankId           string          `json:"aspBankId"`
	KBankId             string          `json:"kbankId"`
	Status              string          `json:"status"`
	ChannelCode         string          `json:"channelCode"`
}

func StageCheckFunc(
	ctx context.Context,
	logger *zap.Logger,
	InsertUnMatedHeaderFunc InsertUnMatedHeaderFunc,
	GetListResultFunc GetListResultFunc,
) error {
	id, err := InsertUnMatedHeaderFunc(ctx, logger)
	logger.Info("Before", zap.Any("id", id))
	result, err := GetListResultFunc(ctx, logger)
	jsonStr, _ := json.Marshal(result)
	logger.Info("After get", zap.Any("json", string(jsonStr)))

	return err
}

type InsertUnMatedHeaderFunc func(ctx context.Context, logger *zap.Logger) (int, error)

func InsertUnMatedHeader(db *pgxpool.Pool) InsertUnMatedHeaderFunc {

	return func(ctx context.Context, logger *zap.Logger) (int, error) {
		var id int
		t := time.Now()
		sql := `
				INSERT INTO tbl_bank_unmatched_header
			(unmatched_date, unmatched_time, status, created_date)
			VALUES($1, $2, 'PENDING', now()) returning unmatched_header_id;
			`

		err := db.QueryRow(ctx, sql, t.Format("20060102"), t.Format("1504")).Scan(&id)
		return id, err
	}

}

type GetListResultFunc func(ctx context.Context, logger *zap.Logger) ([]ResultStruct, error)

func GetListResult(db *pgxpool.Pool) GetListResultFunc {
	return func(ctx context.Context, logger *zap.Logger) ([]ResultStruct, error) {
		sql := `
			select COALESCE( ttr.transaction_id ,-1)as transactionResultId
			, trb.transaction_bank_id   as aspBankId
			, COALESCE( tdkr.transaction_bank_id , '') as kbankId
			, COALESCE( ttr.status ,'')
			, COALESCE( tt.channel_code ,'')as channelCode
			from tbl_reconcile_bank trb   left outer join tbl_daily_kbank_reconcile tdkr 
			on trb.transaction_bank_id  = tdkr.transaction_bank_id  left outer join tbl_transaction tt 
			on trb.transaction_id  = tt.transaction_id left outer join tbl_transaction_result ttr
			on tt.transaction_id  = ttr.transaction_id 
		`
		rows, err := db.Query(ctx, sql)
		if err != nil {
			logger.Error("Error executing ", zap.Any("", err.Error()))
			return []ResultStruct{}, err
		}
		defer rows.Close()

		var results []ResultStruct
		for rows.Next() {
			var result ResultStruct
			err := rows.Scan(&result.TransactionResultId, &result.AspBankId, &result.KBankId, &result.Status, &result.ChannelCode)
			if err != nil {
				logger.Error("Error scanning  row", zap.Any("", err.Error()))
				return []ResultStruct{}, err
			}
			results = append(results, result)
		}

		return results, nil
	}
}
