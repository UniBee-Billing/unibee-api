package payment

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CaptureReq struct {
	g.Meta            `path:"/capture" tags:"Payment" method:"post" summary:"Capture Payment"`
	PaymentId         string `json:"paymentId" dc:"PaymentId" v:"required"`
	ExternalCaptureId string `json:"externalCaptureId" dc:"ExternalCaptureId" v:"required"`
	CaptureAmount     int64  `json:"captureAmount" dc:"CaptureAmount, Cent" v:"required"`
	Currency          string `json:"currency" dc:"Currency"  v:"required"`
}
type CaptureRes struct {
}
