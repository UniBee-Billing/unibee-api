package payment

import "github.com/gogf/gf/v2/frame/g"

type CancelsReq struct {
	g.Meta     `path:"/cancels/{PaymentId}" tags:"Open-Payment-Controller" method:"post" summary:"1.3 取消⽀付单"`
	PaymentId  string `in:"path" dc:"平台支付单号" v:"required"`
	MerchantId int64  `p:"merchantId" dc:"商户号" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
	Reference  string `p:"reference" dc:"取消单号" v:"required"`
}
type CancelsRes struct {
}
