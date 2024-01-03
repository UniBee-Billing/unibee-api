// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantUserAccountDao is the data access object for table merchant_user_account.
type MerchantUserAccountDao struct {
	table   string                     // table is the underlying table name of the DAO.
	group   string                     // group is the database configuration group name of current DAO.
	columns MerchantUserAccountColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantUserAccountColumns defines and stores column names for table merchant_user_account.
type MerchantUserAccountColumns struct {
	Id         string // userId
	GmtCreate  string // 创建时间
	GmtModify  string // 修改时间
	MerchantId string // 用户ID
	IsDeleted  string // 逻辑删除
	Password   string // 密码，加密存储
	UserName   string // 用户名
	Mobile     string // 手机号
	Email      string // 邮箱
	FirstName  string //
	LastName   string //
}

// merchantUserAccountColumns holds the columns for table merchant_user_account.
var merchantUserAccountColumns = MerchantUserAccountColumns{
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
}

// NewMerchantUserAccountDao creates and returns a new DAO object for table data access.
func NewMerchantUserAccountDao() *MerchantUserAccountDao {
	return &MerchantUserAccountDao{
		group:   "oversea_pay",
		table:   "merchant_user_account",
		columns: merchantUserAccountColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantUserAccountDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantUserAccountDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantUserAccountDao) Columns() MerchantUserAccountColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantUserAccountDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantUserAccountDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantUserAccountDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
