package payment

import "github.com/gogf/gf/v2/frame/g"

type PaymentCallbackAgainReq struct {
	g.Meta    `path:"/payment_callback_again" tags:"System-Admin-Controller" method:"post" summary:"Admin Trigger Payment Callback"`
	PaymentId string `p:"paymentId" dc:"PaymentId" v:"required#请输入paymentId"`
}

type PaymentCallbackAgainRes struct {
}
