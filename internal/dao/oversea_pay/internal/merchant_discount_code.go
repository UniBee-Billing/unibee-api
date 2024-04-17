// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantDiscountCodeDao is the data access object for table merchant_discount_code.
type MerchantDiscountCodeDao struct {
	table   string                      // table is the underlying table name of the DAO.
	group   string                      // group is the database configuration group name of current DAO.
	columns MerchantDiscountCodeColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantDiscountCodeColumns defines and stores column names for table merchant_discount_code.
type MerchantDiscountCodeColumns struct {
	Id                 string // ID
	MerchantId         string // merchantId
	Name               string // name
	Code               string // code
	Status             string // status, 1-editable, 2-active, 3-deactive, 4-expire
	BillingType        string // billing_type, 1-one-time, 2-recurring
	DiscountType       string // discount_type, 1-percentage, 2-fixed_amount
	DiscountAmount     string // amount of discount, available when discount_type is fixed_amount
	DiscountPercentage string // percentage of discount, 100=1%, available when discount_type is percentage
	Currency           string // currency of discount, available when discount_type is fixed_amount
	UserLimit          string // the limit of every user apply, 0-unlimited
	SubscriptionLimit  string // the limit of every subscription apply, 0-unlimited
	StartTime          string // start of discount available utc time
	EndTime            string // end of discount available utc time, 0-invalid
	GmtCreate          string // create time
	GmtModify          string // update time
	IsDeleted          string // 0-UnDeleted，1-Deleted
	CreateTime         string // create utc time
	CycleLimit         string // the count limitation of subscription cycle , 0-no limit
	MetaData           string // meta_data(json)
	Type               string // type, 1-external discount code
}

// merchantDiscountCodeColumns holds the columns for table merchant_discount_code.
var merchantDiscountCodeColumns = MerchantDiscountCodeColumns{
	Id:                 "id",
	MerchantId:         "merchant_id",
	Name:               "name",
	Code:               "code",
	Status:             "status",
	BillingType:        "billing_type",
	DiscountType:       "discount_type",
	DiscountAmount:     "discount_amount",
	DiscountPercentage: "discount_percentage",
	Currency:           "currency",
	UserLimit:          "user_limit",
	SubscriptionLimit:  "subscription_limit",
	StartTime:          "start_time",
	EndTime:            "end_time",
	GmtCreate:          "gmt_create",
	GmtModify:          "gmt_modify",
	IsDeleted:          "is_deleted",
	CreateTime:         "create_time",
	CycleLimit:         "cycle_limit",
	MetaData:           "meta_data",
	Type:               "type",
}

// NewMerchantDiscountCodeDao creates and returns a new DAO object for table data access.
func NewMerchantDiscountCodeDao() *MerchantDiscountCodeDao {
	return &MerchantDiscountCodeDao{
		group:   "default",
		table:   "merchant_discount_code",
		columns: merchantDiscountCodeColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantDiscountCodeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantDiscountCodeDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantDiscountCodeDao) Columns() MerchantDiscountCodeColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantDiscountCodeDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantDiscountCodeDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantDiscountCodeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
