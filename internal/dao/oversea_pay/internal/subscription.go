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
	SubscriptionId              string // 订阅id（内部编号）
	UserId                      string // userId
	GmtCreate                   string // 创建时间
	GmtModify                   string // 修改时间
	Amount                      string // 金额,单位：分
	Currency                    string // 货币
	MerchantId                  string // 商户Id
	PlanId                      string // 计划ID
	Quantity                    string // quantity
	AddonData                   string // plan addon json data
	LatestInvoiceId             string // latest_invoice_id
	Type                        string // sub type, 0-channel sub, 1-unibee sub
	ChannelId                   string // 支付渠道Id
	Status                      string // 订阅单状态，0-Init | 1-Create｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete
	ChannelSubscriptionId       string // 支付渠道订阅id
	CustomerName                string // customer_name
	CustomerEmail               string // customer_email
	IsDeleted                   string // 0-UnDeleted，1-Deleted
	Link                        string //
	ChannelStatus               string // 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	ChannelItemData             string // channel_item_data
	CancelAtPeriodEnd           string // 是否在周期结束时取消，0-false | 1-true
	ChannelLatestInvoiceId      string // 渠道最新发票 id
	LastUpdateTime              string //
	CurrentPeriodStart          string // current_period_start
	CurrentPeriodEnd            string // current_period_end
	CurrentPeriodStartTime      string //
	CurrentPeriodEndTime        string //
	BillingCycleAnchor          string // billing_cycle_anchor
	TrialEnd                    string // trial_end
	ReturnUrl                   string //
	FirstPayTime                string // 首次支付时间
	CancelReason                string //
	CountryCode                 string //
	VatNumber                   string //
	TaxScale                    string // Tax税率，万分位，1000 表示 10%
	VatVerifyData               string //
	Data                        string // 渠道额外参数，JSON格式
	ResponseData                string // 渠道返回参数，JSON格式
	PendingUpdateId             string //
	ChannelDefaultPaymentMethod string //
}

// subscriptionColumns holds the columns for table subscription.
var subscriptionColumns = SubscriptionColumns{
	Id:                          "id",
	SubscriptionId:              "subscription_id",
	UserId:                      "user_id",
	GmtCreate:                   "gmt_create",
	GmtModify:                   "gmt_modify",
	Amount:                      "amount",
	Currency:                    "currency",
	MerchantId:                  "merchant_id",
	PlanId:                      "plan_id",
	Quantity:                    "quantity",
	AddonData:                   "addon_data",
	LatestInvoiceId:             "latest_invoice_id",
	Type:                        "type",
	ChannelId:                   "channel_id",
	Status:                      "status",
	ChannelSubscriptionId:       "channel_subscription_id",
	CustomerName:                "customer_name",
	CustomerEmail:               "customer_email",
	IsDeleted:                   "is_deleted",
	Link:                        "link",
	ChannelStatus:               "channel_status",
	ChannelItemData:             "channel_item_data",
	CancelAtPeriodEnd:           "cancel_at_period_end",
	ChannelLatestInvoiceId:      "channel_latest_invoice_id",
	LastUpdateTime:              "last_update_time",
	CurrentPeriodStart:          "current_period_start",
	CurrentPeriodEnd:            "current_period_end",
	CurrentPeriodStartTime:      "current_period_start_time",
	CurrentPeriodEndTime:        "current_period_end_time",
	BillingCycleAnchor:          "billing_cycle_anchor",
	TrialEnd:                    "trial_end",
	ReturnUrl:                   "return_url",
	FirstPayTime:                "first_pay_time",
	CancelReason:                "cancel_reason",
	CountryCode:                 "country_code",
	VatNumber:                   "vat_number",
	TaxScale:                    "tax_scale",
	VatVerifyData:               "vat_verify_data",
	Data:                        "data",
	ResponseData:                "response_data",
	PendingUpdateId:             "pendingUpdateId",
	ChannelDefaultPaymentMethod: "channel_default_payment_method",
}

// NewSubscriptionDao creates and returns a new DAO object for table data access.
func NewSubscriptionDao() *SubscriptionDao {
	return &SubscriptionDao{
		group:   "oversea_pay",
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
