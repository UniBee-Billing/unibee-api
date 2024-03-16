// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RefundDao is the data access object for table refund.
type RefundDao struct {
	table   string        // table is the underlying table name of the DAO.
	group   string        // group is the database configuration group name of current DAO.
	columns RefundColumns // columns contains all the column names of Table for convenient usage.
}

// RefundColumns defines and stores column names for table refund.
type RefundColumns struct {
	Id                   string // id
	CompanyId            string // company id
	MerchantId           string // merchant id
	UserId               string // user_id
	OpenApiId            string // open api id
	GatewayId            string // gateway_id
	BizType              string // biz type, copy from payment.biz_type
	ExternalRefundId     string // external_refund_id
	CountryCode          string // country code
	Currency             string // currency
	PaymentId            string // relative payment id
	RefundId             string // refund id (system generate)
	RefundAmount         string // refund amount, cent
	RefundComment        string // refund comment
	Status               string // status。10-pending，20-success，30-failure, 40-cancel
	RefundTime           string // refund success time
	GmtCreate            string // create time
	GmtModify            string // update time
	GatewayRefundId      string // gateway refund id
	AppId                string // app id
	RefundCommentExplain string // refund comment
	ReturnUrl            string // return url after refund success
	MetaData             string // meta_data(json)
	UniqueId             string // unique id
	SubscriptionId       string // subscription id
	CreateTime           string // create utc time
}

// refundColumns holds the columns for table refund.
var refundColumns = RefundColumns{
	Id:                   "id",
	CompanyId:            "company_id",
	MerchantId:           "merchant_id",
	UserId:               "user_id",
	OpenApiId:            "open_api_id",
	GatewayId:            "gateway_id",
	BizType:              "biz_type",
	ExternalRefundId:     "external_refund_id",
	CountryCode:          "country_code",
	Currency:             "currency",
	PaymentId:            "payment_id",
	RefundId:             "refund_id",
	RefundAmount:         "refund_amount",
	RefundComment:        "refund_comment",
	Status:               "status",
	RefundTime:           "refund_time",
	GmtCreate:            "gmt_create",
	GmtModify:            "gmt_modify",
	GatewayRefundId:      "gateway_refund_id",
	AppId:                "app_id",
	RefundCommentExplain: "refund_comment_explain",
	ReturnUrl:            "return_url",
	MetaData:             "meta_data",
	UniqueId:             "unique_id",
	SubscriptionId:       "subscription_id",
	CreateTime:           "create_time",
}

// NewRefundDao creates and returns a new DAO object for table data access.
func NewRefundDao() *RefundDao {
	return &RefundDao{
		group:   "default",
		table:   "refund",
		columns: refundColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *RefundDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *RefundDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *RefundDao) Columns() RefundColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *RefundDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *RefundDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *RefundDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
