// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ChannelVatRateDao is the data access object for table channel_vat_rate.
type ChannelVatRateDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns ChannelVatRateColumns // columns contains all the column names of Table for convenient usage.
}

// ChannelVatRateColumns defines and stores column names for table channel_vat_rate.
type ChannelVatRateColumns struct {
	Id               string //
	GmtCreate        string // create time
	GmtModify        string // update time
	VatRateId        string // vat_rate_id
	ChannelId        string // channel_id
	ChannelVatRateId string // channel_vat_rate_Id
	IsDeleted        string // 0-UnDeleted，1-Deleted
}

// channelVatRateColumns holds the columns for table channel_vat_rate.
var channelVatRateColumns = ChannelVatRateColumns{
	Id:               "id",
	GmtCreate:        "gmt_create",
	GmtModify:        "gmt_modify",
	VatRateId:        "vat_rate_id",
	ChannelId:        "channel_id",
	ChannelVatRateId: "channel_vat_rate_id",
	IsDeleted:        "is_deleted",
}

// NewChannelVatRateDao creates and returns a new DAO object for table data access.
func NewChannelVatRateDao() *ChannelVatRateDao {
	return &ChannelVatRateDao{
		group:   "oversea_pay",
		table:   "channel_vat_rate",
		columns: channelVatRateColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ChannelVatRateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ChannelVatRateDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ChannelVatRateDao) Columns() ChannelVatRateColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ChannelVatRateDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ChannelVatRateDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ChannelVatRateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
