package eft

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/prior-solution/aurora/standard-platform/common/reconcile_daily_batch/config"
	"golang.org/x/exp/slices"
	"os"

	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

var exceptionInquiryStatus = []string{"Fail", "In Process"}

const (
	SuccessFundTransfer        = "0000"
	OtherExceptionFundTransfer = "9999"
	OauthSuccess               = "approved"
)

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
}

type AccessTokenResponse struct {
	DeveloperEmail string `json:"developer.email"`
	TokenType      string `json:"token_type"`
	ClientID       string `json:"client_id"`
	AccessToken    string `json:"access_token"`
	Scope          string `json:"scope"`
	ExpiresIn      string `json:"expires_in"`
	Status         string `json:"status"`
}

type VerifyDataFundTransferRequest struct {
	MerchantID      string `json:"merchantID"`
	MerchantTransID string `json:"merchantTransID"`
	RequestDateTime string `json:"requestDateTime"`
	ProxyType       string `json:"proxyType"`
	ProxyValue      string `json:"proxyValue"`
	FromAccountNo   string `json:"fromAccountNo"`
	TransType       string `json:"transType"`
	SenderName      string `json:"senderName"`
	SenderTaxID     string `json:"senderTaxID"`
	ToBankCode      string `json:"toBankCode"`
	Amount          string `json:"amount"`
	TypeOfSender    string `json:"typeOfSender"`
}

type VerifyDataFundTransferResponse struct {
	MerchantID       string  `json:"merchantID"`
	MerchantTransID  string  `json:"merchantTransID"`
	RsTransID        string  `json:"rsTransID"`
	ResponseDateTime string  `json:"responseDateTime"`
	ResponseCode     string  `json:"responseCode"`
	ResponseMsg      string  `json:"responseMsg"`
	ProxyType        string  `json:"proxyType"`
	ProxyValue       string  `json:"proxyValue"`
	ToBankCode       string  `json:"toBankCode"`
	ToAccNameTH      string  `json:"toAccNameTH"`
	ToAccNameEN      string  `json:"toAccNameEN"`
	TransType        string  `json:"transType"`
	FromAccountNo    string  `json:"fromAccountNo"`
	SenderName       string  `json:"senderName"`
	SenderTaxID      string  `json:"senderTaxID"`
	TypeOfSender     string  `json:"typeOfSender"`
	Amount           float64 `json:"amount"`
}

type FundTransferRequest struct {
	MerchantID       string `json:"merchantID"`
	RequestDateTime  string `json:"requestDateTime"`
	MerchantTransID  string `json:"merchantTransID"`
	RsTransID        string `json:"rsTransID"`
	CustomerMobileNo string `json:"customerMobileNo"`
}

type FundTransferResponse struct {
	MerchantID       string `json:"merchantID"`
	MerchantTransID  string `json:"merchantTransID"`
	RsTransID        string `json:"rsTransID"`
	ResponseDateTime string `json:"responseDateTime"`
	ResponseCode     string `json:"responseCode"`
	ResponseMsg      string `json:"responseMsg"`
	SettlementDate   string `json:"settlementDate"`
}

type InquiryStatusRequest struct {
	MerchantID      string `json:"merchantID"`
	RequestDateTime string `json:"requestDateTime"`
	MerchantTransID string `json:"merchantTransID"`
	RsTransID       string `json:"rsTransID"`
}

type InquiryStatusResponse struct {
	MerchantID       string `json:"merchantID"`
	MerchantTransID  string `json:"merchantTransID"`
	RsTransID        string `json:"rsTransID"`
	ResponseDateTime string `json:"responseDateTime"`
	ResponseCode     string `json:"responseCode"`
	ResponseMsg      string `json:"responseMsg"`
	TxnStatus        string `json:"txnStatus"`
	SettlementDate   string `json:"settlementDate"`
	FailMsg          string `json:"failMsg"`
}

type HTTPFundTransferFunc func(logger *zap.Logger, req FundTransferRequest, accessToken string) (*FundTransferResponse, error)

func HTTPFundTransfer(client *http.Client, url string, toggle config.ToggleConfiguration) HTTPFundTransferFunc {
	return func(logger *zap.Logger, req FundTransferRequest, accessToken string) (*FundTransferResponse, error) {
		if toggle.IsTest {
			switch toggle.Case {
			case "P":
				return &FundTransferResponse{
					MerchantID:       req.MerchantID,
					MerchantTransID:  req.MerchantTransID,
					RsTransID:        uuid.NewString(),
					ResponseDateTime: time.Now().Format(time.RFC3339),
					ResponseCode:     "0000",
					ResponseMsg:      "Success",
					SettlementDate:   time.Now().Format("20060102"),
				}, nil
			case "F":
				return nil, os.ErrNotExist
			}
		}

		var (
			httpRes  *http.Response
			err      error
			response *FundTransferResponse
		)

		requestBodyJSON, err := json.Marshal(&req)
		if err != nil {
			return nil, err
		}

		bearer := fmt.Sprintf("Bearer %s", accessToken)
		httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBodyJSON))
		if err != nil {
			return nil, fmt.Errorf("unable to New http request: %v", err)
		}
		httpReq.Header.Set("Authorization", bearer)
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("env-id", "OAUTH2")

		httpRes, err = client.Do(httpReq)
		if err != nil {
			logger.Error("Error on call http request", zap.Error(err))
			return nil, err
		}

		if httpRes != nil {
			defer httpRes.Body.Close()

			if httpRes.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("error call %s body: %v", url, httpRes.Body)
			}

			body, err := io.ReadAll(httpRes.Body)
			if err != nil {
				logger.Error("Error on read response", zap.Error(err))
				return nil, err
			}

			err = json.Unmarshal(body, &response)

			return response, err
		}

		return nil, os.ErrNotExist
	}
}

