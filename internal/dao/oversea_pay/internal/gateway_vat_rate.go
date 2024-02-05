// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GatewayVatRateDao is the data access object for table gateway_vat_rate.
type GatewayVatRateDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns GatewayVatRateColumns // columns contains all the column names of Table for convenient usage.
}

// GatewayVatRateColumns defines and stores column names for table gateway_vat_rate.
type GatewayVatRateColumns struct {
	Id               string //
	GmtCreate        string // create time
	GmtModify        string // update time
	VatRateId        string // vat_rate_id
	GatewayId        string // gateway_id
	GatewayVatRateId string // gateway_vat_rate_Id
	IsDeleted        string // 0-UnDeleted，1-Deleted
}

// gatewayVatRateColumns holds the columns for table gateway_vat_rate.
var gatewayVatRateColumns = GatewayVatRateColumns{
	Id:               "id",
	GmtCreate:        "gmt_create",
	GmtModify:        "gmt_modify",
	VatRateId:        "vat_rate_id",
	GatewayId:        "gateway_id",
	GatewayVatRateId: "gateway_vat_rate_id",
	IsDeleted:        "is_deleted",
}

// NewGatewayVatRateDao creates and returns a new DAO object for table data access.
func NewGatewayVatRateDao() *GatewayVatRateDao {
	return &GatewayVatRateDao{
		group:   "oversea_pay",
		table:   "gateway_vat_rate",
		columns: gatewayVatRateColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *GatewayVatRateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *GatewayVatRateDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *GatewayVatRateDao) Columns() GatewayVatRateColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *GatewayVatRateDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *GatewayVatRateDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *GatewayVatRateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
