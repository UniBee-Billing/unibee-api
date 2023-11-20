// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OverseaRefundDao is the data access object for table oversea_refund.
type OverseaRefundDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns OverseaRefundColumns // columns contains all the column names of Table for convenient usage.
}

// OverseaRefundColumns defines and stores column names for table oversea_refund.
type OverseaRefundColumns struct {
	Id                   string // 主键ID
	CompanyId            string // 公司Id
	MerchantId           string // 商户ID
	BizType              string // 业务类型。同pay.biz_type
	BizId                string // 业务ID。同pay.biz_id
	CountryCode          string // 国家代码，指定发起交易的国家的两个字母的ISO 3166国家代码。目前支持SG、MY、PH、ID和TH
	Currency             string // 货币，“SGD” “MYR” “PHP” “IDR” “THB”\n与付款金额关联的货币。指定三个字母的ISO 4217货币代码
	OutTradeNo           string // 支付单号(内部生成，支付单号）
	OutRefundNo          string // 退款单号。可以唯一代表一笔退款（内部生成，退款单号）
	RefundFee            string // 退款金额。单位：分
	RefundComment        string // 退款备注
	RefundStatus         string // 退款状态。10-退款中，20-退款成功，30-退款失败
	RefundTime           string // 退款成功时间
	GmtCreate            string // 创建时间
	GmtModify            string // 更新时间
	ChannelRefundNo      string // 外部退款单号
	AppId                string // 退款使用的APPID
	RefundCommentExplain string // 退款备注说明
	NotifyUrl            string // 退款成功回调Url
	OpenApiId            string // 使用的开放平台配置Id
	ChannelId            string // 退款渠道Id
	ServiceRate          string // 服务费比例，万分位，百分比[0，10000)，精度为0.01%，如3即为0.03%
	AdditionalData       string //
}

// overseaRefundColumns holds the columns for table oversea_refund.
var overseaRefundColumns = OverseaRefundColumns{
	Id:                   "id",
	CompanyId:            "company_id",
	MerchantId:           "merchant_id",
	BizType:              "biz_type",
	BizId:                "biz_id",
	CountryCode:          "country_code",
	Currency:             "currency",
	OutTradeNo:           "out_trade_no",
	OutRefundNo:          "out_refund_no",
	RefundFee:            "refund_fee",
	RefundComment:        "refund_comment",
	RefundStatus:         "refund_status",
	RefundTime:           "refund_time",
	GmtCreate:            "gmt_create",
	GmtModify:            "gmt_modify",
	ChannelRefundNo:      "channel_refund_no",
	AppId:                "app_id",
	RefundCommentExplain: "refund_comment_explain",
	NotifyUrl:            "notify_url",
	OpenApiId:            "open_api_id",
	ChannelId:            "channel_id",
	ServiceRate:          "service_rate",
	AdditionalData:       "additional_data",
}

// NewOverseaRefundDao creates and returns a new DAO object for table data access.
func NewOverseaRefundDao() *OverseaRefundDao {
	return &OverseaRefundDao{
		group:   "oversea_pay",
		table:   "oversea_refund",
		columns: overseaRefundColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *OverseaRefundDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *OverseaRefundDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *OverseaRefundDao) Columns() OverseaRefundColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *OverseaRefundDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *OverseaRefundDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *OverseaRefundDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
