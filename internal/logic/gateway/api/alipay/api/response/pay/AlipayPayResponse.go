package responsePay

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/response"
)

type AlipayPayResponse struct {
	response.AlipayResponse
	PaymentRequestId        string                    `json:"paymentRequestId"`
	PaymentId               string                    `json:"paymentId"`
	PaymentAmount           model.Amount              `json:"paymentAmount"`
	PaymentData             string                    `json:"paymentData"`
	ActualPaymentAmount     model.Amount              `json:"actualPaymentAmount"`
	PaymentQuote            model.Quote               `json:"paymentQuote"`
	PaymentTime             string                    `json:"paymentTime"`
	PaymentCreateTime       string                    `json:"paymentCreateTime"`
	AuthExpiryTime          string                    `json:"authExpiryTime"`
	NonGuaranteeCouponValue model.Amount              `json:"nonGuaranteeCouponValue"`
	PaymentActionForm       string                    `json:"paymentActionForm"`
	PspCustomerInfo         model.PspCustomerInfo     `json:"pspCustomerInfo"`
	ChallengeActionForm     model.ChallengeActionForm `json:"challengeActionForm"`
	RedirectActionForm      model.RedirectActionForm  `json:"redirectActionForm"`
	OrderCodeForm           model.OrderCodeForm       `json:"orderCodeForm"`
	GrossSettlementAmount   model.Amount              `json:"grossSettlementAmount"`
	SettlementQuote         model.Quote               `json:"settlementQuote"`
	ExtendInfo              string                    `json:"extendInfo"`
	NormalUrl               string                    `json:"normalUrl"`
	SchemeUrl               string                    `json:"schemeUrl"`
	ApplinkUrl              string                    `json:"applinkUrl"`
	AppIdentifier           string                    `json:"appIdentifier"`
	PaymentResultInfo       model.PaymentResultInfo   `json:"paymentResultInfo"`
	AcquirerInfo            model.AcquirerInfo        `json:"acquirerInfo"`
	PromotionResult         []model.PromotionResult   `json:"promotionResult"`
}
