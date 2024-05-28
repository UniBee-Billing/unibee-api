package currency

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	_ "unibee/test"
)

func TestGetExchangeCurrencyMap(t *testing.T) {
	ctx := context.Background()
	t.Run("Test for currency rate", func(t *testing.T) {
		rate, err := GetExchangeConversionRates(ctx, "7dea9d6a5bafe83816a6ebdb", "USD", "CNY")
		require.Nil(t, err)
		require.Equal(t, *rate > 0, true)
	})
}
