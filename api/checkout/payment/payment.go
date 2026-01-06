package payment

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type DetailReq struct {
	g.Meta    `path:"/detail" tags:"Checkout" method:"get" summary:"Payment Detail"`
	PaymentId string `json:"paymentId" dc:"The unique id of payment" v:"required"`
}
type DetailRes struct {
	//PaymentDetail *detail.PaymentDetail `json:"paymentDetail" dc:"Payment Detail Object"`
	PaymentStatus int           `json:"paymentStatus" dc:"Payment Status，10-pending，20-success，30-failure, 40-cancel"`
	Payment       *bean.Payment `json:"payment" dc:"Payment"`
	ReturnUrl     string        `json:"returnUrl"`
	CancelUrl     string        `json:"cancelUrl"`
}
