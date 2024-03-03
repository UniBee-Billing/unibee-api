package payment

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CaptureReq struct {
	g.Meta            `path:"/capture/{PaymentId}" tags:"OneTime-Payment-Controller" method:"post" summary:"Capture Payment"`
	PaymentId         string    `in:"path" dc:"PaymentId" v:"required"`
	MerchantCaptureId string    `json:"merchantCaptureId" dc:"MerchantCaptureId" v:"required"`
	Amount            *AmountVo `json:"amount"   in:"query" dc:"Amount, Cent" v:"required"`
}
type CaptureRes struct {
}
