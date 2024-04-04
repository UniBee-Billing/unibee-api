package method

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/internal/consts"
	"unibee/test"
)

func TestPayment(t *testing.T) {
	ctx := context.Background()
	_ = test.TestGateway
	t.Run("Test for QueryList", func(t *testing.T) {
		list := QueryPaymentMethodList(ctx, &PaymentMethodListInternalReq{
			MerchantId: consts.CloudModeManagerMerchantId,
			UserId:     200365887,
			GatewayId:  25,
		})
		require.NotNil(t, list)
		require.Equal(t, true, len(list) > 0)
		one := QueryPaymentMethod(ctx, consts.CloudModeManagerMerchantId, 200365887, 25, list[0].Id)
		require.NotNil(t, one)
		require.Equal(t, one.Id, list[0].Id)
		url, one := NewPaymentMethod(ctx, &NewPaymentMethodInternalReq{
			MerchantId: consts.CloudModeManagerMerchantId,
			UserId:     200365887,
			GatewayId:  25,
			Currency:   "USD",
		})
		require.NotNil(t, url)
		require.Nil(t, one)
	})
}
