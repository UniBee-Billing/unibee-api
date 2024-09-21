package detail

import "unibee/api/bean"

type RefundDetail struct {
	User    *bean.UserAccount `json:"user" dc:"user"`
	Payment *bean.Payment     `json:"payment" dc:"Payment"`
	Refund  *bean.Refund      `json:"refund" dc:"Refund"`
}
