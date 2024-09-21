// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentItemDao is the data access object for table payment_item.
type PaymentItemDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns PaymentItemColumns // columns contains all the column names of Table for convenient usage.
}

// PaymentItemColumns defines and stores column names for table payment_item.
type PaymentItemColumns struct {
	Id             string //
	BizType        string // biz_type 1-onetime payment, 3-subscription
	Status         string // 0-pending, 1-success, 2-failure
	MerchantId     string // merchant id
	UserId         string // userId
	SubscriptionId string // subscription id
	InvoiceId      string // invoice id
	UniqueId       string // unique id
	Currency       string // currency
	Amount         string // amount
	UnitAmount     string // unit_amount
	Quantity       string // quantity
	GatewayId      string // gateway id
	GmtCreate      string // create time
	GmtModify      string // update time
	IsDeleted      string // 0-UnDeleted，1-Deleted
	PaymentId      string // PaymentId
	CreateTime     string // create utc time
	Description    string // description
	Name           string // name
}

// paymentItemColumns holds the columns for table payment_item.
var paymentItemColumns = PaymentItemColumns{
	Id:             "id",
	BizType:        "biz_type",
	Status:         "status",
	MerchantId:     "merchant_id",
	UserId:         "user_id",
	SubscriptionId: "subscription_id",
	InvoiceId:      "invoice_id",
	UniqueId:       "unique_id",
	Currency:       "currency",
	Amount:         "amount",
	UnitAmount:     "unit_amount",
	Quantity:       "quantity",
	GatewayId:      "gateway_id",
	GmtCreate:      "gmt_create",
	GmtModify:      "gmt_modify",
	IsDeleted:      "is_deleted",
	PaymentId:      "payment_id",
	CreateTime:     "create_time",
	Description:    "description",
	Name:           "name",
}

// NewPaymentItemDao creates and returns a new DAO object for table data access.
func NewPaymentItemDao() *PaymentItemDao {
	return &PaymentItemDao{
		group:   "default",
		table:   "payment_item",
		columns: paymentItemColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PaymentItemDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PaymentItemDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PaymentItemDao) Columns() PaymentItemColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PaymentItemDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PaymentItemDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PaymentItemDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
