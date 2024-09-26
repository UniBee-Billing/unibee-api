package vat_gateway

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	gateway2 "unibee/api/bean"
	"unibee/test"
	"unibee/utility"
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

	t.Run("Test MLX Gateway Vat Config", func(t *testing.T) {
		var gatewayVATRules = make([]*gateway2.MerchantVatRule, 0)
		gatewayVATRules = append(gatewayVATRules, &gateway2.MerchantVatRule{
			GatewayNames:      "stripe",
			ValidCountryCodes: "AT,BE,BG,CY,CZ,DE,DK,EE,ES,FI,FR,GR,HR,HU,IE,IT,LT,LU,LV,MT,NL,PL,PT,RO,SE,SI,SK,GB,AE",
		})
		gatewayVATRules = append(gatewayVATRules, &gateway2.MerchantVatRule{
			GatewayNames:      "*",
			ValidCountryCodes: "AT,BE,BG,CY,CZ,DE,DK,EE,ES,FI,FR,GR,HR,HU,IE,IT,LT,LU,LV,MT,NL,PL,PT,RO,SE,SI,SK,GB",
		})
		fmt.Println(utility.MarshalToJsonString(gatewayVATRules))
	})
}
