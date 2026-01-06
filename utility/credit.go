package utility

import (
	"fmt"
	"unibee/internal/consts"
)

func ConvertCreditAmountToDollarStr(cents int64, currency string, AccountType int) string {
	if AccountType == consts.CreditAccountTypePromo {
		return fmt.Sprintf("%d", cents)
	} else {
		return ConvertCentToDollarStr(cents, currency)
	}
}
