// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantChannelConfigDao is the data access object for table merchant_channel_config.
type MerchantChannelConfigDao struct {
	table   string                       // table is the underlying table name of the DAO.
	group   string                       // group is the database configuration group name of current DAO.
	columns MerchantChannelConfigColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantChannelConfigColumns defines and stores column names for table merchant_channel_config.
type MerchantChannelConfigColumns struct {
	Id               string // 主键id
	MerchantId       string //
	EnumKey          string // 支付渠道枚举（内部定义）
	ChannelType      string // 支付渠道类型，null或者 0-Payment 类型 ｜ 1-Subscription类型
	Channel          string // 支付方式枚举（渠道定义）
	Name             string // 支付方式名称
	SubChannel       string // 渠道子支付方式枚举
	BrandData        string //
	Logo             string // 支付方式logo
	Host             string // pay host
	ChannelAccountId string // 渠道账户Id
	ChannelKey       string //
	ChannelSecret    string // secret
	Custom           string // custom
	GmtCreate        string // 创建时间
	GmtModify        string // 修改时间
	Description      string // 支付方式描述
	WebhookKey       string // webhook_key
	WebhookSecret    string // webhook_secret
	UniqueProductId  string // 渠道唯一 ProductId，目前仅限 Paypal 使用
}

// merchantChannelConfigColumns holds the columns for table merchant_channel_config.
var merchantChannelConfigColumns = MerchantChannelConfigColumns{
	Id:               "id",
	MerchantId:       "merchant_id",
	EnumKey:          "enum_key",
	ChannelType:      "channel_type",
	Channel:          "channel",
	Name:             "name",
	SubChannel:       "sub_channel",
	BrandData:        "brand_data",
	Logo:             "logo",
	Host:             "host",
	ChannelAccountId: "channel_account_id",
	ChannelKey:       "channel_key",
	ChannelSecret:    "channel_secret",
	Custom:           "custom",
	GmtCreate:        "gmt_create",
	GmtModify:        "gmt_modify",
	Description:      "description",
	WebhookKey:       "webhook_key",
	WebhookSecret:    "webhook_secret",
	UniqueProductId:  "unique_product_id",
}

// NewMerchantChannelConfigDao creates and returns a new DAO object for table data access.
func NewMerchantChannelConfigDao() *MerchantChannelConfigDao {
	return &MerchantChannelConfigDao{
		group:   "oversea_pay",
		table:   "merchant_channel_config",
		columns: merchantChannelConfigColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantChannelConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantChannelConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantChannelConfigDao) Columns() MerchantChannelConfigColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantChannelConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantChannelConfigDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantChannelConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}