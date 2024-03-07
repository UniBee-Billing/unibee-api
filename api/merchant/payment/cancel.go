package payment

import "github.com/gogf/gf/v2/frame/g"

type CancelReq struct {
	g.Meta           `path:"/cancel" tags:"Payment" method:"post" summary:"Cancel Payment"`
	PaymentId        string `json:"paymentId" dc:"PaymentId" v:"required"`
	ExternalCancelId string `json:"externalCancelId" dc:"ExternalCancelId" v:"required"`
}
type CancelRes struct {
}
