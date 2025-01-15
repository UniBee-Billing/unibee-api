package detail

import (
	"unibee/api/bean"
)

type PaymentDetail struct {
	User    *bean.UserAccount `json:"user" dc:"user"`
	Payment *bean.Payment     `json:"payment" dc:"Payment"`
	Gateway *Gateway          `json:"gateway" dc:"Gateway"`
	Invoice *InvoiceDetail    `json:"invoice" dc:"Invoice"`
}
