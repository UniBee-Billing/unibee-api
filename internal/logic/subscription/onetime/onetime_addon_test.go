package onetime

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/api/bean"
	"unibee/api/merchant/plan"
	"unibee/internal/consts"
	"unibee/internal/logic/plan/service"
	service2 "unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/test"
)

func TestSubscription(t *testing.T) {
	ctx := context.Background()
	var testSubscriptionId string
	t.Run("Test case for subscription create | onetime addon purchase", func(t *testing.T) {
		_, err := service.PlanAddonsBinding(ctx, &plan.AddonsBindingReq{
			PlanId:          test.TestPlan.Id,
			Action:          1,
			AddonIds:        []int64{int64(test.TestRecurringAddon.Id)},
			OnetimeAddonIds: []int64{int64(test.TestOneTimeAddon.Id)},
		})
		require.Nil(t, err)
		create, err := service2.SubscriptionCreate(ctx, &service2.CreateInternalReq{
			MerchantId:      test.TestMerchant.Id,
			PlanId:          test.TestPlan.Id,
			UserId:          test.TestUser.Id,
			Quantity:        1,
			GatewayId:       test.TestGateway.Id,
			PaymentMethodId: "testPaymentMethodId",
			AddonParams:     []*bean.PlanAddonParam{{Quantity: 1, AddonPlanId: test.TestRecurringAddon.Id}},
		})
		require.Nil(t, err)
		one := create.Subscription
		testSubscriptionId = one.SubscriptionId
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

	t.Run("Test for subscription cancel immediately", func(t *testing.T) {
		//cancel immediately
		err := service2.SubscriptionCancel(ctx, testSubscriptionId, false, false, "test cancel")
		require.Nil(t, err)
		one := query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.Status == consts.SubStatusCancelled)
	})
}
