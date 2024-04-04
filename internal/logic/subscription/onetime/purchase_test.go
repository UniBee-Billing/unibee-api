package onetime

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/internal/query"
	_ "unibee/test"
)

func TestSubscription(t *testing.T) {
	ctx := context.Background()
	t.Run("Test case for subscription create preview with plan and addon, vat or not, vat number check", func(t *testing.T) {
		one := query.GetSubscriptionBySubscriptionId(ctx, "sub20240404Jus1gwVwBjE7O4C")
		addon, err := CreateSubscriptionOneTimeAddon(ctx, &SubscriptionCreateOnetimeAddonInternalReq{
			MerchantId:     one.MerchantId,
			SubscriptionId: one.SubscriptionId,
			AddonId:        114,
			Quantity:       1,
			RedirectUrl:    "",
			Metadata:       nil,
		})
		if err != nil {
			fmt.Printf("error:%s", err.Error())
		}
		require.Nil(t, err)
		require.NotNil(t, addon)
	})
}
