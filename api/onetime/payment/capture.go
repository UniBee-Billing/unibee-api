package payment

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CaptureReq struct {
	g.Meta            `path:"/capture/{PaymentId}" tags:"OneTime-Payment" method:"post" summary:"Capture Payment"`
	PaymentId         string `in:"path" dc:"PaymentId" v:"required"`
	MerchantCaptureId string `json:"merchantCaptureId" dc:"MerchantCaptureId" v:"required"`
	CaptureAmount     int64  `json:"captureAmount"   in:"query" dc:"CaptureAmount, Cent" v:"required"`
	Currency          string `json:"currency"   in:"query" dc:"Currency"  v:"required"`
}
type CaptureRes struct {
}
