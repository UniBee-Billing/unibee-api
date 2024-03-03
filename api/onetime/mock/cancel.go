package mock

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CancelReq struct {
	g.Meta    `path:"/cancel" tags:"Open-Mock-Controller" method:"post" summary:"Mock Cancel Payment"`
	PaymentId string `json:"paymentId" dc:"PaymentId" v:"required"`
}
type CancelRes struct {
}
