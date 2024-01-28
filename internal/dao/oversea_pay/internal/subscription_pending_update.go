// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SubscriptionPendingUpdateDao is the data access object for table subscription_pending_update.
type SubscriptionPendingUpdateDao struct {
	table   string                           // table is the underlying table name of the DAO.
	group   string                           // group is the database configuration group name of current DAO.
	columns SubscriptionPendingUpdateColumns // columns contains all the column names of Table for convenient usage.
}

// SubscriptionPendingUpdateColumns defines and stores column names for table subscription_pending_update.
type SubscriptionPendingUpdateColumns struct {
	Id                   string //
	MerchantId           string // 商户Id
	SubscriptionId       string // 订阅id（内部编号）
	UpdateSubscriptionId string // 升级单ID（内部编号）
	ChannelUpdateId      string // 支付渠道订阅更新单id， stripe 适用 channelInvoiceId对应
	GmtCreate            string // 创建时间
	Amount               string // 本周期金额,单位：分
	Status               string // 订阅单状态，0-Init | 1-Create｜2-Finished｜3-Cancelled
	ProrationAmount      string // 下周期金额,单位：分
	UpdateAmount         string // 升级到金额,单位：分
	Currency             string // 货币
	UpdateCurrency       string // 升级到货币
	PlanId               string // 计划ID
	UpdatePlanId         string // 升级到计划ID
	Quantity             string // quantity
	UpdateQuantity       string // 升级到quantity
	AddonData            string // plan addon json data
	UpdateAddonData      string // 升级到plan addon json data
	ChannelId            string // 支付渠道Id
	UserId               string // userId
	GmtModify            string // 修改时间
	IsDeleted            string //
	Paid                 string // 是否已支付，0-否，1-是
	Link                 string // 支付链接
	ChannelStatus        string // 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	MerchantUserId       string // merchant_user_id
	Data                 string // 渠道额外参数，JSON格式
	ResponseData         string // 渠道返回参数，JSON格式
	EffectImmediate      string // 是否马上生效，0-否，1-是
	EffectTime           string // effect_immediate=0, 预计生效时间 unit_time
	AdminNote            string // Admin 修改备注
	ProrationDate        string // merchant_user_id
}

// subscriptionPendingUpdateColumns holds the columns for table subscription_pending_update.
var subscriptionPendingUpdateColumns = SubscriptionPendingUpdateColumns{
	Id:                   "id",
	MerchantId:           "merchant_id",
	SubscriptionId:       "subscription_id",
	UpdateSubscriptionId: "update_subscription_id",
	ChannelUpdateId:      "channel_update_id",
	GmtCreate:            "gmt_create",
	Amount:               "amount",
	Status:               "status",
	ProrationAmount:      "proration_amount",
	UpdateAmount:         "update_amount",
	Currency:             "currency",
	UpdateCurrency:       "update_currency",
	PlanId:               "plan_id",
	UpdatePlanId:         "update_plan_id",
	Quantity:             "quantity",
	UpdateQuantity:       "update_quantity",
	AddonData:            "addon_data",
	UpdateAddonData:      "update_addon_data",
	ChannelId:            "channel_id",
	UserId:               "user_id",
	GmtModify:            "gmt_modify",
	IsDeleted:            "is_deleted",
	Paid:                 "paid",
	Link:                 "link",
	ChannelStatus:        "channel_status",
	MerchantUserId:       "merchant_user_id",
	Data:                 "data",
	ResponseData:         "response_data",
	EffectImmediate:      "effect_immediate",
	EffectTime:           "effect_time",
	AdminNote:            "admin_note",
	ProrationDate:        "proration_date",
}

// NewSubscriptionPendingUpdateDao creates and returns a new DAO object for table data access.
func NewSubscriptionPendingUpdateDao() *SubscriptionPendingUpdateDao {
	return &SubscriptionPendingUpdateDao{
		group:   "oversea_pay",
		table:   "subscription_pending_update",
		columns: subscriptionPendingUpdateColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SubscriptionPendingUpdateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SubscriptionPendingUpdateDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SubscriptionPendingUpdateDao) Columns() SubscriptionPendingUpdateColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SubscriptionPendingUpdateDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SubscriptionPendingUpdateDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SubscriptionPendingUpdateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
