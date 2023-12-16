// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OverseaPayChannelDao is the data access object for table oversea_pay_channel.
type OverseaPayChannelDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of current DAO.
	columns OverseaPayChannelColumns // columns contains all the column names of Table for convenient usage.
}

// OverseaPayChannelColumns defines and stores column names for table oversea_pay_channel.
type OverseaPayChannelColumns struct {
	Id               string // 主键id
	EnumKey          string // 支付渠道枚举（内部定义）
	Channel          string // 支付方式枚举（渠道定义）
	Name             string // 支付方式名称
	SubChannel       string // 渠道子支付方式枚举
	BrandData        string //
	Logo             string // 支付方式logo
	ChannelAccountId string // 渠道账户Id
	ChannelKey       string //
	ChannelSecret    string // secret
	Custom           string // custom
	GmtCreate        string // 创建时间
	GmtModify        string // 修改时间
	Description      string // 支付方式描述
}

// overseaPayChannelColumns holds the columns for table oversea_pay_channel.
var overseaPayChannelColumns = OverseaPayChannelColumns{
	Id:               "id",
	EnumKey:          "enum_key",
	Channel:          "channel",
	Name:             "name",
	SubChannel:       "sub_channel",
	BrandData:        "brand_data",
	Logo:             "logo",
	ChannelAccountId: "channel_account_id",
	ChannelKey:       "channel_key",
	ChannelSecret:    "channel_secret",
	Custom:           "custom",
	GmtCreate:        "gmt_create",
	GmtModify:        "gmt_modify",
	Description:      "description",
}

// NewOverseaPayChannelDao creates and returns a new DAO object for table data access.
func NewOverseaPayChannelDao() *OverseaPayChannelDao {
	return &OverseaPayChannelDao{
		group:   "oversea_pay",
		table:   "oversea_pay_channel",
		columns: overseaPayChannelColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *OverseaPayChannelDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *OverseaPayChannelDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *OverseaPayChannelDao) Columns() OverseaPayChannelColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *OverseaPayChannelDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *OverseaPayChannelDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *OverseaPayChannelDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
