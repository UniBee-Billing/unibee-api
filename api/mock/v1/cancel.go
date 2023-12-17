package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CancelReq struct {
	g.Meta              `path:"/cancel" tags:"Mock-Controller" method:"post" summary:"1.2取消支付单"`
	PaymentPspReference string `p:"paymentPspReference" dc:"平台支付单号" v:"required"`
	MerchantId          int64  `p:"merchantAccount" d:"15621" dc:"商户号" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
}
type CancelRes struct {
}
