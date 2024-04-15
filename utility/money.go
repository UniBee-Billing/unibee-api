package utility

import (
	"fmt"
	"strconv"
	"strings"
)

func IsNoCentCurrency(currency string) bool {
	NOCentCurrencies := []string{"JPY", "KRW"}
	for _, element := range NOCentCurrencies {
		if currency == element {
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

func ConvertDollarStrToCent(dollarStr string, currency string) int64 {
	dollars, err := strconv.ParseFloat(dollarStr, 64)
	if err != nil {
		panic(fmt.Sprintf("ConvertDollarStrToCent panic:%s", dollarStr))
	}
	if IsNoCentCurrency(currency) {
		return int64(dollars)
	}
	cents := int64(dollars * 100)
	return cents
}

func ConvertTaxPercentageToPercentageString(taxPercentage int64) string {
	return fmt.Sprintf("%.1f", float64(taxPercentage)/100)
}

func ConvertTaxPercentageToPercentageFloat(taxPercentage int64) float64 {
	return float64(taxPercentage) / 100
}

func ConvertTaxPercentageToInternalFloat(taxPercentage int64) float64 {
	return float64(taxPercentage) / 10000
}
