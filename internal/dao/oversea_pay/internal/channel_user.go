// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ChannelUserDao is the data access object for table channel_user.
type ChannelUserDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns ChannelUserColumns // columns contains all the column names of Table for convenient usage.
}

// ChannelUserColumns defines and stores column names for table channel_user.
type ChannelUserColumns struct {
	Id                          string //
	GmtCreate                   string // 创建时间
	GmtModify                   string // 修改时间
	UserId                      string // userId
	ChannelId                   string // 支付渠道Id
	ChannelUserId               string // 支付渠道user_Id
	IsDeleted                   string // 0-UnDeleted，1-Deleted
	ChannelDefaultPaymentMethod string //
}

// channelUserColumns holds the columns for table channel_user.
var channelUserColumns = ChannelUserColumns{
	Id:                          "id",
	GmtCreate:                   "gmt_create",
	GmtModify:                   "gmt_modify",
	UserId:                      "user_id",
	ChannelId:                   "channel_id",
	ChannelUserId:               "channel_user_id",
	IsDeleted:                   "is_deleted",
	ChannelDefaultPaymentMethod: "channel_default_payment_method",
}

// NewChannelUserDao creates and returns a new DAO object for table data access.
func NewChannelUserDao() *ChannelUserDao {
	return &ChannelUserDao{
		group:   "oversea_pay",
		table:   "channel_user",
		columns: channelUserColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ChannelUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ChannelUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ChannelUserDao) Columns() ChannelUserColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ChannelUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ChannelUserDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ChannelUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
