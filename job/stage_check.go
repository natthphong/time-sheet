package job

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/config"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/eft"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/internal/cache"
	"go.uber.org/zap"
	"time"
)

type BankUnmatchedDetail struct {
	UnmatchedHeaderID int
	ChannelCode       string
	TransactionBankID string
	BankStatus        string
	ASPStatus         string
	Reason            string
	CreatedDate       time.Time
}

type ResultStruct struct {
	TransactionResultId decimal.Decimal `json:"transactionResultId"`
	AspBankId           string          `json:"aspBankId"`
	KBankId             string          `json:"kbankId"`
	Status              string          `json:"status"`
	ChannelCode         string          `json:"channelCode"`
	RsTransID           string          `json:"rsTransId"`
}

func StageCheckFunc(
	ctx context.Context,
	logger *zap.Logger,
	InsertUnMatedHeaderFunc InsertUnMatedHeaderFunc,
	GetListResultFunc GetListResultFunc,
	InsertUnMatedDetailFunc InsertUnMatedDetailFunc,
) error {
	id, err := InsertUnMatedHeaderFunc(ctx, logger)
	rs, err := GetListResultFunc(ctx, logger)

	err = InsertUnMatedDetailFunc(ctx, logger, rs, id)
	fmt.Println(err.Error())
	//TODO insert detail
	return err
}

type InsertUnMatedDetailFunc func(ctx context.Context, logger *zap.Logger, rs []ResultStruct, id int) error

func InsertUnMatedDetail(
	fundTransferConfig config.FundTransferConfig,
	exceptionConfig config.Exception,
	getRedis cache.GetRedisFunc,
	oauthFundTransferHttp eft.HTTPOauthFundTransferHttpFunc,
	inquiryStatusFundTransferHttp eft.HTTPInquiryStatusFundTransferFunc,
) InsertUnMatedDetailFunc {
	return func(ctx context.Context, logger *zap.Logger, rs []ResultStruct, id int) error {

		var unmatchedDetails []BankUnmatchedDetail

		fundTransferTokenKey := cache.FundTransferTokenKey
		accessToken, err := getRedis(ctx, fundTransferTokenKey)
		if err != nil {
			tokenResponse, err := oauthFundTransferHttp(logger, fundTransferConfig.Auth, time.Second*45)
			if err != nil {
				//TODO
				_ = exceptionConfig.Description.SystemError
				_ = exceptionConfig.Code.SystemError
				logger.Error("callOauthFundTransferHttp", zap.Error(err))
			}
			if tokenResponse != nil {
				accessToken = tokenResponse.AccessToken
			}
		}

		if accessToken != "" {
			auth := fmt.Sprintf("Bearer %s", accessToken)
			for _, r := range rs {
				startTime := time.Now()
				requestDateTime := startTime.Format(time.RFC3339)
				if r.KBankId == "" {
					//case TimeOut
					if !r.TransactionResultId.IsNegative() {
						//TODO sendToRevert
						temp := BankUnmatchedDetail{
							UnmatchedHeaderID: id,
							ChannelCode:       r.ChannelCode,
							TransactionBankID: r.KBankId,
							BankStatus:        "",
							ASPStatus:         r.Status,
							Reason:            "REVERT",
							CreatedDate:       startTime,
						}
						unmatchedDetails = append(unmatchedDetails, temp)
						logger.Info("REVERT CASE TIME OUT", zap.Any("json", r))
					}
				} else {
					inquiryStatusRequest := eft.InquiryStatusRequest{
						MerchantID:      fundTransferConfig.MerchantID,
						RequestDateTime: requestDateTime,
						MerchantTransID: r.KBankId,
						RsTransID:       r.RsTransID,
					}
					inquiryStatusRes, err := inquiryStatusFundTransferHttp(logger, inquiryStatusRequest, auth, time.Second*45)
					if err != nil {
						//TODO
					}
					status := inquiryStatusRes.TxnStatus
					if status == "FAIL" {
						temp := BankUnmatchedDetail{
							UnmatchedHeaderID: id,
							ChannelCode:       r.ChannelCode,
							TransactionBankID: r.KBankId,
							BankStatus:        status,
							ASPStatus:         r.Status,
							Reason:            "FAIL",
							CreatedDate:       startTime,
						}
						unmatchedDetails = append(unmatchedDetails, temp)
						logger.Info("REVERT CASE STATUS FAIL", zap.Any("json", r))
						//TODO REVERT
					} else if r.TransactionResultId.IsNegative() {
						temp := BankUnmatchedDetail{
							UnmatchedHeaderID: id,
							ChannelCode:       r.ChannelCode,
							TransactionBankID: r.KBankId,
							BankStatus:        status,
							ASPStatus:         r.Status,
							Reason:            "SUCCESS",
							CreatedDate:       startTime,
						}
						unmatchedDetails = append(unmatchedDetails, temp)
						logger.Info("TOO FINAL", zap.Any("json", r))
						//TODO FINAL
					}

				}

			}

			json, _ := json2.Marshal(unmatchedDetails)
			logger.Info("unmatchedDetails", zap.Any("json", string(json)))

		}
		return nil
	}
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
			,COALESCE( tdkr.rs_trans_id  , '') as rsTransId
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
			err := rows.Scan(&result.TransactionResultId, &result.AspBankId, &result.KBankId, &result.Status, &result.ChannelCode, &result.RsTransID)
			if err != nil {
				logger.Error("Error scanning  row", zap.Any("", err.Error()))
				return []ResultStruct{}, err
			}
			results = append(results, result)
		}

		return results, nil
	}
}
