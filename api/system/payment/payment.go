package payment

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type PaymentCallbackAgainReq struct {
	g.Meta    `path:"/payment_callback_again" tags:"System-Admin-Controller" method:"post" summary:"Admin Trigger Payment Callback"`
	PaymentId string `p:"paymentId" dc:"PaymentId" v:"required#请输入paymentId"`
}

type PaymentCallbackAgainRes struct {
}

type GatewayPaymentMethodListReq struct {
	g.Meta    `path:"/gateway_payment_method_list" tags:"System-Admin-Controller" method:"post" summary:"Admin Query Custom Gateway Payment Callback"`
	PaymentId string `p:"paymentId" dc:"PaymentId" v:"required#请输入paymentId"`
}

type GatewayPaymentMethodListRes struct {
	MethodList []*ro.PaymentMethod
}
