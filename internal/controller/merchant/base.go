package merchant

import (
	"strings"
	"unibee/utility"
)

func currencyNumberCheck(amount int64, currency string) {
	if strings.Compare(currency, "JPY") == 0 {
		utility.Assert(amount%100 == 0, "this currency No decimals allowed，made it divisible by 100")
	}
}
