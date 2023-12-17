package v1

import "github.com/gogf/gf/v2/frame/g"

type CaptureReq struct {
	g.Meta              `path:"/capture" tags:"Mock-Controller" method:"post" summary:"1.3捕获支付单"`
	PaymentPspReference string `p:"paymentPspReference" dc:"平台支付单号" v:"required"`
	MerchantId          int64  `p:"merchantAccount" d:"15621" dc:"商户号" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
	Currency            string `p:"currency" dc:"currency 货币" d:"JPY" v:"required"`
	Amount              int64  `p:"amount" dc:"amount 金额(需x100，对比RMB到分）" v:"required"`
}
type CaptureRes struct {
}
