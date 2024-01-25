package job

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

type PaymentMessage struct {
	TransactionId                decimal.Decimal `json:"TransactionId"`
	TransactionType              string          `json:"TransactionType"`
	OrderType                    string          `json:"OrderType"`
	RequestRef                   string          `json:"RequestRef"`
	StatusCode                   string          `json:"StatusCode"`
	StatusDescription            string          `json:"StatusDescription"`
	ChannelCode                  string          `json:"ChannelCode"`
	OrderRequest                 json.RawMessage `json:"OrderRequest"`
	OrderResponse                json.RawMessage `json:"OrderResponse"`
	BankingFeature               string          `json:"BankingFeature"`
	TopicResponse                string          `json:"TopicResponse"`
	ToRevert                     bool            `json:"ToRevert"`
	ChannelFeature               string          `json:"ChannelFeature"`
	CustomerProfile              json.RawMessage `json:"CustomerProfile"`
	FundTransferTransactionModel json.RawMessage `json:"FundTransferTransaction"`
}
