package v1

import "github.com/gogf/gf/v2/frame/g"

type PaymentsReq struct {
	g.Meta `path:"/payments" tags:"Out-Controller" method:"post" summary:"1.1 用于接收商户段的请求（包括 token 请求）"`
}
type PaymentsRes struct {
}

type PaymentMethodsReq struct {
	g.Meta `path:"/paymentMethods" tags:"Out-Controller" method:"post" summary:"1.0 根据配置⽀付⽅式的信息，通过请求字段筛选可以返回的⽀付⽅式(Klarna、Evonet支持）"`
}
type PaymentMethodsRes struct {
}

type PaymentDetailsReq struct {
	g.Meta               `path:"/paymentDetails/{PaymentsPspReference}" tags:"Out-Controller" method:"post" summary:"1.5 查询当前交易状态及详情"`
	PaymentsPspReference string `in:"path" dc:"平台支付单号" v:"required|length:4,30#请输入平台支付单号长度为:{min}到:{max}位"`
}
type PaymentDetailsRes struct {
}
