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
	Id                   string // 主键ID
	CompanyId            string // 公司Id
	MerchantId           string // 商户ID
	UserId               string // user_id
	OpenApiId            string // 使用的开放平台配置Id
	ChannelId            string // 退款渠道Id
	BizType              string // 业务类型。同pay.biz_type
	BizId                string // 业务ID。同pay.biz_id
	CountryCode          string // 国家代码，指定发起交易的国家的两个字母的ISO 3166国家代码。目前支持SG、MY、PH、ID和TH
	Currency             string // 货币，“SGD” “MYR” “PHP” “IDR” “THB”\n与付款金额关联的货币。指定三个字母的ISO 4217货币代码
	PaymentId            string // 支付单号(内部生成，支付单号）
	RefundId             string // 退款单号。可以唯一代表一笔退款（内部生成，退款单号）
	RefundFee            string // 退款金额。单位：分
	RefundComment        string // 退款备注
	Status               string // 退款状态。10-退款中，20-退款成功，30-退款失败
	RefundTime           string // 退款成功时间
	GmtCreate            string // 创建时间
	GmtModify            string // 更新时间
	ChannelRefundId      string // 外部退款单号
	AppId                string // 退款使用的APPID
	RefundCommentExplain string // 退款备注说明
	ReturnUrl            string // 退款成功回调Url
	AdditionalData       string //
	UniqueId             string // 唯一键，以同步为逻辑加入使用自定义唯一键
	SubscriptionId       string // 订阅id（内部编号）
	ChannelPaymentId     string // 外部支付渠道订单号，支付成功回调返回
}

// refundColumns holds the columns for table refund.
var refundColumns = RefundColumns{
	Id:                   "id",
	CompanyId:            "company_id",
	MerchantId:           "merchant_id",
	UserId:               "user_id",
	OpenApiId:            "open_api_id",
	ChannelId:            "channel_id",
	BizType:              "biz_type",
	BizId:                "biz_id",
	CountryCode:          "country_code",
	Currency:             "currency",
	PaymentId:            "payment_id",
	RefundId:             "refund_id",
	RefundFee:            "refund_fee",
	RefundComment:        "refund_comment",
	Status:               "status",
	RefundTime:           "refund_time",
	GmtCreate:            "gmt_create",
	GmtModify:            "gmt_modify",
	ChannelRefundId:      "channel_refund_id",
	AppId:                "app_id",
	RefundCommentExplain: "refund_comment_explain",
	ReturnUrl:            "return_url",
	AdditionalData:       "additional_data",
	UniqueId:             "unique_id",
	SubscriptionId:       "subscription_id",
	ChannelPaymentId:     "channel_payment_id",
}

// NewRefundDao creates and returns a new DAO object for table data access.
func NewRefundDao() *RefundDao {
	return &RefundDao{
		group:   "oversea_pay",
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
