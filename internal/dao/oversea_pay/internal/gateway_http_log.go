// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GatewayHttpLogDao is the data access object for table gateway_http_log.
type GatewayHttpLogDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns GatewayHttpLogColumns // columns contains all the column names of Table for convenient usage.
}

// GatewayHttpLogColumns defines and stores column names for table gateway_http_log.
type GatewayHttpLogColumns struct {
	Id        string // id
	Url       string // request url
	Request   string // request body(json)
	Response  string // response(json)
	RequestId string // request_id
	Mamo      string // mamo
	GatewayId string // gateway_id
	GmtCreate string // create time
	GmtModify string // update time
}

// gatewayHttpLogColumns holds the columns for table gateway_http_log.
var gatewayHttpLogColumns = GatewayHttpLogColumns{
	Id:        "id",
	Url:       "url",
	Request:   "request",
	Response:  "response",
	RequestId: "request_id",
	Mamo:      "mamo",
	GatewayId: "gateway_id",
	GmtCreate: "gmt_create",
	GmtModify: "gmt_modify",
}

// NewGatewayHttpLogDao creates and returns a new DAO object for table data access.
func NewGatewayHttpLogDao() *GatewayHttpLogDao {
	return &GatewayHttpLogDao{
		group:   "oversea_pay",
		table:   "gateway_http_log",
		columns: gatewayHttpLogColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *GatewayHttpLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *GatewayHttpLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *GatewayHttpLogDao) Columns() GatewayHttpLogColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *GatewayHttpLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *GatewayHttpLogDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *GatewayHttpLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
