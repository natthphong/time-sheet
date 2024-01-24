package job

import (
	"encoding/json"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/internal/httputil"
)

type CallInquiryStatusFundTransferHttp func(url string, req InquiryStatusRequest, token string) (*InquiryStatusResponse, error)

func NewCallInquiryStatusFundTransferHttp(httpClientFunc httputil.HTTPPostPaymentRequestFunc) CallInquiryStatusFundTransferHttp {
	return func(url string, req InquiryStatusRequest, token string) (*InquiryStatusResponse, error) {

		var response *InquiryStatusResponse

		requestBodyJSON, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}

		result, err := httpClientFunc(url, requestBodyJSON, token, "application/json")
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(result, &response)
		if err != nil {
			return nil, err
		}

		//if response.TxnStatus != "Success" {
		//	return nil, fmt.Errorf("status: %s, body: %s", response.ResponseCode, response)
		//}
		return response, nil
	}
}
