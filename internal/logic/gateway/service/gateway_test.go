package service

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/api/bean"
	"unibee/internal/query"
	"unibee/test"
)

func TestEditGateway(t *testing.T) {
	ctx := context.Background()
	t.Run("Wire Transfer", func(t *testing.T) {
		gatewayName := "wire transfer"
		one := query.GetGatewayByGatewayName(ctx, test.TestMerchant.Id, gatewayName)
		if one == nil {
			one = SetupWireTransferGateway(ctx, &WireTransferSetupReq{
				MerchantId:    test.TestMerchant.Id,
				Currency:      "USD",
				MinimumAmount: 100,
				Bank: &bean.GatewayBank{
					AccountHolder: "testAccountHolder",
					BIC:           "testBic",
					IBAN:          "testIBAN",
					Address:       "testAddress",
				},
			})
			require.NotNil(t, one)
			require.Equal(t, one.Currency, "USD")
			require.Equal(t, one.MinimumAmount, int64(100))
		}
		one = EditWireTransferGateway(ctx, &WireTransferSetupReq{
			GatewayId:     one.Id,
			MerchantId:    test.TestMerchant.Id,
			Currency:      "USD",
			MinimumAmount: 200,
			Bank: &bean.GatewayBank{
				AccountHolder: "testAccountHolder",
				BIC:           "testBic",
				IBAN:          "testIBAN",
				Address:       "testAddress",
			},
		})
		require.NotNil(t, one)
		require.Equal(t, one.Currency, "USD")
		require.Equal(t, one.MinimumAmount, int64(200))
	})
}
