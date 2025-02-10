package utility

import (
	"fmt"
	"strconv"
	"strings"
)

func IsNoCentCurrency(currency string) bool {
	NOCentCurrencies := []string{"JPY", "KRW"}
	for _, element := range NOCentCurrencies {
		if strings.ToUpper(currency) == element {
			return true
		}
	}
	return false
}

func ConvertCentToDollarStr(cents int64, currency string) string {
	if IsNoCentCurrency(strings.ToUpper(currency)) {
		return fmt.Sprintf("%d", cents)
	}
	dollars := float64(cents) / 100.0
	return strings.Replace(fmt.Sprintf("%.2f", dollars), ".00", "", -1)
}

func ExchangeCurrencyConvert(from int64, fromCurrency string, toCurrency string, exchangeRate float64) int64 {
	if IsNoCentCurrency(strings.ToUpper(fromCurrency)) && !IsNoCentCurrency(strings.ToUpper(toCurrency)) {
		return int64(float64(from*100) * exchangeRate)
	} else if !IsNoCentCurrency(strings.ToUpper(fromCurrency)) && IsNoCentCurrency(strings.ToUpper(toCurrency)) {
		return int64(float64(from/100) * exchangeRate)
	}
	return int64(float64(from) * exchangeRate)
}

func ConvertDollarStrToCent(dollarStr string, currency string) int64 {
	dollars, err := strconv.ParseFloat(dollarStr, 64)
	if err != nil {
		panic(fmt.Sprintf("ConvertDollarStrToCent panic dollarStr:%s currency:%s err:%s", dollarStr, currency, err.Error()))
	}
	if IsNoCentCurrency(currency) {
		return int64(dollars)
	}
	cents := int64(dollars * 100)
	return cents
}

func ConvertCentToDollarFloat(cents int64, currency string) float64 {
	if IsNoCentCurrency(strings.ToUpper(currency)) {
		return float64(cents)
	}
	dollars := float64(cents) / 100.0
	return dollars
}

func ConvertDollarFloatToInt64Cent(dollar float64, currency string) int64 {
	if IsNoCentCurrency(currency) {
		return int64(dollar)
	}
	cents := int64(dollar * 100)
	return cents
}

func ConvertTaxPercentageToPercentageString(taxPercentage int64) string {
	return fmt.Sprintf("%.1f", float64(taxPercentage)/100)
}

func ConvertTaxPercentageToPercentageFloat(taxPercentage int64) float64 {
	return float64(taxPercentage) / 100
}

func ConvertTaxPercentageToInternalFloat(taxPercentage int64) float64 {
	if taxPercentage == 0 {
		return 0
	}
	return float64(taxPercentage) / 10000
}
