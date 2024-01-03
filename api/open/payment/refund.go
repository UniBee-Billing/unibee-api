package payment

import (
	"github.com/gogf/gf/v2/frame/g"
)

type RefundsReq struct {
	g.Meta               `path:"/refunds/{PaymentsPspReference}" tags:"Open-Payment-Controller" method:"post" summary:"1.4 退款，请款之后"`
	PaymentsPspReference string       `in:"path" dc:"平台支付单号" v:"required|length:4,30#请输入平台支付单号长度为:{min}到:{max}位"`
	MerchantId           int64        `p:"merchantId" dc:"商户号" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
	Reference            string       `p:"reference" dc:"取消单号" v:"required"`
	Amount               *PayAmountVo `json:"amount"   in:"query" dc:"具体金额" v:"required"`
}
type RefundsRes struct {
	Status              string `p:"status" dc:"交易状态"`
	PspReference        string `p:"pspReference" dc:"系统交易唯一编码-平台订单号"`
	Reference           string `p:"reference" dc:"商户订单号"`
	PaymentPspReference string `p:"paymentPspReference" dc:"系统交易唯一编码-平台支付订单号"`
}
