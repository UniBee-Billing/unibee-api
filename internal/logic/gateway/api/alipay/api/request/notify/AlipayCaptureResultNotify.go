package notify

import "unibee/internal/logic/gateway/api/alipay/api/model"

type AlipayCaptureResultNotify struct {
	AlipayNotify
	CaptureRequestId    string              `json:"captureRequestId,omitempty"`
	PaymentId           string              `json:"paymentId,omitempty"`
	CaptureId           string              `json:"captureId,omitempty"`
	CaptureAmount       *model.Amount       `json:"captureAmount,omitempty"`
	CaptureTime         string              `json:"captureTime,omitempty"`
	AcquirerReferenceNo string              `json:"acquirerReferenceNo,omitempty"`
	AcquirerInfo        *model.AcquirerInfo `json:"acquirerInfo,omitempty"`
}
