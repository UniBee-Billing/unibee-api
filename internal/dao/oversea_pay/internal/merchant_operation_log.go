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
	Id                 string // id
	CompanyId          string // company id
	MerchantId         string // merchant Id
	MemberId           string // member_id
	UserId             string // user_id
	OptAccount         string // admin account
	ClientType         string // client type
	BizType            string // biz_type
	OptTarget          string // operation target
	OptContent         string // operation content
	CreateTime         string // operation create utc time
	IsDelete           string // 0-UnDeleted，1-Deleted
	GmtCreate          string // create time
	GmtModify          string // update time
	QueryportRequestId string // queryport id
	ServerType         string // server type
	ServerTypeDesc     string // server type description
	SubscriptionId     string // subscription_id
	InvoiceId          string // invoice id
	PlanId             string // plan id
	DiscountCode       string // discount_code
}

// merchantOperationLogColumns holds the columns for table merchant_operation_log.
var merchantOperationLogColumns = MerchantOperationLogColumns{
	Id:                 "id",
	CompanyId:          "company_id",
	MerchantId:         "merchant_id",
	MemberId:           "member_id",
	UserId:             "user_id",
	OptAccount:         "opt_account",
	ClientType:         "client_type",
	BizType:            "biz_type",
	OptTarget:          "opt_target",
	OptContent:         "opt_content",
	CreateTime:         "create_time",
	IsDelete:           "is_delete",
	GmtCreate:          "gmt_create",
	GmtModify:          "gmt_modify",
	QueryportRequestId: "queryport_request_id",
	ServerType:         "server_type",
	ServerTypeDesc:     "server_type_desc",
	SubscriptionId:     "subscription_id",
	InvoiceId:          "invoice_id",
	PlanId:             "plan_id",
	DiscountCode:       "discount_code",
}

// NewMerchantOperationLogDao creates and returns a new DAO object for table data access.
func NewMerchantOperationLogDao() *MerchantOperationLogDao {
	return &MerchantOperationLogDao{
		group:   "default",
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
