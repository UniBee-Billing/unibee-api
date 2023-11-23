package v1

import "github.com/gogf/gf/v2/frame/g"

type CapturesReq struct {
	g.Meta               `path:"/captures/{PaymentsPspReference}" tags:"Out-Controller" method:"post" summary:"1.2 如果支付方式支持分布授权（请款）"`
	PaymentsPspReference string `in:"path" dc:"平台支付单号" v:"required|length:4,30#请输入平台支付单号长度为:{min}到:{max}位"`
}
type CapturesRes struct {
}
