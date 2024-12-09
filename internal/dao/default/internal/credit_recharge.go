// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CreditRechargeDao is the data access object for table credit_recharge.
type CreditRechargeDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns CreditRechargeColumns // columns contains all the column names of Table for convenient usage.
}

// CreditRechargeColumns defines and stores column names for table credit_recharge.
type CreditRechargeColumns struct {
	Id                string // Id
	UserId            string // user_id
	CreditId          string // id of credit account
	RechargeId        string // unique recharge id for credit account
	RechargeStatus    string // recharge status, 10-in charging，20-recharge success，30-recharge failed
	Currency          string // currency
	TotalAmount       string // recharge total amount, cent
	PaymentAmount     string // the payment amount for recharge
	Name              string // recharge title
	Description       string // recharge description
	PaidTime          string // paid time
	GatewayId         string // payment gateway id
	InvoiceId         string // invoice_id
	PaymentId         string // paymentId
	TotalRefundAmount string // total refund amount,cent
	GmtCreate         string // create time
	GmtModify         string // update time
	CreateTime        string // create utc time
	MerchantId        string // merchant id
	AccountType       string // type of credit account, 1-main recharge account, 2-promo credit account
}

// creditRechargeColumns holds the columns for table credit_recharge.
var creditRechargeColumns = CreditRechargeColumns{
	Id:                "id",
	UserId:            "user_id",
	CreditId:          "credit_id",
	RechargeId:        "recharge_id",
	RechargeStatus:    "recharge_status",
	Currency:          "currency",
	TotalAmount:       "total_amount",
	PaymentAmount:     "payment_amount",
	Name:              "name",
	Description:       "description",
	PaidTime:          "paid_time",
	GatewayId:         "gateway_id",
	InvoiceId:         "invoice_id",
	PaymentId:         "payment_id",
	TotalRefundAmount: "total_refund_amount",
	GmtCreate:         "gmt_create",
	GmtModify:         "gmt_modify",
	CreateTime:        "create_time",
	MerchantId:        "merchant_id",
	AccountType:       "account_type",
}

// NewCreditRechargeDao creates and returns a new DAO object for table data access.
func NewCreditRechargeDao() *CreditRechargeDao {
	return &CreditRechargeDao{
		group:   "default",
		table:   "credit_recharge",
		columns: creditRechargeColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CreditRechargeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CreditRechargeDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CreditRechargeDao) Columns() CreditRechargeColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CreditRechargeDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CreditRechargeDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CreditRechargeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
