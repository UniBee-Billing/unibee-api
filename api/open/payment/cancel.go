package payment

import "github.com/gogf/gf/v2/frame/g"

type CancelsReq struct {
	g.Meta           `path:"/cancels/{PaymentId}" tags:"Open-Payment-Controller" method:"post" summary:"Cancel Payment"`
	PaymentId        string `in:"path" dc:"PaymentId" v:"required"`
	MerchantId       int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	MerchantCancelId string `p:"merchantCancelId" dc:"MerchantCancelId" v:"required"`
}
type CancelsRes struct {
}
