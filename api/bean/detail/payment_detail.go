package detail

import (
	"unibee/api/bean"
)

type PaymentDetail struct {
	User    *bean.UserAccountSimplify `json:"user" dc:"user"`
	Payment *bean.PaymentSimplify     `json:"payment" dc:"Payment"`
	//Invoice *InvoiceDetail            `json:"invoice" dc:"Invoice"`
}
