// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GatewayUserDao is the data access object for table gateway_user.
type GatewayUserDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns GatewayUserColumns // columns contains all the column names of Table for convenient usage.
}

// GatewayUserColumns defines and stores column names for table gateway_user.
type GatewayUserColumns struct {
	Id                          string //
	GmtCreate                   string // create time
	GmtModify                   string // update time
	UserId                      string // userId
	GatewayId                   string // gateway_id
	GatewayUserId               string // gateway_user_Id
	IsDeleted                   string // 0-UnDeleted，1-Deleted
	GatewayDefaultPaymentMethod string // gateway_default_payment_method
}

// gatewayUserColumns holds the columns for table gateway_user.
var gatewayUserColumns = GatewayUserColumns{
	Id:                          "id",
	GmtCreate:                   "gmt_create",
	GmtModify:                   "gmt_modify",
	UserId:                      "user_id",
	GatewayId:                   "gateway_id",
	GatewayUserId:               "gateway_user_id",
	IsDeleted:                   "is_deleted",
	GatewayDefaultPaymentMethod: "gateway_default_payment_method",
}

// NewGatewayUserDao creates and returns a new DAO object for table data access.
func NewGatewayUserDao() *GatewayUserDao {
	return &GatewayUserDao{
		group:   "oversea_pay",
		table:   "gateway_user",
		columns: gatewayUserColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *GatewayUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *GatewayUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *GatewayUserDao) Columns() GatewayUserColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *GatewayUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *GatewayUserDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *GatewayUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
