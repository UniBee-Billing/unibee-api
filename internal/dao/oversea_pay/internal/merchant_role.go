// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantRoleDao is the data access object for table merchant_role.
type MerchantRoleDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns MerchantRoleColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantRoleColumns defines and stores column names for table merchant_role.
type MerchantRoleColumns struct {
	Id             string // userId
	GmtCreate      string // create time
	GmtModify      string // update time
	MerchantId     string // merchant id
	IsDeleted      string // 0-UnDeleted，1-Deleted
	Role           string // role
	PermissionData string // permission_data（json）
	CreateTime     string // create utc time
}

// merchantRoleColumns holds the columns for table merchant_role.
var merchantRoleColumns = MerchantRoleColumns{
	Id:             "id",
	GmtCreate:      "gmt_create",
	GmtModify:      "gmt_modify",
	MerchantId:     "merchant_id",
	IsDeleted:      "is_deleted",
	Role:           "role",
	PermissionData: "permission_data",
	CreateTime:     "create_time",
}

// NewMerchantRoleDao creates and returns a new DAO object for table data access.
func NewMerchantRoleDao() *MerchantRoleDao {
	return &MerchantRoleDao{
		group:   "default",
		table:   "merchant_role",
		columns: merchantRoleColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantRoleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantRoleDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantRoleDao) Columns() MerchantRoleColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantRoleDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantRoleDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantRoleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
