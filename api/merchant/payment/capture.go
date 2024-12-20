package payment

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CaptureReq struct {
	g.Meta            `path:"/capture" tags:"Payment" method:"post" summary:"CapturePayment"`
	PaymentId         string `json:"paymentId" dc:"The unique id of payment" v:"required"`
	ExternalCaptureId string `json:"externalCaptureId" dc:"The external id of payment capture" v:"required"`
}
type CaptureRes struct {
}
