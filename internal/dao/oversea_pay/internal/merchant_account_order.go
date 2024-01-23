// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantAccountOrderDao is the data access object for table merchant_account_order.
type MerchantAccountOrderDao struct {
	table   string                      // table is the underlying table name of the DAO.
	group   string                      // group is the database configuration group name of current DAO.
	columns MerchantAccountOrderColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantAccountOrderColumns defines and stores column names for table merchant_account_order.
type MerchantAccountOrderColumns struct {
	Id                            string //
	CompanyId                     string //
	MerchantId                    string // 商户ID
	MainId                        string // 结算单id，表merchant_account_main的id
	BizId                         string // 业务ID。可能是payId、refundId
	BizType                       string // 业务ID类型，1-pay,2-refund
	OrderType                     string // 账单类型；1支付；2 退款
	OrigCurrency                  string // 原始货币类型
	OrigTradeFee                  string // 原始货币交易金额。单位：分
	CurrencyRate                  string // 汇率，万分位
	CurrencyRateDataJson          string // 汇率数据JSON结构
	Currency                      string
	TradeFee                      string // 交易金额（正值，退款代表退还金额）。单位：分
	DeductPoint                   string // 服务费扣点，万分位
	DeductFee                     string // 扣点金额（正值，退款代表退服务费）。单位：分
	BillFee                       string // 结算金额（正值，退款代表退还金额）。单位：分
	GmtCreate                     string //
	GmtModify                     string //
	MerchantReference             string // 客户订单号
	ChannelOrderNo                string //
	MerchantOrderNo               string //
	Channel                       string //
	PaymentCurrency               string //
	Authorised                    string //
	Captured                      string //
	CurrencyRateMarkup            string //
	PaymentMethodVariant          string //
	ModificationMerchantReference string //
	MerchantOrderReference        string //
	DelayCaptureTime              string //
	NotificationUrl               string //
	OpenAppId                     string //
}

// merchantAccountOrderColumns holds the columns for table merchant_account_order.
var merchantAccountOrderColumns = MerchantAccountOrderColumns{
	Id:                            "id",
	CompanyId:                     "company_id",
	MerchantId:                    "merchant_id",
	MainId:                        "main_id",
	BizId:                         "biz_id",
	BizType:                       "biz_type",
	OrderType:                     "order_type",
	OrigCurrency:                  "orig_currency",
	OrigTradeFee:                  "orig_trade_fee",
	CurrencyRate:                  "currency_rate",
	CurrencyRateDataJson:          "currency_rate_data_json",
	Currency:                      "currency",
	TradeFee:                      "trade_fee",
	DeductPoint:                   "deduct_point",
	DeductFee:                     "deduct_fee",
	BillFee:                       "bill_fee",
	GmtCreate:                     "gmt_create",
	GmtModify:                     "gmt_modify",
	MerchantReference:             "merchant_reference",
	ChannelOrderNo:                "channel_order_no",
	MerchantOrderNo:               "merchant_order_no",
	Channel:                       "channel",
	PaymentCurrency:               "payment_currency",
	Authorised:                    "authorised",
	Captured:                      "captured",
	CurrencyRateMarkup:            "currency_rate_markup",
	PaymentMethodVariant:          "payment_method_variant",
	ModificationMerchantReference: "modification_merchant_reference",
	MerchantOrderReference:        "merchant_order_reference",
	DelayCaptureTime:              "delay_capture_time",
	NotificationUrl:               "notification_url",
	OpenAppId:                     "open_app_id",
}

// NewMerchantAccountOrderDao creates and returns a new DAO object for table data access.
func NewMerchantAccountOrderDao() *MerchantAccountOrderDao {
	return &MerchantAccountOrderDao{
		group:   "oversea_pay",
		table:   "merchant_account_order",
		columns: merchantAccountOrderColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantAccountOrderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantAccountOrderDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantAccountOrderDao) Columns() MerchantAccountOrderColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantAccountOrderDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantAccountOrderDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantAccountOrderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
