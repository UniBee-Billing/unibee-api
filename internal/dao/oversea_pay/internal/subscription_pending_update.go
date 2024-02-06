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
	Id                   string // id
	MerchantId           string // merchant id
	SubscriptionId       string // subscription id
	UpdateSubscriptionId string // pending update unique id
	GatewayUpdateId      string // gateway update payment id assosiate to this update, use payment.paymentId
	GmtCreate            string // create time
	Amount               string // amount of this period, cent
	Status               string // status，0-Init | 1-Create｜2-Finished｜3-Cancelled
	ProrationAmount      string // proration amount of this pending update , cent
	UpdateAmount         string // the amount after update
	Currency             string // currency of this period
	UpdateCurrency       string // the currency after update
	PlanId               string // the plan id of this period
	UpdatePlanId         string // the plan id after update
	Quantity             string // quantity of this period
	UpdateQuantity       string // quantity after update
	AddonData            string // plan addon data (json) of this period
	UpdateAddonData      string // plan addon data (json) after update
	GatewayId            string // gateway_id
	UserId               string // userId
	GmtModify            string // update time
	IsDeleted            string // 0-UnDeleted，1-Deleted
	Paid                 string // paid，0-no，1-yes
	Link                 string // payment link
	GatewayStatus        string // gateway status
	MerchantUserId       string // merchant_user_id
	Data                 string //
	ResponseData         string //
	EffectImmediate      string // effect immediate，0-no，1-yes
	EffectTime           string // effect_immediate=0, effect time, utc_time
	Note                 string // note
	ProrationDate        string // merchant_user_id
	CreateAt             string // create utc time
}

// subscriptionPendingUpdateColumns holds the columns for table subscription_pending_update.
var subscriptionPendingUpdateColumns = SubscriptionPendingUpdateColumns{
	Id:                   "id",
	MerchantId:           "merchant_id",
	SubscriptionId:       "subscription_id",
	UpdateSubscriptionId: "update_subscription_id",
	GatewayUpdateId:      "gateway_update_id",
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
	GatewayId:            "gateway_id",
	UserId:               "user_id",
	GmtModify:            "gmt_modify",
	IsDeleted:            "is_deleted",
	Paid:                 "paid",
	Link:                 "link",
	GatewayStatus:        "gateway_status",
	MerchantUserId:       "merchant_user_id",
	Data:                 "data",
	ResponseData:         "response_data",
	EffectImmediate:      "effect_immediate",
	EffectTime:           "effect_time",
	Note:                 "note",
	ProrationDate:        "proration_date",
	CreateAt:             "create_at",
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
