package notify

import "unibee/internal/logic/gateway/api/alipay/api/model"

type AlipayPayResultNotify struct {
	AlipayNotify
	PaymentRequestId         string                   `json:"paymentRequestId,omitempty"`
	PaymentId                string                   `json:"paymentId,omitempty"`
	PaymentAmount            *model.Amount            `json:"paymentAmount,omitempty"`
	PaymentCreateTime        string                   `json:"paymentCreateTime,omitempty"`
	PaymentTime              string                   `json:"paymentTime,omitempty"`
	CustomsDeclarationAmount *model.Amount            `json:"customsDeclarationAmount,omitempty"`
	GrossSettlementAmount    *model.Amount            `json:"grossSettlementAmount,omitempty"`
	SettlementQuote          *model.Quote             `json:"settlementQuote,omitempty"`
	PspCustomerInfo          *model.PspCustomerInfo   `json:"pspCustomerInfo,omitempty"`
	AcquirerReferenceNo      string                   `json:"acquirerReferenceNo,omitempty"`
	PaymentResultInfo        *model.PaymentResultInfo `json:"paymentResultInfo,omitempty"`
	AcquirerInfo             *model.AcquirerInfo      `json:"acquirerInfo,omitempty"`
	PromotionResult          []*model.PromotionResult `json:"promotionResult,omitempty"`
	PaymentMethodType        string                   `json:"paymentMethodType,omitempty"`
}
