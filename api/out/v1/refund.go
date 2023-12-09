package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/api/out/vo"
)

type RefundsReq struct {
	g.Meta               `path:"/refunds/{PaymentsPspReference}" tags:"Out-Controller" method:"post" summary:"1.4 退款，请款之后"`
	PaymentsPspReference string          `in:"path" dc:"平台支付单号" v:"required|length:4,30#请输入平台支付单号长度为:{min}到:{max}位"`
	MerchantAccount      string          `p:"merchantAccount" dc:"商户号" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
	Reference            string          `p:"reference" dc:"取消单号" v:"required"`
	Amount               *vo.PayAmountVo `json:"amount"   in:"query" dc:"具体金额" v:"required"`
}
type RefundsRes struct {
}
