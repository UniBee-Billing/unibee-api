package payment

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CaptureReq struct {
	g.Meta            `path:"/capture/{PaymentId}" tags:"OneTime-Payment-Controller" method:"post" summary:"Capture Payment"`
	PaymentId         string    `in:"path" dc:"PaymentId" v:"required"`
	MerchantId        uint64    `p:"merchantId" dc:"MerchantId" v:"required"`
	MerchantCaptureId string    `p:"merchantCaptureId" dc:"MerchantCaptureId" v:"required"`
	Amount            *AmountVo `json:"amount"   in:"query" dc:"Amount, Cent" v:"required"`
}
type CaptureRes struct {
}
