package cycle

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/api/bean"
	"unibee/internal/consts"
	service2 "unibee/internal/logic/invoice/detail"
	"unibee/internal/logic/invoice/invoice_compute"
	service3 "unibee/internal/logic/invoice/service"
	"unibee/internal/logic/subscription/config"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/test"
	"unibee/utility/unibee"
)

// success testcases
// case: billing cycle without pendingUpdate and check dunning time invoice
// case: billing cycle with pendingUpdate and check dunning time invoice
// case: set subscription trialEnd and billing cycle effected, check trialEnd radius, should after max(now,periodEnd) -- todo set time not may cause sub new cycle invoice and payment
// case: upgrade|downgrade subscription after periodEnd and before trialEnd

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
			err := service.SubscriptionCancel(ctx, one.SubscriptionId, false, false, "test cancel")
			require.Nil(t, err)
		}
		preview, err := service.SubscriptionCreatePreview(ctx, &service.CreatePreviewInternalReq{
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
		invoice_compute.VerifyInvoiceSimplify(preview.Invoice)
		require.Equal(t, true, preview.TotalAmount == (test.TestPlan.Amount*testQuantity)+test.TestRecurringAddon.Amount*testQuantity)
		require.Equal(t, true, preview.Currency == test.TestPlan.Currency)
		require.Equal(t, true, len(preview.Gateways) > 0)
		err = vat_gateway.SetupMerchantVatConfig(ctx, test.TestMerchant.Id, "github", "github", true)
		require.Nil(t, err)
		preview, err = service.SubscriptionCreatePreview(ctx, &service.CreatePreviewInternalReq{
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
		invoice_compute.VerifyInvoiceSimplify(preview.Invoice)
		require.Equal(t, true, preview.TotalAmount == preview.Invoice.TotalAmountExcludingTax)
		require.Equal(t, true, preview.TotalAmount == preview.Invoice.TotalAmountExcludingTax+preview.Invoice.TaxAmount)
		require.Equal(t, true, preview.Invoice.TotalAmountExcludingTax == ((test.TestPlan.Amount*testQuantity)+(test.TestRecurringAddon.Amount*testQuantity)))
		require.Equal(t, true, preview.Currency == test.TestPlan.Currency)

		preview, err = service.SubscriptionCreatePreview(ctx, &service.CreatePreviewInternalReq{
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
		invoice_compute.VerifyInvoiceSimplify(preview.Invoice)
		require.Equal(t, true, preview.TotalAmount == preview.Invoice.TotalAmountExcludingTax)
		require.Equal(t, true, preview.Invoice.TotalAmountExcludingTax == ((test.TestPlan.Amount*testQuantity)+(test.TestRecurringAddon.Amount*testQuantity)))
		require.Equal(t, true, preview.Currency == test.TestPlan.Currency)
	})
	t.Run("Test for vat config clean", func(t *testing.T) {
		require.Nil(t, vat_gateway.CleanMerchantDefaultVatConfig(ctx, test.TestMerchant.Id))
	})
	t.Run("Test for subscription create|cancelAtPeriodEnd|billing cycle effected", func(t *testing.T) {
		create, err := service.SubscriptionCreate(ctx, &service.CreateInternalReq{
			MerchantId:      test.TestMerchant.Id,
			PlanId:          test.TestPlan.Id,
			UserId:          test.TestUser.Id,
			Quantity:        testQuantity,
			GatewayId:       test.TestGateway.Id,
			PaymentMethodId: "testPaymentMethodId",
			AddonParams:     []*bean.PlanAddonParam{{Quantity: testQuantity, AddonPlanId: test.TestRecurringAddon.Id}},
			Discount: &bean.ExternalDiscountParam{
				Recurring:          unibee.Bool(false),
				DiscountAmount:     nil,
				DiscountPercentage: unibee.Int64(100),
				Metadata:           map[string]interface{}{"freeTraffic": 5},
			},
		})
		require.Nil(t, err)
		require.NotNil(t, create)
		require.NotNil(t, create.Subscription)
		require.NotNil(t, create.Link)
		require.NotNil(t, create.Paid)
		require.Equal(t, true, create.Paid)
		testSubscriptionId = create.Subscription.SubscriptionId
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.Status == consts.SubStatusActive)
		invoice := query.GetInvoiceByInvoiceId(ctx, one.LatestInvoiceId)
		require.NotNil(t, invoice)
		require.Equal(t, true, invoice.Status == consts.InvoiceStatusPaid)
		invoiceDetail := service2.InvoiceDetail(ctx, one.LatestInvoiceId)
		require.NotNil(t, invoiceDetail)
		err = CycleWalkForSubTest(ctx, testSubscriptionId, one.CurrentPeriodEnd-config.GetMerchantSubscriptionConfig(ctx, test.TestMerchant.Id).TryAutomaticPaymentBeforePeriodEnd-1, "testcase")
		require.Nil(t, err)
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, len(one.LatestInvoiceId) > 0)
		require.Equal(t, true, one.Status == consts.SubStatusActive)
		invoice = query.GetInvoiceByInvoiceId(ctx, one.LatestInvoiceId)
		require.NotNil(t, invoice)
		require.Equal(t, true, invoice.Status == consts.InvoiceStatusProcessing)
		err = CycleWalkForSubTest(ctx, testSubscriptionId, one.CurrentPeriodEnd+1, "testcase")
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, len(one.LatestInvoiceId) > 0)
		require.Equal(t, true, one.Status == consts.SubStatusActive)
		invoice = query.GetInvoiceByInvoiceId(ctx, one.LatestInvoiceId)
		require.NotNil(t, invoice)
		require.Equal(t, true, invoice.Status == consts.InvoiceStatusPaid)
		//start test cancelAtPeriodEnd
		err = service.SubscriptionCancelAtPeriodEnd(ctx, testSubscriptionId, false, 0)
		require.Nil(t, err)
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.CancelAtPeriodEnd == 1)
		err = service.SubscriptionCancelLastCancelAtPeriodEnd(ctx, testSubscriptionId, false)
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.CancelAtPeriodEnd == 0)
		err = service.SubscriptionCancelAtPeriodEnd(ctx, testSubscriptionId, false, 0)
		require.Nil(t, err)
		err = CycleWalkForSubTest(ctx, testSubscriptionId, one.CurrentPeriodEnd+1, "testcase")
		require.Nil(t, err)
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.Status == consts.SubStatusCancelled)
	})
	t.Run("Test for subscription cancel immediately", func(t *testing.T) {
		//cancel immediately
		err := service.SubscriptionCancel(ctx, testSubscriptionId, false, false, "test cancel")
		require.Nil(t, err)
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.Status == consts.SubStatusCancelled)
	})
	t.Run("Test for subscription trialEnd", func(t *testing.T) {
		create, err := service.SubscriptionCreate(ctx, &service.CreateInternalReq{
			MerchantId:      test.TestMerchant.Id,
			PlanId:          test.TestPlan.Id,
			UserId:          test.TestUser.Id,
			Quantity:        testQuantity,
			GatewayId:       test.TestGateway.Id,
			PaymentMethodId: "testPaymentMethodId",
			AddonParams:     []*bean.PlanAddonParam{{Quantity: testQuantity, AddonPlanId: test.TestRecurringAddon.Id}},
			TrialEnd:        gtime.Now().Timestamp() + 86400,
		})
		require.Nil(t, err)
		require.NotNil(t, create)
		require.NotNil(t, create.Subscription)
		require.NotNil(t, create.Link)
		require.NotNil(t, create.Paid)
		require.Equal(t, true, create.Paid)
		testSubscriptionId = create.Subscription.SubscriptionId
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.Equal(t, one.Status, consts.SubStatusActive)
		//cancel immediately
		err = service.SubscriptionCancel(ctx, testSubscriptionId, false, false, "test cancel")
		require.Nil(t, err)
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.Status == consts.SubStatusCancelled)
	})
	t.Run("Test for subscription upgrade|downgrade", func(t *testing.T) {
		create, err := service.SubscriptionCreate(ctx, &service.CreateInternalReq{
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
		require.Equal(t, true, create.Paid)
		testSubscriptionId = create.Subscription.SubscriptionId
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.Status == consts.SubStatusActive)
		//upgrade
		preview, err := service.SubscriptionUpdatePreview(ctx, &service.UpdatePreviewInternalReq{
			SubscriptionId: testSubscriptionId,
			NewPlanId:      test.TestPlan.Id,
			Quantity:       2,
			GatewayId:      one.GatewayId,
		}, 0, 0)
		require.Nil(t, err)
		invoice_compute.VerifyInvoiceSimplify(preview.Invoice)
		invoice_compute.VerifyInvoiceSimplify(preview.NextPeriodInvoice)
		_, err = service.SubscriptionUpdate(ctx, &service.UpdateInternalReq{
			SubscriptionId: testSubscriptionId,
			NewPlanId:      test.TestPlan.Id,
			Quantity:       3, //todo mark if set to 2 will cause a bug
			GatewayId:      one.GatewayId,
		}, 0)
		require.Nil(t, err)
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, int64(3), one.Quantity)

		//err = service.SubscriptionCancelAtPeriodEnd(ctx, testSubscriptionId, false, 0)
		//require.Nil(t, err)
		//one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		//require.NotNil(t, one)
		//require.Equal(t, true, one.CancelAtPeriodEnd == 1)

		_, err = service.SubscriptionUpdate(ctx, &service.UpdateInternalReq{
			SubscriptionId: testSubscriptionId,
			NewPlanId:      test.TestPlan.Id,
			Quantity:       1,
			GatewayId:      one.GatewayId,
		}, 0)
		require.Nil(t, err)
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, int64(3), one.Quantity)
		//require.Equal(t, true, one.CancelAtPeriodEnd == 0)
		err = CycleWalkForSubTest(ctx, testSubscriptionId, one.CurrentPeriodEnd+1, "testcase")
		require.Nil(t, err)
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.Status == consts.SubStatusActive)
		require.Equal(t, int64(1), one.Quantity)
	})
	t.Run("Test for subscription cancel immediately", func(t *testing.T) {
		//cancel immediately
		err := service.SubscriptionCancel(ctx, testSubscriptionId, false, false, "test cancel")
		require.Nil(t, err)
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.Status == consts.SubStatusCancelled)
	})
	t.Run("Test for wire transfer subscription upgrade|downgrade", func(t *testing.T) {
		create, err := service.SubscriptionCreate(ctx, &service.CreateInternalReq{
			MerchantId:      test.TestMerchant.Id,
			PlanId:          test.TestPlan.Id,
			UserId:          test.TestUser.Id,
			Quantity:        testQuantity,
			GatewayId:       test.TestWireTransferGateway.Id,
			PaymentMethodId: "testPaymentMethodId",
			AddonParams:     []*bean.PlanAddonParam{{Quantity: testQuantity, AddonPlanId: test.TestRecurringAddon.Id}},
		})
		require.Nil(t, err)
		require.NotNil(t, create)
		require.NotNil(t, create.Subscription)
		require.NotNil(t, create.Link)
		require.NotNil(t, create.Paid)
		require.Equal(t, false, create.Paid)
		testSubscriptionId = create.Subscription.SubscriptionId
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.Status == consts.SubStatusPending)
		require.NotNil(t, one.LatestInvoiceId)
		_, err = service3.MarkWireTransferInvoiceAsSuccess(ctx, one.LatestInvoiceId, "automatic_transfer_number")
		require.Nil(t, err)
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.Status == consts.SubStatusActive)
	})
	t.Run("Test for wire transfer subscription cancel immediately", func(t *testing.T) {
		//cancel immediately
		err := service.SubscriptionCancel(ctx, testSubscriptionId, false, false, "test cancel")
		require.Nil(t, err)
		one = query.GetSubscriptionBySubscriptionId(ctx, testSubscriptionId)
		require.NotNil(t, one)
		require.Equal(t, true, one.Status == consts.SubStatusCancelled)
	})
}

func CycleWalkForSubTest(ctx context.Context, subscriptionId string, time int64, source string) error {
	for {
		result, err := SubPipeBillingCycleWalk(ctx, subscriptionId, time, source)
		if err != nil {
			return err
		} else {
			if !result.WalkUnfinished {
				return nil
			}
		}
	}
}
