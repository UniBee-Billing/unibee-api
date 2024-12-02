// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserAdminNoteDao is the data access object for table user_admin_note.
type UserAdminNoteDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns UserAdminNoteColumns // columns contains all the column names of Table for convenient usage.
}

// UserAdminNoteColumns defines and stores column names for table user_admin_note.
type UserAdminNoteColumns struct {
	Id               string // id
	GmtCreate        string // create_time
	GmtModify        string // modify_time
	UserId           string // user_id
	MerchantMemberId string // merchant_user_id
	Note             string // note
	IsDeleted        string // 0-UnDeleted，1-Deleted
	CreateTime       string // create utc time
}

// userAdminNoteColumns holds the columns for table user_admin_note.
var userAdminNoteColumns = UserAdminNoteColumns{
	Id:               "id",
	GmtCreate:        "gmt_create",
	GmtModify:        "gmt_modify",
	UserId:           "user_id",
	MerchantMemberId: "merchant_member_id",
	Note:             "note",
	IsDeleted:        "is_deleted",
	CreateTime:       "create_time",
}

// NewUserAdminNoteDao creates and returns a new DAO object for table data access.
func NewUserAdminNoteDao() *UserAdminNoteDao {
	return &UserAdminNoteDao{
		group:   "default",
		table:   "user_admin_note",
		columns: userAdminNoteColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UserAdminNoteDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UserAdminNoteDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UserAdminNoteDao) Columns() UserAdminNoteColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UserAdminNoteDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UserAdminNoteDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UserAdminNoteDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
