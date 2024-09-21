package payment

import (
	"github.com/gogf/gf/v2/frame/g"
)

type PaymentCallbackAgainReq struct {
	g.Meta    `path:"/payment_callback_again" tags:"System-Admin" method:"post" summary:"Admin Trigger Payment Callback"`
	PaymentId string `json:"paymentId" dc:"PaymentId" v:"required#Require paymentId"`
}

type PaymentCallbackAgainRes struct {
}
