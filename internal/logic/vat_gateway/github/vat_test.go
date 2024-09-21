package vat

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGithubVat(t *testing.T) {
	t.Run("Test for github vat", func(t *testing.T) {
		one := &Github{
			Password: "",
			Name:     "github",
		}
		require.Equal(t, "github", one.GetGatewayName())
		list, err := one.ListAllCountries()
		require.Nil(t, err)
		require.NotNil(t, list)
		require.Equal(t, 0, len(list))
		rates, err := one.ListAllRates()
		require.Nil(t, err)
		require.NotNil(t, list)
		require.Equal(t, true, len(rates) > 0)
		number, err := one.ValidateVatNumber("IE6388047V", "")
		require.Nil(t, err)
		require.NotNil(t, number)
		require.Equal(t, true, number.Valid)
		number, err = one.ValidateVatNumber("NL123456789B01", "")
		require.NotNil(t, err)
		require.Nil(t, number)
		number, err = one.ValidateEoriNumber("NL123456789B01")
		require.NotNil(t, err)
		require.Nil(t, number)
	})
}
