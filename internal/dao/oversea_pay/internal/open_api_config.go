// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OpenApiConfigDao is the data access object for table open_api_config.
type OpenApiConfigDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns OpenApiConfigColumns // columns contains all the column names of Table for convenient usage.
}

// OpenApiConfigColumns defines and stores column names for table open_api_config.
type OpenApiConfigColumns struct {
	Id              string //
	Qps             string // 开放平台Api qps总控制
	GmtCreate       string // create time
	GmtModify       string // 修改时间
	MerchantId      string // 商户Id
	Hmac            string // 回调加密hmac
	Callback        string // 回调Url
	ApiKey          string // 开放平台Key
	Token           string // 开放平台token
	IsDeleted       string // 0-UnDeleted，1-Deleted
	Validips        string //
	ChannelCallback string // 渠道支付原信息回调地址
	CompanyId       string // 公司ID
}

// openApiConfigColumns holds the columns for table open_api_config.
var openApiConfigColumns = OpenApiConfigColumns{
	Id:              "id",
	Qps:             "qps",
	GmtCreate:       "gmt_create",
	GmtModify:       "gmt_modify",
	MerchantId:      "merchant_id",
	Hmac:            "hmac",
	Callback:        "callback",
	ApiKey:          "api_key",
	Token:           "token",
	IsDeleted:       "is_deleted",
	Validips:        "validips",
	ChannelCallback: "channel_callback",
	CompanyId:       "company_id",
}

// NewOpenApiConfigDao creates and returns a new DAO object for table data access.
func NewOpenApiConfigDao() *OpenApiConfigDao {
	return &OpenApiConfigDao{
		group:   "oversea_pay",
		table:   "open_api_config",
		columns: openApiConfigColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *OpenApiConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *OpenApiConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *OpenApiConfigDao) Columns() OpenApiConfigColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *OpenApiConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *OpenApiConfigDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *OpenApiConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
