// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantUserDiscountCodeDao is the data access object for table merchant_user_discount_code.
type MerchantUserDiscountCodeDao struct {
	table   string                          // table is the underlying table name of the DAO.
	group   string                          // group is the database configuration group name of current DAO.
	columns MerchantUserDiscountCodeColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantUserDiscountCodeColumns defines and stores column names for table merchant_user_discount_code.
type MerchantUserDiscountCodeColumns struct {
	Id             string // ID
	MerchantId     string // merchantId
	UserId         string // user_id
	Code           string // code
	Status         string // status, 1-normal, 2-rollback
	PlanId         string // plan_id
	SubscriptionId string // subscription_id
	PaymentId      string // payment_id
	InvoiceId      string // invoice_id
	UniqueId       string // unique_id
	GmtCreate      string // create time
	GmtModify      string // update time
	IsDeleted      string // 0-UnDeleted，1-Deleted
	CreateTime     string // create utc time
	ApplyAmount    string // apply_amount
	Currency       string // currency
}

// merchantUserDiscountCodeColumns holds the columns for table merchant_user_discount_code.
var merchantUserDiscountCodeColumns = MerchantUserDiscountCodeColumns{
	Id:             "id",
	MerchantId:     "merchant_id",
	UserId:         "user_id",
	Code:           "code",
	Status:         "status",
	PlanId:         "plan_id",
	SubscriptionId: "subscription_id",
	PaymentId:      "payment_id",
	InvoiceId:      "invoice_id",
	UniqueId:       "unique_id",
	GmtCreate:      "gmt_create",
	GmtModify:      "gmt_modify",
	IsDeleted:      "is_deleted",
	CreateTime:     "create_time",
	ApplyAmount:    "apply_amount",
	Currency:       "currency",
}

// NewMerchantUserDiscountCodeDao creates and returns a new DAO object for table data access.
func NewMerchantUserDiscountCodeDao() *MerchantUserDiscountCodeDao {
	return &MerchantUserDiscountCodeDao{
		group:   "default",
		table:   "merchant_user_discount_code",
		columns: merchantUserDiscountCodeColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantUserDiscountCodeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantUserDiscountCodeDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantUserDiscountCodeDao) Columns() MerchantUserDiscountCodeColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantUserDiscountCodeDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantUserDiscountCodeDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantUserDiscountCodeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
