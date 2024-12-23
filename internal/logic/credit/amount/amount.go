package amount

import (
	"context"
	"math"
	"unibee/internal/consts"
	"unibee/internal/query"
	"unibee/utility"
)

func ConvertCreditAmountToCurrency(ctx context.Context, merchantId uint64, creditType int, currency string, creditAmount int64) (currencyAmount int64, exchangeRate int64) {
	one := query.GetCreditConfig(ctx, merchantId, creditType, currency)
	if one == nil {
		return 0, 0
	}
	if one.Type == consts.CreditAccountTypePromo {
		return utility.ConvertDollarFloatToInt64Cent(float64(creditAmount)*(float64(one.ExchangeRate)/100), currency), one.ExchangeRate
	} else {
		return creditAmount, one.ExchangeRate
	}
}

func ConvertCurrencyAmountToCreditAmount(ctx context.Context, merchantId uint64, creditType int, currency string, currencyAmount int64) (creditAmount int64, exchangeRate int64) {
	one := query.GetCreditConfig(ctx, merchantId, creditType, currency)
	if one == nil {
		return 0, 0
	}
	if one.Type == consts.CreditAccountTypePromo {
		return int64(math.Ceil(utility.ConvertCentToDollarFloat(currencyAmount, currency) / (float64(one.ExchangeRate) / 100))), one.ExchangeRate
	} else {
		return currencyAmount, one.ExchangeRate
	}
}
