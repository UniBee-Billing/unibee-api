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
	Id                    string //
	GmtCreate             string // 创建时间
	GmtModify             string // 修改时间
	CompanyId             string // 公司ID
	MerchantId            string // 商户Id
	PlanId                string // 计划ID
	ChannelId             string // 支付渠道Id
	UserId                string // userId
	Quantity              string // quantity
	SubscriptionId        string // 内部订阅id
	ChannelSubscriptionId string // 支付渠道订阅id
	Data                  string // 渠道额外参数，JSON格式
	ResponseData          string // 渠道返回参数，JSON格式
	IsDeleted             string //
	Status                string // 订阅单状态，0-Init | 1-Create｜2-Active｜3-Inactive
	ChannelUserId         string // 渠道用户 Id
	CustomerName          string // customer_name
	CustomerEmail         string // customer_email
	Link                  string //
}

// subscriptionColumns holds the columns for table subscription.
var subscriptionColumns = SubscriptionColumns{
	Id:                    "id",
	GmtCreate:             "gmt_create",
	GmtModify:             "gmt_modify",
	CompanyId:             "company_id",
	MerchantId:            "merchant_id",
	PlanId:                "plan_id",
	ChannelId:             "channel_id",
	UserId:                "user_id",
	Quantity:              "quantity",
	SubscriptionId:        "subscription_id",
	ChannelSubscriptionId: "channel_subscription_id",
	Data:                  "data",
	ResponseData:          "response_data",
	IsDeleted:             "is_deleted",
	Status:                "status",
	ChannelUserId:         "channel_user_id",
	CustomerName:          "customer_name",
	CustomerEmail:         "customer_email",
	Link:                  "link",
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
