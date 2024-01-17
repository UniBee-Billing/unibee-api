// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SubscriptionVatRateChannelDao is the data access object for table subscription_vat_rate_channel.
type SubscriptionVatRateChannelDao struct {
	table   string                            // table is the underlying table name of the DAO.
	group   string                            // group is the database configuration group name of current DAO.
	columns SubscriptionVatRateChannelColumns // columns contains all the column names of Table for convenient usage.
}

// SubscriptionVatRateChannelColumns defines and stores column names for table subscription_vat_rate_channel.
type SubscriptionVatRateChannelColumns struct {
	Id               string //
	GmtCreate        string // 创建时间
	GmtModify        string // 修改时间
	VatRateId        string // vat_rate_id
	ChannelId        string // 支付渠道Id
	ChannelVatRateId string // 支付渠道vat_rate_Id
	IsDeleted        string //
}

// subscriptionVatRateChannelColumns holds the columns for table subscription_vat_rate_channel.
var subscriptionVatRateChannelColumns = SubscriptionVatRateChannelColumns{
	Id:               "id",
	GmtCreate:        "gmt_create",
	GmtModify:        "gmt_modify",
	VatRateId:        "vat_rate_id",
	ChannelId:        "channel_id",
	ChannelVatRateId: "channel_vat_rate_id",
	IsDeleted:        "is_deleted",
}

// NewSubscriptionVatRateChannelDao creates and returns a new DAO object for table data access.
func NewSubscriptionVatRateChannelDao() *SubscriptionVatRateChannelDao {
	return &SubscriptionVatRateChannelDao{
		group:   "oversea_pay",
		table:   "subscription_vat_rate_channel",
		columns: subscriptionVatRateChannelColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SubscriptionVatRateChannelDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SubscriptionVatRateChannelDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SubscriptionVatRateChannelDao) Columns() SubscriptionVatRateChannelColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SubscriptionVatRateChannelDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SubscriptionVatRateChannelDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SubscriptionVatRateChannelDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
