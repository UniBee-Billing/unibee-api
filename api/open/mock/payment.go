package mock

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type SamplePaymentNetherlandsReq struct {
	g.Meta     `path:"/quick_pay_sample_netherlands" tags:"Open-Mock-Controller" method:"post" summary:"1.1 netherlands一键创建支付单(自动填充用户） https://docs.klarna.com/resources/test-environment/sample-customer-data/#netherlands"`
	Currency   string `p:"currency" dc:"currency 货币" d:"JPY" v:"required"`
	Amount     int64  `p:"amount" dc:"amount 金额(需x100，对比RMB到分）" v:"required"`
	MerchantId int64  `p:"merchantId" d:"15621" dc:"商户号" v:"required长度为:{min}到:{max}位"`
	Channel    string `p:"channel" d:"paypay" dc:"支付方式，klarna_paynow|klarna|klarna_account|paypay" v:"required"`
	ReturnUrl  string `p:"returnUrl" dc:"支付之后回跳商户Url" v:""`
	//PaymentBrandAddtion string `p:"paymentBrandAddtion" dc:"paymentBrandAddtion" v:""`
}
type SamplePaymentNetherlandsRes struct {
	Status    string      `p:"status" dc:"交易状态"`
	PaymentId string      `p:"paymentId" dc:"系统交易唯一编码-平台订单号"`
	Reference string      `p:"reference" dc:"商户订单号"`
	Action    *gjson.Json `p:"action" dc:"action"`
}

type DetailPayReq struct {
	g.Meta     `path:"/detail_pay" tags:"Open-Mock-Controller" method:"post" summary:"1.5支付单详情"`
	PaymentId  string `p:"paymentId" dc:"平台支付单号" v:"required"`
	MerchantId int64  `p:"merchantId" d:"15621" dc:"商户号" v:"required长度为:{min}到:{max}位"`
}
type DetailPayRes struct {
}
