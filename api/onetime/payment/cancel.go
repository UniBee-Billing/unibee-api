package payment

import "github.com/gogf/gf/v2/frame/g"

type CancelReq struct {
	g.Meta           `path:"/cancel/{PaymentId}" tags:"OneTime-Payment-Controller" method:"post" summary:"Cancel Payment"`
	PaymentId        string `in:"path" dc:"PaymentId" v:"required"`
	MerchantId       uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
	MerchantCancelId string `p:"merchantCancelId" dc:"MerchantCancelId" v:"required"`
}
type CancelRes struct {
}
