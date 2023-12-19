// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SubscriptionPlanDao is the data access object for table subscription_plan.
type SubscriptionPlanDao struct {
	table   string                  // table is the underlying table name of the DAO.
	group   string                  // group is the database configuration group name of current DAO.
	columns SubscriptionPlanColumns // columns contains all the column names of Table for convenient usage.
}

// SubscriptionPlanColumns defines and stores column names for table subscription_plan.
type SubscriptionPlanColumns struct {
	Id           string //
	GmtCreate    string // 创建时间
	GmtModify    string // 修改时间
	CompanyId    string // 公司ID
	MerchantId   string // 商户Id
	PlanName     string // 计划名称
	Amount       string // 金额,单位：分
	Currency     string // 货币
	IntervalUnit string // 周期,day|month|year|week
	Description  string //
	IsDeleted    string //
}

// subscriptionPlanColumns holds the columns for table subscription_plan.
var subscriptionPlanColumns = SubscriptionPlanColumns{
	Id:           "id",
	GmtCreate:    "gmt_create",
	GmtModify:    "gmt_modify",
	CompanyId:    "company_id",
	MerchantId:   "merchant_id",
	PlanName:     "plan_name",
	Amount:       "amount",
	Currency:     "currency",
	IntervalUnit: "interval_unit",
	Description:  "description",
	IsDeleted:    "is_deleted",
}

// NewSubscriptionPlanDao creates and returns a new DAO object for table data access.
func NewSubscriptionPlanDao() *SubscriptionPlanDao {
	return &SubscriptionPlanDao{
		group:   "oversea_pay",
		table:   "subscription_plan",
		columns: subscriptionPlanColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SubscriptionPlanDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SubscriptionPlanDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SubscriptionPlanDao) Columns() SubscriptionPlanColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SubscriptionPlanDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SubscriptionPlanDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SubscriptionPlanDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
