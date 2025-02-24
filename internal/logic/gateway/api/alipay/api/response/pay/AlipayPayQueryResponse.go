package responsePay

import (
	"time"
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/response"
)

type AlipayPayQueryResponse struct {
	response.AlipayResponse
	PaymentStatus            model.TransactionStatusType `json:"paymentStatus"`
	PaymentResultCode        string                      `json:"paymentResultCode"`
	PaymentResultMessage     string                      `json:"paymentResultMessage"`
	PaymentRequestId         string                      `json:"paymentRequestId"`
	PaymentId                string                      `json:"paymentId"`
	AuthPaymentId            string                      `json:"authPaymentId"`
	PaymentAmount            model.Amount                `json:"paymentAmount"`
	ActualPaymentAmount      model.Amount                `json:"actualPaymentAmount"`
	PaymentQuote             model.Quote                 `json:"paymentQuote"`
	AuthExpiryTime           string                      `json:"authExpiryTime"`
	PaymentCreateTime        time.Time                   `json:"paymentCreateTime"`
	PaymentTime              time.Time                   `json:"paymentTime"`
	NonGuaranteeCouponAmount model.Amount                `json:"nonGuaranteeCouponAmount"`
	PspCustomerInfo          model.PspCustomerInfo       `json:"pspCustomerInfo"`
	RedirectActionForm       model.RedirectActionForm    `json:"redirectActionForm"`
	CardInfo                 model.CardInfo              `json:"cardInfo"`
	AcquirerReferenceNo      string                      `json:"acquirerReferenceNo"`
	ExtendInfo               string                      `json:"extendInfo"`
	Transactions             []model.Transaction         `json:"transactions"`
	CustomsDeclarationAmount model.Amount                `json:"customsDeclarationAmount"`
	GrossSettlementAmount    model.Amount                `json:"grossSettlementAmount"`
	SettlementQuote          model.Quote                 `json:"settlementQuote"`
	PaymentResultInfo        model.PaymentResultInfo     `json:"paymentResultInfo"`
	AcquirerInfo             model.AcquirerInfo          `json:"acquirerInfo"`
	MerchantAccountId        string                      `json:"merchantAccountId"`
	PromotionResults         []model.PromotionResult     `json:"promotionResult"`
	EarliestSettlementTime   string                      `json:"earliestSettlementTime"`
	PaymentMethodType        string                      `json:"paymentMethodType,omitempty"`
}
