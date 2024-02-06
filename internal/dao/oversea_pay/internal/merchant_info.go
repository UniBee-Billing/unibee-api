// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantInfoDao is the data access object for table merchant_info.
type MerchantInfoDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns MerchantInfoColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantInfoColumns defines and stores column names for table merchant_info.
type MerchantInfoColumns struct {
	Id          string // merchant_id
	CompanyId   string // company_id
	UserId      string // create_user_id
	Type        string // type
	CompanyName string // company_name
	Email       string // email
	BusinessNum string // business_num
	Name        string // name
	Idcard      string // idcard
	Location    string // location
	Address     string // address
	GmtCreate   string // create time
	GmtModify   string // update_time
	IsDeleted   string // 0-UnDeleted，1-Deleted
	CompanyLogo string // company_logo
	HomeUrl     string //
	FirstName   string // first_name
	LastName    string // last_name
	Phone       string // phone
	CreateAt    string // create utc time
	TimeZone    string // merchant default time zone
}

// merchantInfoColumns holds the columns for table merchant_info.
var merchantInfoColumns = MerchantInfoColumns{
	Id:          "id",
	CompanyId:   "company_id",
	UserId:      "user_id",
	Type:        "type",
	CompanyName: "company_name",
	Email:       "email",
	BusinessNum: "business_num",
	Name:        "name",
	Idcard:      "idcard",
	Location:    "location",
	Address:     "address",
	GmtCreate:   "gmt_create",
	GmtModify:   "gmt_modify",
	IsDeleted:   "is_deleted",
	CompanyLogo: "company_logo",
	HomeUrl:     "home_url",
	FirstName:   "first_name",
	LastName:    "last_name",
	Phone:       "phone",
	CreateAt:    "create_at",
	TimeZone:    "time_zone",
}

// NewMerchantInfoDao creates and returns a new DAO object for table data access.
func NewMerchantInfoDao() *MerchantInfoDao {
	return &MerchantInfoDao{
		group:   "oversea_pay",
		table:   "merchant_info",
		columns: merchantInfoColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantInfoDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantInfoDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantInfoDao) Columns() MerchantInfoColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantInfoDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantInfoDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantInfoDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
