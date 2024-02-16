package mock

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CancelReq struct {
	g.Meta     `path:"/cancel" tags:"Open-Mock-Controller" method:"post" summary:"Mock Cancel Payment"`
	PaymentId  string `p:"paymentId" dc:"PaymentId" v:"required"`
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required"`
}
type CancelRes struct {
}
