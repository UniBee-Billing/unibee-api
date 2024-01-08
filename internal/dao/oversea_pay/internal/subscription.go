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
	Id                     string //
	SubscriptionId         string // 订阅id（内部编号）
	UpdateSubscriptionId   string // 升级来源订阅 ID（内部编号）
	GmtCreate              string // 创建时间
	Amount                 string // 金额,单位：分
	Currency               string // 货币
	MerchantId             string // 商户Id
	PlanId                 string // 计划ID
	Quantity               string // quantity
	AddonData              string // plan addon json data
	ChannelId              string // 支付渠道Id
	Status                 string // 订阅单状态，0-Init | 1-Create｜2-Active｜3-Suspend | 4-Cancel | 5-Expire
	UserId                 string // userId
	ChannelSubscriptionId  string // 支付渠道订阅id
	Data                   string // 渠道额外参数，JSON格式
	ResponseData           string // 渠道返回参数，JSON格式
	ChannelUserId          string // 渠道用户 Id
	CustomerName           string // customer_name
	CustomerEmail          string // customer_email
	GmtModify              string // 修改时间
	IsDeleted              string //
	Link                   string //
	ChannelStatus          string // 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	ChannelLatestInvoiceId string // 渠道最新发票 id
	CancelAtPeriodEnd      string // 是否在周期结束时取消，0-false | 1-true
	CurrentPeriodStart     string // current_period_start
	CurrentPeriodEnd       string // current_period_end
	TrailEnd               string // trail_end
	ReturnUrl              string //
}

// subscriptionColumns holds the columns for table subscription.
var subscriptionColumns = SubscriptionColumns{
	Id:                     "id",
	SubscriptionId:         "subscription_id",
	UpdateSubscriptionId:   "update_subscription_id",
	GmtCreate:              "gmt_create",
	Amount:                 "amount",
	Currency:               "currency",
	MerchantId:             "merchant_id",
	PlanId:                 "plan_id",
	Quantity:               "quantity",
	AddonData:              "addon_data",
	ChannelId:              "channel_id",
	Status:                 "status",
	UserId:                 "user_id",
	ChannelSubscriptionId:  "channel_subscription_id",
	Data:                   "data",
	ResponseData:           "response_data",
	ChannelUserId:          "channel_user_id",
	CustomerName:           "customer_name",
	CustomerEmail:          "customer_email",
	GmtModify:              "gmt_modify",
	IsDeleted:              "is_deleted",
	Link:                   "link",
	ChannelStatus:          "channel_status",
	ChannelLatestInvoiceId: "channel_latest_invoice_id",
	CancelAtPeriodEnd:      "cancel_at_period_end",
	CurrentPeriodStart:     "current_period_start",
	CurrentPeriodEnd:       "current_period_end",
	TrailEnd:               "trail_end",
	ReturnUrl:              "return_url",
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
