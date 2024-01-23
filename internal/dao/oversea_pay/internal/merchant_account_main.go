// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantAccountMainDao is the data access object for table merchant_account_main.
type MerchantAccountMainDao struct {
	table   string                     // table is the underlying table name of the DAO.
	group   string                     // group is the database configuration group name of current DAO.
	columns MerchantAccountMainColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantAccountMainColumns defines and stores column names for table merchant_account_main.
type MerchantAccountMainColumns struct {
	Id          string // 主键ID
	CompanyId   string //
	MerchantId  string // 商户ID
	Currency    string
	TotalTrade  string // 交易金额总计
	TotalRefund string // 退款金额总计
	TotalCut    string // 服务扣点金额总计
	TotalSend   string // 结算金额总计
	Year        string // 结算key-年
	Month       string // 结算key-月
	Day         string // 结算key-日
	GmtCreate   string //
	GmtModify   string //
	Statistic   string // 统计使用
	IsDeleted   string // 是否删除，0-未删除，1-已删除
}

// merchantAccountMainColumns holds the columns for table merchant_account_main.
var merchantAccountMainColumns = MerchantAccountMainColumns{
	Id:          "id",
	CompanyId:   "company_id",
	MerchantId:  "merchant_id",
	Currency:    "currency",
	TotalTrade:  "total_trade",
	TotalRefund: "total_refund",
	TotalCut:    "total_cut",
	TotalSend:   "total_send",
	Year:        "year",
	Month:       "month",
	Day:         "day",
	GmtCreate:   "gmt_create",
	GmtModify:   "gmt_modify",
	Statistic:   "statistic",
	IsDeleted:   "is_deleted",
}

// NewMerchantAccountMainDao creates and returns a new DAO object for table data access.
func NewMerchantAccountMainDao() *MerchantAccountMainDao {
	return &MerchantAccountMainDao{
		group:   "oversea_pay",
		table:   "merchant_account_main",
		columns: merchantAccountMainColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantAccountMainDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantAccountMainDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantAccountMainDao) Columns() MerchantAccountMainColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantAccountMainDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantAccountMainDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantAccountMainDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
