// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SubscriptionAdminNoteDao is the data access object for table subscription_admin_note.
type SubscriptionAdminNoteDao struct {
	table   string                       // table is the underlying table name of the DAO.
	group   string                       // group is the database configuration group name of current DAO.
	columns SubscriptionAdminNoteColumns // columns contains all the column names of Table for convenient usage.
}

// SubscriptionAdminNoteColumns defines and stores column names for table subscription_admin_note.
type SubscriptionAdminNoteColumns struct {
	Id               string // id
	GmtCreate        string // create_time
	GmtModify        string // modify_time
	SubscriptionId   string // subscription_id
	MerchantMemberId string // merchant_user_id
	Note             string // note
	IsDeleted        string // 0-UnDeleted，1-Deleted
	CreateTime       string // create utc time
}

// subscriptionAdminNoteColumns holds the columns for table subscription_admin_note.
var subscriptionAdminNoteColumns = SubscriptionAdminNoteColumns{
	Id:               "id",
	GmtCreate:        "gmt_create",
	GmtModify:        "gmt_modify",
	SubscriptionId:   "subscription_id",
	MerchantMemberId: "merchant_member_id",
	Note:             "note",
	IsDeleted:        "is_deleted",
	CreateTime:       "create_time",
}

// NewSubscriptionAdminNoteDao creates and returns a new DAO object for table data access.
func NewSubscriptionAdminNoteDao() *SubscriptionAdminNoteDao {
	return &SubscriptionAdminNoteDao{
		group:   "oversea_pay",
		table:   "subscription_admin_note",
		columns: subscriptionAdminNoteColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SubscriptionAdminNoteDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SubscriptionAdminNoteDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SubscriptionAdminNoteDao) Columns() SubscriptionAdminNoteColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SubscriptionAdminNoteDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SubscriptionAdminNoteDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SubscriptionAdminNoteDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
