// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ChannelHttpLogDao is the data access object for table channel_http_log.
type ChannelHttpLogDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns ChannelHttpLogColumns // columns contains all the column names of Table for convenient usage.
}

// ChannelHttpLogColumns defines and stores column names for table channel_http_log.
type ChannelHttpLogColumns struct {
	Id        string // id
	Url       string // 请求url
	Request   string // 请求body参数(json格式)
	Response  string // 请求返回结果(json格式)
	RequestId string // reuqest_id
	Mamo      string // 备注
	ChannelId string // channel_id
	GmtCreate string // 创建时间
	GmtModify string // 更新时间
}

// channelHttpLogColumns holds the columns for table channel_http_log.
var channelHttpLogColumns = ChannelHttpLogColumns{
	Id:        "id",
	Url:       "url",
	Request:   "request",
	Response:  "response",
	RequestId: "request_id",
	Mamo:      "mamo",
	ChannelId: "channel_id",
	GmtCreate: "gmt_create",
	GmtModify: "gmt_modify",
}

// NewChannelHttpLogDao creates and returns a new DAO object for table data access.
func NewChannelHttpLogDao() *ChannelHttpLogDao {
	return &ChannelHttpLogDao{
		group:   "oversea_pay",
		table:   "channel_http_log",
		columns: channelHttpLogColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ChannelHttpLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ChannelHttpLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ChannelHttpLogDao) Columns() ChannelHttpLogColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ChannelHttpLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ChannelHttpLogDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ChannelHttpLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
