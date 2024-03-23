package vat_gateway

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/test"
)

func TestVat(t *testing.T) {
	ctx := context.Background()
	var err error
	t.Run("Test for vat interface api", func(t *testing.T) {
		one := GetDefaultVatGateway(ctx, test.TestMerchant.Id)
		require.Nil(t, one)
		_, err = ValidateVatNumberByDefaultGateway(ctx, test.TestMerchant.Id, test.TestUser.Id, "", "")
		require.NotNil(t, err)
		err = InitMerchantDefaultVatGateway(ctx, test.TestMerchant.Id)
		require.NotNil(t, err)
		_, err = MerchantCountryRateList(ctx, test.TestMerchant.Id)
		require.NotNil(t, err)
		_, err = QueryVatCountryRateByMerchant(ctx, test.TestMerchant.Id, "CN")
		require.NotNil(t, err)
		err = SetupMerchantVatConfig(ctx, test.TestMerchant.Id, "github", "github", true)
		require.Nil(t, err)
		res, err := ValidateVatNumberByDefaultGateway(ctx, test.TestMerchant.Id, test.TestUser.Id, "IE6388047V", "")
		require.Nil(t, err)
		require.NotNil(t, res)
		require.Equal(t, true, res.Valid)
		res, err = ValidateVatNumberByDefaultGateway(ctx, test.TestMerchant.Id, test.TestUser.Id, "IE6388047V"+uuid.New().String(), "")
		require.NotNil(t, err)
		require.Nil(t, res)
		err = InitMerchantDefaultVatGateway(ctx, test.TestMerchant.Id)
		require.Nil(t, err)
		_, err = MerchantCountryRateList(ctx, test.TestMerchant.Id)
		require.Nil(t, err)
		_, err = QueryVatCountryRateByMerchant(ctx, test.TestMerchant.Id, "NL")
		require.Nil(t, err)
	})
	t.Run("Test for vat config clean", func(t *testing.T) {
		require.Nil(t, CleanMerchantDefaultVatConfig(ctx, test.TestMerchant.Id))
	})
}
