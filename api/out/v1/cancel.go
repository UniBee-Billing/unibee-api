package v1

import "github.com/gogf/gf/v2/frame/g"

type CancelsReq struct {
	g.Meta               `path:"/cancels/{PaymentsPspReference}" tags:"Out-Controller" method:"post" summary:"1.3 取消⽀付单"`
	PaymentsPspReference string `in:"path" dc:"平台支付单号" v:"required|length:4,30#请输入平台支付单号长度为:{min}到:{max}位"`
}
type CancelsRes struct {
}
