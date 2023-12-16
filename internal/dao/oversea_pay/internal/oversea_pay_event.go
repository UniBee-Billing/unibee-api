// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OverseaPayEventDao is the data access object for table oversea_pay_event.
type OverseaPayEventDao struct {
	table   string                 // table is the underlying table name of the DAO.
	group   string                 // group is the database configuration group name of current DAO.
	columns OverseaPayEventColumns // columns contains all the column names of Table for convenient usage.
}

// OverseaPayEventColumns defines and stores column names for table oversea_pay_event.
type OverseaPayEventColumns struct {
	Id              string // 主键id
	BizType         string // biz_type=0，oversea_pay表
	BizId           string // biz_type=0，oversea_pay表Id；
	Fee             string // 金额（分）
	EventType       string // 0-未知
	Event           string // 事件
	RelativeTradeNo string // 关联单号
	UniqueNo        string // 唯一键
	GmtCreate       string // 创建时间
	GmtModify       string // 更新时间
	OpenApiId       string // 使用的开放平台配置Id
	TerminalIp      string // 实时交易终端IP
	Message         string // message
}

// overseaPayEventColumns holds the columns for table oversea_pay_event.
var overseaPayEventColumns = OverseaPayEventColumns{
	Id:              "id",
	BizType:         "biz_type",
	BizId:           "biz_id",
	Fee:             "fee",
	EventType:       "event_type",
	Event:           "event",
	RelativeTradeNo: "relative_trade_no",
	UniqueNo:        "unique_no",
	GmtCreate:       "gmt_create",
	GmtModify:       "gmt_modify",
	OpenApiId:       "open_api_id",
	TerminalIp:      "terminal_ip",
	Message:         "message",
}

// NewOverseaPayEventDao creates and returns a new DAO object for table data access.
func NewOverseaPayEventDao() *OverseaPayEventDao {
	return &OverseaPayEventDao{
		group:   "oversea_pay",
		table:   "oversea_pay_event",
		columns: overseaPayEventColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *OverseaPayEventDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *OverseaPayEventDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *OverseaPayEventDao) Columns() OverseaPayEventColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *OverseaPayEventDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *OverseaPayEventDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *OverseaPayEventDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
