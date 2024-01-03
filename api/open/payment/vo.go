package payment

type PayAmountVo struct {
	Currency string `json:"currency"   in:"query" dc:"币种"  v:"required"`
	Value    int64  `json:"value"   in:"query" dc:"金额，单位分"  v:"required"`
}
