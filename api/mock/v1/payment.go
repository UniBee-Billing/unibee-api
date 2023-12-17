package v1

import "github.com/gogf/gf/v2/frame/g"

type SamplePaymentNetherlandsReq struct {
	g.Meta              `path:"/quick_pay_sample_netherlands" tags:"Mock-Controller" method:"post" summary:"1.1 netherlands一键创建支付单(自动填充用户） https://docs.klarna.com/resources/test-environment/sample-customer-data/#netherlands"`
	Amount              int64  `p:"amount" dc:"amount 金额(需x100，对比RMB到分）" v:"required"`
	MerchantAccount     string `p:"merchantAccount" d:"15621" dc:"商户号" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
	PaymentMethod       string `p:"paymentMethod" d:"paypay" dc:"支付方式，klarna_paynow|klarna|klarna_account|paypay" v:"required"`
	Reference           string `p:"reference" dc:"客户单号" v:"required"`
	ReturnUrl           string `p:"returnUrl" dc:"支付之后回跳商户Url" v:""`
	PaymentBrandAddtion string `p:"paymentBrandAddtion" dc:"paymentBrandAddtion" v:""`
}
type SamplePaymentNetherlandsRes struct {
}
