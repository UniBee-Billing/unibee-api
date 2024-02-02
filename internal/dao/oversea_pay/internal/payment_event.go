// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentEventDao is the data access object for table payment_event.
type PaymentEventDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns PaymentEventColumns // columns contains all the column names of Table for convenient usage.
}

// PaymentEventColumns defines and stores column names for table payment_event.
type PaymentEventColumns struct {
	Id              string // id
	UserId          string // user_id
	MerchantUserId  string // merchant_user_id
	OpenApiId       string // open api id
	TerminalIp      string // terminal_ip
	BizType         string // biz_type=1，Payment表
	BizId           string // biz_type=1，pay；
	Fee             string // amount, cent
	EventType       string // 0-unknown
	Event           string // event
	RelativeTradeNo string // relative trade no
	UniqueNo        string // unique no
	GmtCreate       string // create time
	GmtModify       string // update time
	Message         string // message
}

// paymentEventColumns holds the columns for table payment_event.
var paymentEventColumns = PaymentEventColumns{
	Id:              "id",
	UserId:          "user_id",
	MerchantUserId:  "merchant_user_id",
	OpenApiId:       "open_api_id",
	TerminalIp:      "terminal_ip",
	BizType:         "biz_type",
	BizId:           "biz_id",
	Fee:             "fee",
	EventType:       "event_type",
	Event:           "event",
	RelativeTradeNo: "relative_trade_no",
	UniqueNo:        "unique_no",
	GmtCreate:       "gmt_create",
	GmtModify:       "gmt_modify",
	Message:         "message",
}

// NewPaymentEventDao creates and returns a new DAO object for table data access.
func NewPaymentEventDao() *PaymentEventDao {
	return &PaymentEventDao{
		group:   "oversea_pay",
		table:   "payment_event",
		columns: paymentEventColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PaymentEventDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PaymentEventDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PaymentEventDao) Columns() PaymentEventColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PaymentEventDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PaymentEventDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PaymentEventDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