type HTTPOauthFundTransferHttpFunc func(logger *zap.Logger, auth string, wait time.Duration) (*AccessTokenResponse, error)

func HTTPOauthFundTransferHttp(client *http.Client, url string, toggle config.ToggleConfiguration, retry int) HTTPOauthFundTransferHttpFunc {
	return func(logger *zap.Logger, auth string, wait time.Duration) (*AccessTokenResponse, error) {
		if toggle.IsTest {
			switch toggle.Case {
			case "P":
				return &AccessTokenResponse{
					DeveloperEmail: "FreedomX-10@hotmail.com",
					TokenType:      "Bearer",
					ClientID:       "t3rrPWnrt2jsOdjFrliIJcPslE76q09B",
					AccessToken:    "AccessToken",
					Scope:          "Any",
					ExpiresIn:      "1799",
					Status:         "approved",
				}, nil
			case "F":
				return nil, errors.New("error on HTTPOauthFundTransferHttp")
			}
		}

		var (
			httpRes  *http.Response
			err      error
			response *AccessTokenResponse
		)

		basicAuth := fmt.Sprintf("Basic %s", auth)
		data := "grant_type=client_credentials"
		dataByte := []byte(data)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(dataByte))
		if err != nil {
			return nil, fmt.Errorf("unable to New http request: %v", err)
		}
		req.Header.Set("Authorization", basicAuth)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("env-id", "OAUTH2")

		newRetry := retry

		for newRetry > 0 {
			httpRes, err = client.Do(req)
			if err != nil {
				logger.Error("Error on call http request", zap.Error(err))
				retry--
				time.Sleep(wait)
				continue
			}

			if httpRes != nil {
				if !(httpRes.StatusCode < 200 && httpRes.StatusCode > 299) {
					logger.Error(fmt.Sprintf("HTTP status code out of range (%d)", httpRes.StatusCode))
					retry--
					time.Sleep(wait)
					continue
				}
				body, err := io.ReadAll(httpRes.Body)
				if err != nil {
					logger.Error("Error on read response", zap.Error(err))
					retry--
					time.Sleep(wait)
					continue
				}

				err = json.Unmarshal(body, &response)
				if err != nil {
					logger.Error("Unmarshal", zap.Error(err))
					retry--
					time.Sleep(wait)
					continue
				}

				if response.Status != OauthSuccess {
					logger.Error(fmt.Sprintf("Call Oauth node success (%s)", response.Status))
					retry--
					time.Sleep(wait)
					continue
				}
			}

		}

		if response != nil {
			return response, nil
		}

		return nil, fmt.Errorf("unable to request %s", url)
	}
}

type HTTPInquiryStatusFundTransferFunc func(logger *zap.Logger, req InquiryStatusRequest, accessToken string, wait time.Duration) (*InquiryStatusResponse, error)

func HTTPInquiryStatusFundTransfer(client *http.Client, url string, toggle config.ToggleConfiguration, retry int) HTTPInquiryStatusFundTransferFunc {
	return func(logger *zap.Logger, req InquiryStatusRequest, accessToken string, wait time.Duration) (*InquiryStatusResponse, error) {
		if toggle.IsTest {
			switch toggle.Case {
			case "P":
				return &InquiryStatusResponse{
					MerchantID:       req.MerchantID,
					MerchantTransID:  req.MerchantTransID,
					RsTransID:        uuid.NewString(),
					ResponseDateTime: time.Now().Format(time.RFC3339),
					ResponseCode:     "0000",
					ResponseMsg:      "Success",
					TxnStatus:        "Success",
					SettlementDate:   time.Now().Format("20060102"),
					FailMsg:          "FailMsg",
				}, nil
			case "F":
				return nil, errors.New("error on fund transfer")
			}
		}
		var (
			httpRes  *http.Response
			httpReq  *http.Request
			err      error
			response *InquiryStatusResponse
		)

		requestBodyJSON, err := json.Marshal(&req)
		if err != nil {
			return nil, err
		}

		bearer := fmt.Sprintf("Bearer %s", accessToken)
		newRetry := retry

		for newRetry > 0 {
			httpReq, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBodyJSON))
			if err != nil {
				return nil, fmt.Errorf("unable to New http request: %v", err)
			}

			httpReq.Header.Set("Authorization", bearer)
			httpReq.Header.Set("Content-Type", "application/json")

			httpRes, err = client.Do(httpReq)
			if err != nil {
				logger.Error("Error on call http request", zap.Error(err))
				newRetry--
				time.Sleep(wait)
				continue
			}
			if httpRes != nil {
				if httpRes.StatusCode < 200 && httpRes.StatusCode > 299 {
					logger.Error(fmt.Sprintf("HTTP status code out of range (%d)", httpRes.StatusCode))
					newRetry--
					time.Sleep(wait)
					continue
				}
				body, err := io.ReadAll(httpRes.Body)
				if err != nil {
					logger.Error("Error on read response", zap.Error(err))
					newRetry--
					time.Sleep(wait)
					continue
				}

				err = json.Unmarshal(body, &response)
				if err != nil {
					logger.Error("Unmarshal", zap.Error(err))
					newRetry--
					time.Sleep(wait)
					continue
				}

				if slices.Contains(exceptionInquiryStatus, response.TxnStatus) {
					logger.Error(fmt.Sprintf("transaction is (%s)", response.TxnStatus))
					newRetry--
					time.Sleep(wait)
					continue
				}
			}
		}

		return response, err
	}
}
