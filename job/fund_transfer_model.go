package job

const (
	SuccessFundTransfer        = "0000"
	OtherExceptionFundTransfer = "9999"
	OauthSuccess               = "approved"
)

const (
	OauthFundTransferUrl         = "https://openapi-test.kasikornbank.com/v2/oauth/token"
	VerifyFundTransferUrl        = "https://openapi-test.kasikornbank.com/v1/fundtransfer/verifydata"
	FundTransferUrl              = "https://openapi-test.kasikornbank.com/v1/fundtransfer/fundtransfer"
	InquiryStatusFundTransferUrl = "https://openapi-test.kasikornbank.com/v1/fundtransfer/inqtxnstatus"

	MerchantID    = "ARRT"
	FromAccountNo = "0481418100"
	Auth          = "Basic dDNyclBXbnJ0MmpzT2RqRnJsaUlKY1BzbEU3NnEwOUI6NURoVDVJc2wyalVPZ0ZHUA=="

	SenderName = "AURORA TRADING CO.LTD."
	//SenderTaxID = "0105553020335"

	/*
		Test Case 01
		MSG REQUEST
		fromAccountNo  = 048-1-41810-0
		ProxyType =   10
		ProxyValue =   0118324366
		toBackCode  = 004
		amount  =  depend on merchant
		senderTaxID =  TAX ID of merchant
	*/

	ProxyType    = "10"
	ProxyValue   = "0118324366"
	TransType    = "K2K"
	ToBankCode   = "004"
	Amount       = "1.99"
	TypeOfSender = "K"
	SenderTaxID  = "1540907567898"

	/*
		Test Case 02
		ทดสอบโอนเงินสำเร็จ
		-กรณีโอนเงินไปยังบัญชีปลายทาง ธนาคารอื่น (K-to-O)
		ด้วยหมายเลขบัญชีธนาคาร

		MSG REQUEST
		fromAccountNo  = 048-1-41810-0
		ProxyType =   10
		ProxyValue =   020000159010
		toBackCode  = 030
		amount  =  depend on merchant
		senderTaxID =  TAX ID of merchant
	*/

	//ProxyType    = "10"
	//ProxyValue   = "020000159010"
	//TransType    = "K2O"
	//ToBankCode   = "030"
	//Amount       = "1.99"
	//TypeOfSender = "K"

	/*
		Test Case 3
		MSG REQUEST
		fromAccountNo  = 048-1-41810-0
		ProxyType =   10
		ProxyValue =   0019562702
		toBackCode  = 025
		amount  =  depend on merchant
		senderTaxID =  TAX ID of merchant
	*/

	//ProxyType    = "10"
	//ProxyValue   = "0019562702"
	//TransType    = "K2O"
	//ToBankCode   = "025"
	//Amount       = "1.99"
	//TypeOfSender = "K"

	/*
		Test Case 4
		MSG REQUEST
		fromAccountNo  = 048-1-41810-0
		ProxyType =   10
		ProxyValue =   0010004901
		toBackCode  = 002
		amount  =  depend on merchant
		senderTaxID =  TAX ID of merchant
	*/
	//ProxyType    = "10"
	//ProxyValue   = "0010004901"
	//TransType    = "K2O"
	//ToBankCode   = "002"
	//Amount       = "1.99"
	//TypeOfSender = "K"

	/*
		Test Case 5
		MSG REQUEST
		fromAccountNo  = 048-1-41810-0
		ProxyType =   02
		ProxyValue =   0814428232
		toBackCode  = N/A
		amount  =  depend on merchant
		senderTaxID =  TAX ID of merchant
	*/
	//ProxyType    = "02"
	//ProxyValue   = "0814428232"
	//TransType    = "K2O"
	//ToBankCode   = ""
	//Amount       = "1.99"
	//TypeOfSender = "K"

	/*
		Test Case 6
		MSG REQUEST
		fromAccountNo  = 048-1-41810-0
		ProxyType =   01
		ProxyValue =   1121821823009
		toBackCode  = N/A
		amount  =  depend on merchant
		senderTaxID =  TAX ID of merchant
	*/
	//ProxyType    = "01"
	//ProxyValue   = "1121821823009"
	//TransType    = "K2O"
	//ToBankCode   = ""
	//Amount       = "1.99"
	//TypeOfSender = "K"

	/*
		Test Case 7
		ทดสอบทำรายการโอนเงิน กรณีไม่สำเร็จ เนื่องจาก ทำรายการโอนเงินต่อรายการ เกิน Limit Amount / Transaction
		 -กรณีโอนเงินไปยังบัญชีปลายทาง KBANK (K-to-K)

		MSG REQUEST
		fromAccountNo  = 048-1-41810-0
		ProxyType =   10
		ProxyValue =   0011248136
		toBackCode  = 004
		amount  =  5,100,000 THB
		senderTaxID =  TAX ID of merchant

		MSG RESPONSE
		Response code = 1008 (Over limit amount/transaction)
	*/

	//ProxyType    = "10"
	//ProxyValue   = "0011248136"
	//TransType    = "K2K"
	//ToBankCode   = "004"
	//Amount       = "5100000.00"
	//TypeOfSender = "K"

	/*
		Test Case 8
		ทดสอบทำรายการโอนเงิน กรณีไม่สำเร็จ เนื่องจาก ทำรายการโอนเงินต่อรายการ เกิน Limit Amount / Transaction
		 -กรณีโอนเงินไปยังบัญชีปลายทาง ธนาคารอื่น (K-to-O)

		MSG REQUEST
		fromAccountNo  = 048-1-41810-0
		ProxyType =   10
		ProxyValue =   7770002696
		toBackCode  = 025
		amount  =  3,100,000 THB
		senderTaxID =  TAX ID of merchant

		MSG RESPONSE
		 Response code = 1008 (Over limit amount/transaction)
	*/

	//ProxyType    = "10"
	//ProxyValue   = "7770002696"
	//TransType    = "K2O"
	//ToBankCode   = "025"
	//Amount       = "3100000.00"
	//TypeOfSender = "K"

	/*
			Test Case 9
			ทดสอบทำรายการโอนเงิน กรณีไม่สำเร็จ เนื่องจาก บัญชีปลายทางมีปัญหา ไม่สามารถทำรายการโอนเงินได้
			 -กรณีโอนเงินไปยังบัญชีปลายทาง KBANK (K-to-K)

			MSG REQUEST
			fromAccountNo  =  048-1-41810-0
			ProxyType =   10
			ProxyValue =   0013346917
			toBackCode  = 004
			amount  =  100,000 THB
			senderTaxID =  TAX ID of merchant

			MSG RESPONSE
		    Response code = 4000 (Bank Account is not allow to process)
	*/

	//ProxyType    = "10"
	//ProxyValue   = "0013346917"
	//TransType    = "K2K"
	//ToBankCode   = "004"
	//Amount       = "100000.00"
	//TypeOfSender = "K"
	//SenderTaxID  = "3410200567196"

	/*
			Test case 10
			ทดสอบทำรายการโอนเงิน กรณีไม่สำเร็จ เนื่องจาก บัญชี Operative account ไม่พอตัดเงิน
			 -กรณีโอนเงินไปยังบัญชีปลายทาง KBANK (K-to-K)

			MSG REQUEST
			fromAccountNo  = 0011232329
			ProxyType =   10
			ProxyValue =   0503653931
			toBackCode  = 004
			amount  =  1,000,000 THB
			senderTaxID =  TAX ID of merchant

			MSG RESPONSE
		    Response code = 4006 (Insufficient Balance)
	*/

	//ProxyType     = "10"
	//ProxyValue    = "0503653931"
	//TransType     = "K2K"
	//ToBankCode    = "004"
	//Amount        = "1000000.00"
	//TypeOfSender  = "K"
	//SenderTaxID   = "3730207045443"
	//FromAccountNo = "0011232329"
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
