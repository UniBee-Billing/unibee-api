package service

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/test"
)

// success testcases
// case: set cancelAtPeriodEnd subscription and billing cycle effected, and check upgrade|downgrade will resume it
// case: upgrade subscription with addon
// case: billing cycle without pendingUpdate and check dunning time invoice
// case: downgrade subscription with addon
// case: billing cycle with pendingUpdate and check dunning time invoice
// case: set subscription trialEnd and billing cycle effected, check trialEnd radius, should after max(now,periodEnd) -- todo set time not may cause sub new cycle invoice and payment
// case: upgrade|downgrade subscription after periodEnd and before trialEnd
// case: cancel subscription immediately

// failure testcases
// case1: create subscription with payment failure and check expired cycle
// case2: billing cycle with payment failure after periodEnd, sub should change to incomplete, else may set trialEnd
// case3: incomplete status situations todo

func TestSubscription(t *testing.T) {
	ctx := context.Background()
	var testQuantity int64 = 1
	var testSubscriptionId string
	var one *entity.Subscription
	t.Run("Test case for subscription create preview with plan and addon, vat or not, vat number check", func(t *testing.T) {
		one = query.GetLatestActiveOrIncompleteOrCreateSubscriptionByUserId(ctx, test.TestUser.Id, test.TestMerchant.Id)
		if one != nil {
			err := SubscriptionCancel(ctx, one.SubscriptionId, false, false, "test cancel")
			require.Nil(t, err)
		}
		preview, err := SubscriptionCreatePreview(ctx, &CreatePreviewInternalReq{
			MerchantId:  test.TestMerchant.Id,
			PlanId:      test.TestPlan.Id,
			UserId:      test.TestUser.Id,
			Quantity:    testQuantity,
			GatewayId:   test.TestGateway.Id,
			AddonParams: []*bean.PlanAddonParam{{Quantity: testQuantity, AddonPlanId: test.TestRecurringAddon.Id}},
		})
		require.Nil(t, err)
		require.Nil(t, preview.VatNumberValidate)
		require.NotNil(t, preview)
		require.NotNil(t, preview.Gateways)
		require.NotNil(t, preview.Invoice)
		require.Equal(t, true, preview.TotalAmount == (test.TestPlan.Amount*testQuantity)+test.TestRecurringAddon.Amount*testQuantity)
		require.Equal(t, true, preview.Currency == test.TestPlan.Currency)
		require.Equal(t, true, len(preview.Gateways) > 0)
		err = vat_gateway.SetupMerchantVatConfig(ctx, test.TestMerchant.Id, "github", "github", true)
		require.Nil(t, err)
		preview, err = SubscriptionCreatePreview(ctx, &CreatePreviewInternalReq{
			MerchantId:     test.TestMerchant.Id,
			PlanId:         test.TestPlan.Id,
			UserId:         test.TestUser.Id,
			Quantity:       testQuantity,
			GatewayId:      test.TestGateway.Id,
			AddonParams:    []*bean.PlanAddonParam{{Quantity: testQuantity, AddonPlanId: test.TestRecurringAddon.Id}},
			VatCountryCode: "AT",
		})
		require.Nil(t, err)
		require.Nil(t, preview.VatNumberValidate)
		require.Equal(t, true, preview.TotalAmount > preview.Invoice.TotalAmountExcludingTax)
		require.Equal(t, true, preview.TotalAmount == preview.Invoice.TotalAmountExcludingTax+preview.Invoice.TaxAmount)
		require.Equal(t, true, preview.Invoice.TotalAmountExcludingTax == ((test.TestPlan.Amount*testQuantity)+(test.TestRecurringAddon.Amount*testQuantity)))
		require.Equal(t, true, preview.Currency == test.TestPlan.Currency)

		preview, err = SubscriptionCreatePreview(ctx, &CreatePreviewInternalReq{
			MerchantId:     test.TestMerchant.Id,
			PlanId:         test.TestPlan.Id,
			UserId:         test.TestUser.Id,
			Quantity:       testQuantity,
			GatewayId:      test.TestGateway.Id,
			AddonParams:    []*bean.PlanAddonParam{{Quantity: testQuantity, AddonPlanId: test.TestRecurringAddon.Id}},
			VatCountryCode: "AT",
			VatNumber:      "IE6388047V",
		})
		require.Nil(t, err)
		require.NotNil(t, preview.VatNumberValidate)
		require.Equal(t, true, preview.TotalAmount == preview.Invoice.TotalAmountExcludingTax)
		require.Equal(t, true, preview.Invoice.TotalAmountExcludingTax == ((test.TestPlan.Amount*testQuantity)+(test.TestRecurringAddon.Amount*testQuantity)))
		require.Equal(t, true, preview.Currency == test.TestPlan.Currency)
	})
	t.Run("Test for vat config clean", func(t *testing.T) {
		require.Nil(t, vat_gateway.CleanMerchantDefaultVatConfig(ctx, test.TestMerchant.Id))
	})
	t.Run("Test for subscription create|cancelAtPeriodEnd|billing cycle effected|upgrade|downgrade|resume cancelAtPeriodEnd", func(t *testing.T) {
		create, err := SubscriptionCreate(ctx, &CreateInternalReq{
			MerchantId:      test.TestMerchant.Id,
			PlanId:          test.TestPlan.Id,
			UserId:          test.TestUser.Id,
			Quantity:        testQuantity,
			GatewayId:       test.TestGateway.Id,
			PaymentMethodId: "testPaymentMethodId",
			AddonParams:     []*bean.PlanAddonParam{{Quantity: testQuantity, AddonPlanId: test.TestRecurringAddon.Id}},
		})
		require.Nil(t, err)
		require.NotNil(t, create)
		require.NotNil(t, create.Subscription)
		require.NotNil(t, create.Link)
		require.NotNil(t, create.Paid)
		testSubscriptionId = create.Subscription.SubscriptionId
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		err = SubscriptionCancelAtPeriodEnd(ctx, testSubscriptionId, false, 0)
		require.Nil(t, err)
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.CancelAtPeriodEnd == 1)
		err = SubscriptionCancelLastCancelAtPeriodEnd(ctx, testSubscriptionId, false)
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.CancelAtPeriodEnd == 0)

	})
	t.Run("Test for subscription cancel", func(t *testing.T) {
		err := SubscriptionCancel(ctx, testSubscriptionId, false, false, "test cancel")
		require.Nil(t, err)
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.Status == consts.SubStatusCancelled)
	})
}
