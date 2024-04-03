package payment

import (
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type PaymentDetail struct {
	User    *bean.UserAccountSimplify `json:"user" dc:"user"`
	Payment *bean.PaymentSimplify     `json:"payment" dc:"Payment"`
	Invoice *detail.InvoiceDetail     `json:"invoice" dc:"Invoice"`
}
