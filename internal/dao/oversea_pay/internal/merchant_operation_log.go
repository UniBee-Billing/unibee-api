// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantOperationLogDao is the data access object for table merchant_operation_log.
type MerchantOperationLogDao struct {
	table   string                      // table is the underlying table name of the DAO.
	group   string                      // group is the database configuration group name of current DAO.
	columns MerchantOperationLogColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantOperationLogColumns defines and stores column names for table merchant_operation_log.
type MerchantOperationLogColumns struct {
	Id                 string // 主键id
	CompanyId          string // 公司ID
	MerchantId         string // merchantId
	UserId             string // 操作userId，系统自动操作可能没有
	OptAccount         string // 操作账号
	ClientType         string // 操作渠道 0:云店后台 1:云管家app 2:Java服务 3:小程序
	BizType            string // 操作业务 0:菜单 1:商品 2:门店 3:订单 4:账号|会员 5:优惠券转赠中 6:优惠券转赠领取成功 7:优惠券转赠自动取消
	OptTarget          string // 操作对象
	OptContent         string // 操作内容
	OptCreate          string // 操作发生时间
	IsDelete           string // 0-UnDeleted，1-Deleted
	GmtCreate          string // create time
	GmtModify          string // 修改时间
	QueryportRequestId string // queryport请求Id，可在request_security_log查询请求信息
	ServerType         string // 操作终端，参看 message-api包 OperationLogServerTypeEnum的code
	ServerTypeDesc     string // 操作终端描述，参看 message-api包 OperationLogServerTypeEnum的desc
}

// merchantOperationLogColumns holds the columns for table merchant_operation_log.
var merchantOperationLogColumns = MerchantOperationLogColumns{
	Id:                 "id",
	CompanyId:          "company_id",
	MerchantId:         "merchant_id",
	UserId:             "user_id",
	OptAccount:         "opt_account",
	ClientType:         "client_type",
	BizType:            "biz_type",
	OptTarget:          "opt_target",
	OptContent:         "opt_content",
	OptCreate:          "opt_create",
	IsDelete:           "is_delete",
	GmtCreate:          "gmt_create",
	GmtModify:          "gmt_modify",
	QueryportRequestId: "queryport_request_id",
	ServerType:         "server_type",
	ServerTypeDesc:     "server_type_desc",
}

// NewMerchantOperationLogDao creates and returns a new DAO object for table data access.
func NewMerchantOperationLogDao() *MerchantOperationLogDao {
	return &MerchantOperationLogDao{
		group:   "oversea_pay",
		table:   "merchant_operation_log",
		columns: merchantOperationLogColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantOperationLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantOperationLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantOperationLogDao) Columns() MerchantOperationLogColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantOperationLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantOperationLogDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantOperationLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
