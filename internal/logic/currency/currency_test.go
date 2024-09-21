package currency

import (
	"github.com/stretchr/testify/require"
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
