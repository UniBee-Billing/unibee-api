// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantPayChannelMappingDao is the data access object for table merchant_pay_channel_mapping.
type MerchantPayChannelMappingDao struct {
	table   string                           // table is the underlying table name of the DAO.
	group   string                           // group is the database configuration group name of current DAO.
	columns MerchantPayChannelMappingColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantPayChannelMappingColumns defines and stores column names for table merchant_pay_channel_mapping.
type MerchantPayChannelMappingColumns struct {
	Id         string //
	GmtCreate  string // create time
	GmtModify  string // 修改时间
	MerchantId string // 商户Id
	ChannelId  string // oversea_pay_channel表的id
	IsDeleted  string // 0-UnDeleted，1-Deleted
}

// merchantPayChannelMappingColumns holds the columns for table merchant_pay_channel_mapping.
var merchantPayChannelMappingColumns = MerchantPayChannelMappingColumns{
	Id:         "id",
	GmtCreate:  "gmt_create",
	GmtModify:  "gmt_modify",
	MerchantId: "merchant_id",
	ChannelId:  "channel_id",
	IsDeleted:  "is_deleted",
}

// NewMerchantPayChannelMappingDao creates and returns a new DAO object for table data access.
func NewMerchantPayChannelMappingDao() *MerchantPayChannelMappingDao {
	return &MerchantPayChannelMappingDao{
		group:   "oversea_pay",
		table:   "merchant_pay_channel_mapping",
		columns: merchantPayChannelMappingColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantPayChannelMappingDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantPayChannelMappingDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantPayChannelMappingDao) Columns() MerchantPayChannelMappingColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantPayChannelMappingDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantPayChannelMappingDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantPayChannelMappingDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
