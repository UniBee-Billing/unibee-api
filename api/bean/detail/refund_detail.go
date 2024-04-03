package detail

import "unibee/api/bean"

type RefundDetail struct {
	User    *bean.UserAccountSimplify `json:"user" dc:"user"`
	Payment *bean.PaymentSimplify     `json:"payment" dc:"Payment"`
	Refund  *bean.RefundSimplify      `json:"refund" dc:"Refund"`
}
