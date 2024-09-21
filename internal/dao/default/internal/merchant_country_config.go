// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantCountryConfigDao is the data access object for table merchant_country_config.
type MerchantCountryConfigDao struct {
	table   string                       // table is the underlying table name of the DAO.
	group   string                       // group is the database configuration group name of current DAO.
	columns MerchantCountryConfigColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantCountryConfigColumns defines and stores column names for table merchant_country_config.
type MerchantCountryConfigColumns struct {
	Id          string //
	MerchantId  string //
	CountryCode string //
	Name        string //
	GmtCreate   string // create time
	GmtModify   string // update time
	IsDeleted   string // 0-UnDeleted，1-Deleted
	CreateTime  string // create utc time
}

// merchantCountryConfigColumns holds the columns for table merchant_country_config.
var merchantCountryConfigColumns = MerchantCountryConfigColumns{
	Id:          "id",
	MerchantId:  "merchant_id",
	CountryCode: "country_code",
	Name:        "name",
	GmtCreate:   "gmt_create",
	GmtModify:   "gmt_modify",
	IsDeleted:   "is_deleted",
	CreateTime:  "create_time",
}

// NewMerchantCountryConfigDao creates and returns a new DAO object for table data access.
func NewMerchantCountryConfigDao() *MerchantCountryConfigDao {
	return &MerchantCountryConfigDao{
		group:   "default",
		table:   "merchant_country_config",
		columns: merchantCountryConfigColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantCountryConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantCountryConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantCountryConfigDao) Columns() MerchantCountryConfigColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantCountryConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantCountryConfigDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantCountryConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
