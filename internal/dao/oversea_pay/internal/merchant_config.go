// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantConfigDao is the data access object for table merchant_config.
type MerchantConfigDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns MerchantConfigColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantConfigColumns defines and stores column names for table merchant_config.
type MerchantConfigColumns struct {
	Id          string // ID
	MerchantId  string // merchantId
	ConfigKey   string // config_key
	ConfigValue string // config_value
	GmtCreate   string // create time
	GmtModify   string // update time
	IsDeleted   string // 0-UnDeleted，1-Deleted
	CreateAt    string // create utc time
}

// merchantConfigColumns holds the columns for table merchant_config.
var merchantConfigColumns = MerchantConfigColumns{
	Id:          "id",
	MerchantId:  "merchant_id",
	ConfigKey:   "config_key",
	ConfigValue: "config_value",
	GmtCreate:   "gmt_create",
	GmtModify:   "gmt_modify",
	IsDeleted:   "is_deleted",
	CreateAt:    "create_at",
}

// NewMerchantConfigDao creates and returns a new DAO object for table data access.
func NewMerchantConfigDao() *MerchantConfigDao {
	return &MerchantConfigDao{
		group:   "oversea_pay",
		table:   "merchant_config",
		columns: merchantConfigColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantConfigDao) Columns() MerchantConfigColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantConfigDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
