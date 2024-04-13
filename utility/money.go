package utility

import (
	"fmt"
	"strconv"
	"strings"
)

func ConvertCentToDollarStr(cents int64, currency string) string {
	if strings.Compare(strings.ToUpper(currency), "JPY") == 0 {
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
	if strings.Compare(strings.ToUpper(currency), "JPY") == 0 {
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
