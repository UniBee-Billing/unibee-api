package currency

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/currency"
	"strings"
	"testing"
)

func TestGetMerchantCurrencyMap(t *testing.T) {
	t.Run("Test for currency", func(t *testing.T) {
		list := GetMerchantCurrencies()
		require.NotNil(t, list)
		require.Greater(t, len(list), 1)
		currencyMap := GetMerchantCurrencyMap()
		require.NotNil(t, currencyMap)
		require.Equal(t, true, IsFiatCurrencySupport("USD"))
		require.Equal(t, true, IsFiatCurrencySupport("usd"))
		require.Equal(t, true, IsFiatCurrencySupport("EUR"))
		require.Equal(t, true, IsFiatCurrencySupport("eur"))
	})
}

func TestCurrencySymbol(t *testing.T) {
	t.Run("Test for symbol", func(t *testing.T) {
		list := GetMerchantCurrencies()
		require.NotNil(t, list)
		require.Greater(t, len(list), 1)
		for _, one := range GetMerchantCurrencies() {
			g.Log().Infof(context.Background(), "%s 窄符号:%s,完整符号:%s", one.Currency, fmt.Sprintf("%s", currency.NarrowSymbol(currency.MustParseISO(strings.ToUpper(one.Currency)))), one.Symbol)
		}
	})
}
