package merchant

import (
	"unibee/utility"
)

func currencyNumberCheck(amount int64, currency string) {
	if utility.IsNoCentCurrency(currency) {
		utility.Assert(amount%100 == 0, "this currency No decimals allowed，made it divisible by 100")
	}
}
