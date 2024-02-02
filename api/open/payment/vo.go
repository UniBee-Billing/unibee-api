package payment

type AmountVo struct {
	Currency string `json:"currency"   in:"query" dc:"Currency"  v:"required"`
	Amount   int64  `json:"amount"   in:"query" dc:"Amount，Cent"  v:"required"`
}
