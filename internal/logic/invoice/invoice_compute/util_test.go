package invoice_compute

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"testing"
	"unibee/api/bean"
	"unibee/test"
	"unibee/utility"
)

func TestInvoiceSimplifyCreation(t *testing.T) {
	t.Run("Test Creation", func(t *testing.T) {
		ctx := context.Background()
		invoice := ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &CalculateInvoiceReq{
			UserId:       0,
			InvoiceName:  "SubscriptionCreate",
			DiscountCode: "",
			TimeNow:      gtime.Now().Timestamp(),
			Currency:     "USD",
			PlanId:       test.TestPlan.Id,
			Quantity:     1,
			AddonJsonData: utility.MarshalToJsonString([]*bean.PlanAddonParam{
				{
					Quantity:    1,
					AddonPlanId: test.TestRecurringAddon.Id,
				},
			}),
			CountryCode:        "UAE",
			VatNumber:          "",
			TaxPercentage:      900,
			PeriodStart:        gtime.Now().Timestamp(),
			PeriodEnd:          gtime.Now().AddDate(0, 1, 0).Timestamp(),
			FinishTime:         gtime.Now().Timestamp(),
			ProductData:        nil,
			BillingCycleAnchor: gtime.Now().Timestamp(),
			Metadata:           nil,
		})
		VerifyInvoiceSimplify(invoice)
	})
}
