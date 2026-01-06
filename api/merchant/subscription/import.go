package subscription

import (
	"unibee/api/bean/detail"

	"github.com/gogf/gf/v2/frame/g"
)

type ActiveSubscriptionImportReq struct {
	g.Meta                 `path:"/active_subscription_import" tags:"Subscription Import" method:"post" summary:"Active Subscription Import(allows repetition imports)" dc:"Each repetition import overrides existing subscriptions with matching ExternalSubscriptionId."`
	ExternalSubscriptionId string                 `json:"externalSubscriptionId"    dc:"Required, The external id of subscription" required:"true"   `
	PlanId                 uint64                 `json:"planId" dc:"The id of plan, one of planId or ExternalPlanId is required, plan should created at UniBee at first "`
	ExternalPlanId         string                 `json:"externalPlanId"   dc:"The external id of plan, one of planId or ExternalPlanId is required, plan should created at UniBee at first"   `
	Email                  string                 `json:"email"  dc:"The email of user, one of Email or ExternalUserId is required" `
	ExternalUserId         string                 `json:"externalUserId"    dc:"The external id of user, one of Email or ExternalUserId is required "    `
	Quantity               int64                  `json:"quantity"      dc:"the quantity of plan, default 1 if not provided " required:"true"        `
	CountryCode            string                 `json:"countryCode"    dc:"Required. Specifies the ISO 3166-1 alpha-2 country code for the subscription (e.g., EE, RU). This code determines the applicable tax rules for the subscription." required:"true"  `
	VatNumber              string                 `json:"vatNumber"    dc:"The Vat Number of user"  `
	TaxPercentage          int64                  `json:"taxPercentage" dc:"The tax percentage. Only applicable when the system VAT gateway not setup. Value is in thousandths (e.g., 1000 = 10%)."`
	Gateway                string                 `json:"gateway" dc:"Required, should one of stripe|paypal|wire_transfer|changelly " required:"true"           `
	CurrentPeriodStart     string                 `json:"currentPeriodStart" dc:"Required, UTC time, the current period start time of subscription, format '2006-01-02 15:04:05'" required:"true"`
	CurrentPeriodEnd       string                 `json:"currentPeriodEnd"   dc:"Required, UTC time, the current period end time of subscription, format '2006-01-02 15:04:05'" required:"true"`
	BillingCycleAnchor     string                 `json:"billingCycleAnchor"   dc:"Required, UTC time, The reference point that aligns future billing cycle dates. It sets the day of week for week intervals, the day of month for month and year intervals, and the month of year for year intervals, format '2006-01-02 15:04:05'" required:"true"`
	FirstPaidTime          string                 `json:"firstPaidTime"   dc:"UTC time, the first payment success time of subscription, format '2006-01-02 15:04:05'"   `
	CreateTime             string                 `json:"createTime"      dc:"Required, UTC time, the creation time of subscription, format '2006-01-02 15:04:05'" required:"true"   `
	Features               string                 `json:"features"    dc:"In json format, additional features data of subscription, will join user's metric data in user api if provided'"     `
	ExpectedTotalAmount    int64                  `json:"expectedTotalAmount" dc:"Optional. Unit: cents. If greater than 0, the system will verify the calculated total amount against this value"`
	Metadata               map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}
type ActiveSubscriptionImportRes struct {
	Subscription *detail.SubscriptionDetail `json:"subscription" dc:"Subscription"`
}

type HistorySubscriptionImportReq struct {
	g.Meta                 `path:"/history_subscription_import" tags:"Subscription Import" method:"post" summary:"History Subscription Import(Allows repetition imports)" dc:"Each repetition import overrides existing subscriptions with matching ExternalSubscriptionId."`
	ExternalSubscriptionId string                 `json:"externalSubscriptionId"    dc:"Required, The external id of subscription" required:"true"   `
	PlanId                 uint64                 `json:"planId" dc:"The id of plan, one of planId or ExternalPlanId is required, plan should created at UniBee at first "`
	ExternalPlanId         string                 `json:"externalPlanId"   dc:"The external id of plan, one of planId or ExternalPlanId is required, plan should created at UniBee at first"   `
	Email                  string                 `json:"email"  dc:"The email of user, one of Email or ExternalUserId is required" `
	ExternalUserId         string                 `json:"externalUserId"    dc:"The external id of user, one of Email or ExternalUserId is required "    `
	Quantity               int64                  `json:"quantity"      dc:"the quantity of plan, default 1 if not provided " required:"true"        `
	CountryCode            string                 `json:"countryCode"    dc:"Required. Specifies the ISO 3166-1 alpha-2 country code for the subscription (e.g., EE, RU). This code determines the applicable tax rules for the subscription." required:"true"  `
	TaxPercentage          int64                  `json:"taxPercentage" dc:"The TaxPercentage of subscription, Only applicable when the system VAT gateway not setup, 1000 = 10%"`
	Gateway                string                 `json:"gateway" dc:"Required, should one of stripe|paypal|wire_transfer|changelly " required:"true"           `
	CurrentPeriodStart     string                 `json:"currentPeriodStart" dc:"Required, UTC time, the current period start time of subscription, format '2006-01-02 15:04:05'" required:"true"`
	CurrentPeriodEnd       string                 `json:"currentPeriodEnd"   dc:"Required, UTC time, the current period end time of subscription, format '2006-01-02 15:04:05'" required:"true"`
	TotalAmount            int64                  `json:"totalAmount" dc:"Required. Unit: cents." required:"true"`
	Metadata               map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}
type HistorySubscriptionImportRes struct {
	Subscription *detail.SubscriptionDetail `json:"subscription" dc:"Subscription"`
}
