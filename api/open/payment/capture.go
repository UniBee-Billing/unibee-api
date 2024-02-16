package payment

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CapturesReq struct {
	g.Meta            `path:"/captures/{PaymentId}" tags:"Open-Payment-Controller" method:"post" summary:"Capture Payment"`
	PaymentId         string    `in:"path" dc:"PaymentId" v:"required"`
	MerchantId        int64     `p:"merchantId" dc:"MerchantId" v:"required"`
	MerchantCaptureId string    `p:"merchantCaptureId" dc:"MerchantCaptureId" v:"required"`
	Amount            *AmountVo `json:"amount"   in:"query" dc:"Amount, Cent" v:"required"`
}
type CapturesRes struct {
}
