package payment

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CapturesReq struct {
	g.Meta     `path:"/captures/{PaymentId}" tags:"Open-Payment-Controller" method:"post" summary:"1.2 如果支付方式支持分布授权（请款）"`
	PaymentId  string       `in:"path" dc:"平台支付单号" v:"required|length:4,30#请输入平台支付单号长度为:{min}到:{max}位"`
	MerchantId int64        `p:"merchantId" dc:"商户号" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
	Reference  string       `p:"reference" dc:"取消单号" v:"required"`
	Amount     *PayAmountVo `json:"amount"   in:"query" dc:"具体金额" v:"required"`
}
type CapturesRes struct {
}
