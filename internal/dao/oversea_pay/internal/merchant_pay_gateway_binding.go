// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantPayGatewayBindingDao is the data access object for table merchant_pay_gateway_binding.
type MerchantPayGatewayBindingDao struct {
	table   string                           // table is the underlying table name of the DAO.
	group   string                           // group is the database configuration group name of current DAO.
	columns MerchantPayGatewayBindingColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantPayGatewayBindingColumns defines and stores column names for table merchant_pay_gateway_binding.
type MerchantPayGatewayBindingColumns struct {
	Id         string //
	GmtCreate  string // create time
	GmtModify  string // update time
	MerchantId string // merchant id
	GatewayId  string // gateway_id
	IsDeleted  string // 0-UnDeleted，1-Deleted
	CreateTime string // create utc time
}

// merchantPayGatewayBindingColumns holds the columns for table merchant_pay_gateway_binding.
var merchantPayGatewayBindingColumns = MerchantPayGatewayBindingColumns{
	Id:         "id",
	GmtCreate:  "gmt_create",
	GmtModify:  "gmt_modify",
	MerchantId: "merchant_id",
	GatewayId:  "gateway_id",
	IsDeleted:  "is_deleted",
	CreateTime: "create_time",
}

// NewMerchantPayGatewayBindingDao creates and returns a new DAO object for table data access.
func NewMerchantPayGatewayBindingDao() *MerchantPayGatewayBindingDao {
	return &MerchantPayGatewayBindingDao{
		group:   "oversea_pay",
		table:   "merchant_pay_gateway_binding",
		columns: merchantPayGatewayBindingColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantPayGatewayBindingDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantPayGatewayBindingDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantPayGatewayBindingDao) Columns() MerchantPayGatewayBindingColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantPayGatewayBindingDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantPayGatewayBindingDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantPayGatewayBindingDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
