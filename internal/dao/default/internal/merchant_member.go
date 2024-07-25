// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantMemberDao is the data access object for table merchant_member.
type MerchantMemberDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns MerchantMemberColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantMemberColumns defines and stores column names for table merchant_member.
type MerchantMemberColumns struct {
	Id         string // userId
	GmtCreate  string // create time
	GmtModify  string // update time
	MerchantId string // merchant id
	IsDeleted  string // 0-UnDeleted，1-Deleted
	Password   string // password
	UserName   string // user name
	Mobile     string // mobile
	Email      string // email
	FirstName  string // first name
	LastName   string // last name
	CreateTime string // create utc time
	Role       string // role
	Status     string // 0-Active, 2-Suspend
}

// merchantMemberColumns holds the columns for table merchant_member.
var merchantMemberColumns = MerchantMemberColumns{
	Id:         "id",
	GmtCreate:  "gmt_create",
	GmtModify:  "gmt_modify",
	MerchantId: "merchant_id",
	IsDeleted:  "is_deleted",
	Password:   "password",
	UserName:   "user_name",
	Mobile:     "mobile",
	Email:      "email",
	FirstName:  "first_name",
	LastName:   "last_name",
	CreateTime: "create_time",
	Role:       "role",
	Status:     "status",
}

// NewMerchantMemberDao creates and returns a new DAO object for table data access.
func NewMerchantMemberDao() *MerchantMemberDao {
	return &MerchantMemberDao{
		group:   "default",
		table:   "merchant_member",
		columns: merchantMemberColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantMemberDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantMemberDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantMemberDao) Columns() MerchantMemberColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantMemberDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantMemberDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantMemberDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
