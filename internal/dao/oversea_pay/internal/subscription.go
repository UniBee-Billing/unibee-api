// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SubscriptionDao is the data access object for table subscription.
type SubscriptionDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns SubscriptionColumns // columns contains all the column names of Table for convenient usage.
}

// SubscriptionColumns defines and stores column names for table subscription.
type SubscriptionColumns struct {
	Id                          string //
	SubscriptionId              string // subscription id
	UserId                      string // userId
	GmtCreate                   string // create time
	GmtModify                   string // update time
	TaskTime                    string // task_time
	Amount                      string // amount, cent
	Currency                    string // currency
	MerchantId                  string // merchant id
	PlanId                      string // plan id
	Quantity                    string // quantity
	AddonData                   string // plan addon json data
	LatestInvoiceId             string // latest_invoice_id
	Type                        string // sub type, 0-gateway sub, 1-unibee sub
	GatewayId                   string // gateway_id
	Status                      string // status，0-Init | 1-Pending｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete
	GatewaySubscriptionId       string // gateway subscription id
	CustomerName                string // customer_name
	CustomerEmail               string // customer_email
	IsDeleted                   string // 0-UnDeleted，1-Deleted
	GatewayDefaultPaymentMethod string //
	Link                        string //
	GatewayStatus               string // gateway status，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	GatewayItemData             string // gateway_item_data
	CancelAtPeriodEnd           string // whether cancel at period end，0-false | 1-true
	DiscountCode                string // discount_code
	LastUpdateTime              string //
	CurrentPeriodStart          string // current_period_start, utc time
	CurrentPeriodEnd            string // current_period_end, utc time
	CurrentPeriodStartTime      string //
	CurrentPeriodEndTime        string //
	BillingCycleAnchor          string // billing_cycle_anchor
	DunningTime                 string // dunning_time, utc time
	TrialEnd                    string // trial_end, utc time
	ReturnUrl                   string //
	FirstPaidTime               string // first success payment time
	CancelReason                string //
	CountryCode                 string //
	VatNumber                   string //
	TaxPercentage               string // taxPercentage，1000 = 10%
	VatVerifyData               string //
	Data                        string //
	ResponseData                string //
	PendingUpdateId             string //
	CreateTime                  string // create utc time
	TestClock                   string // test_clock, simulator clock for subscription, if set , sub will out of cronjob controll
	MetaData                    string // meta_data(json)
	GasPayer                    string // who pay the gas, merchant|user
}

// subscriptionColumns holds the columns for table subscription.
var subscriptionColumns = SubscriptionColumns{
	Id:                          "id",
	SubscriptionId:              "subscription_id",
	UserId:                      "user_id",
	GmtCreate:                   "gmt_create",
	GmtModify:                   "gmt_modify",
	TaskTime:                    "task_time",
	Amount:                      "amount",
	Currency:                    "currency",
	MerchantId:                  "merchant_id",
	PlanId:                      "plan_id",
	Quantity:                    "quantity",
	AddonData:                   "addon_data",
	LatestInvoiceId:             "latest_invoice_id",
	Type:                        "type",
	GatewayId:                   "gateway_id",
	Status:                      "status",
	GatewaySubscriptionId:       "gateway_subscription_id",
	CustomerName:                "customer_name",
	CustomerEmail:               "customer_email",
	IsDeleted:                   "is_deleted",
	GatewayDefaultPaymentMethod: "gateway_default_payment_method",
	Link:                        "link",
	GatewayStatus:               "gateway_status",
	GatewayItemData:             "gateway_item_data",
	CancelAtPeriodEnd:           "cancel_at_period_end",
	DiscountCode:                "discount_code",
	LastUpdateTime:              "last_update_time",
	CurrentPeriodStart:          "current_period_start",
	CurrentPeriodEnd:            "current_period_end",
	CurrentPeriodStartTime:      "current_period_start_time",
	CurrentPeriodEndTime:        "current_period_end_time",
	BillingCycleAnchor:          "billing_cycle_anchor",
	DunningTime:                 "dunning_time",
	TrialEnd:                    "trial_end",
	ReturnUrl:                   "return_url",
	FirstPaidTime:               "first_paid_time",
	CancelReason:                "cancel_reason",
	CountryCode:                 "country_code",
	VatNumber:                   "vat_number",
	TaxPercentage:               "tax_percentage",
	VatVerifyData:               "vat_verify_data",
	Data:                        "data",
	ResponseData:                "response_data",
	PendingUpdateId:             "pendingUpdateId",
	CreateTime:                  "create_time",
	TestClock:                   "test_clock",
	MetaData:                    "meta_data",
	GasPayer:                    "gas_payer",
}

// NewSubscriptionDao creates and returns a new DAO object for table data access.
func NewSubscriptionDao() *SubscriptionDao {
	return &SubscriptionDao{
		group:   "default",
		table:   "subscription",
		columns: subscriptionColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SubscriptionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SubscriptionDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SubscriptionDao) Columns() SubscriptionColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SubscriptionDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SubscriptionDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SubscriptionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
