// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantCheckoutDao is the data access object for table merchant_checkout.
type MerchantCheckoutDao struct {
	table   string                  // table is the underlying table name of the DAO.
	group   string                  // group is the database configuration group name of current DAO.
	columns MerchantCheckoutColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantCheckoutColumns defines and stores column names for table merchant_checkout.
type MerchantCheckoutColumns struct {
	Id          string // ID
	MerchantId  string // merchantId
	Name        string // name
	Description string // description
	Data        string // data(json)
	Staging     string // staging_data(json)
	GmtCreate   string // create time
	GmtModify   string // update time
	IsDeleted   string // 0-UnDeleted，1-Deleted
	CreateTime  string // create utc time
}

// merchantCheckoutColumns holds the columns for table merchant_checkout.
var merchantCheckoutColumns = MerchantCheckoutColumns{
	Id:          "id",
	MerchantId:  "merchant_id",
	Name:        "name",
	Description: "description",
	Data:        "data",
	Staging:     "staging",
	GmtCreate:   "gmt_create",
	GmtModify:   "gmt_modify",
	IsDeleted:   "is_deleted",
	CreateTime:  "create_time",
}

// NewMerchantCheckoutDao creates and returns a new DAO object for table data access.
func NewMerchantCheckoutDao() *MerchantCheckoutDao {
	return &MerchantCheckoutDao{
		group:   "default",
		table:   "merchant_checkout",
		columns: merchantCheckoutColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantCheckoutDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantCheckoutDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantCheckoutDao) Columns() MerchantCheckoutColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantCheckoutDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantCheckoutDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantCheckoutDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
