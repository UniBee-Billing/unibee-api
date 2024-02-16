// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantVatNumberVerifyHistoryDao is the data access object for table merchant_vat_number_verify_history.
type MerchantVatNumberVerifyHistoryDao struct {
	table   string                                // table is the underlying table name of the DAO.
	group   string                                // group is the database configuration group name of current DAO.
	columns MerchantVatNumberVerifyHistoryColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantVatNumberVerifyHistoryColumns defines and stores column names for table merchant_vat_number_verify_history.
type MerchantVatNumberVerifyHistoryColumns struct {
	Id              string // Id
	MerchantId      string // merchantId
	VatNumber       string // vat_number
	Valid           string // 0-Invalid，1-Valid
	ValidateGateway string // validate_gateway
	CountryCode     string // country_code
	CompanyName     string // company_name
	CompanyAddress  string // company_address
	GmtCreate       string // create time
	GmtModify       string // update time
	IsDeleted       string // 0-UnDeleted，1-Deleted
	ValidateMessage string // validate_message
	CreateTime      string // create utc time
}

// merchantVatNumberVerifyHistoryColumns holds the columns for table merchant_vat_number_verify_history.
var merchantVatNumberVerifyHistoryColumns = MerchantVatNumberVerifyHistoryColumns{
	Id:              "id",
	MerchantId:      "merchant_id",
	VatNumber:       "vat_number",
	Valid:           "valid",
	ValidateGateway: "validate_gateway",
	CountryCode:     "country_code",
	CompanyName:     "company_name",
	CompanyAddress:  "company_address",
	GmtCreate:       "gmt_create",
	GmtModify:       "gmt_modify",
	IsDeleted:       "is_deleted",
	ValidateMessage: "validate_message",
	CreateTime:      "create_time",
}

// NewMerchantVatNumberVerifyHistoryDao creates and returns a new DAO object for table data access.
func NewMerchantVatNumberVerifyHistoryDao() *MerchantVatNumberVerifyHistoryDao {
	return &MerchantVatNumberVerifyHistoryDao{
		group:   "oversea_pay",
		table:   "merchant_vat_number_verify_history",
		columns: merchantVatNumberVerifyHistoryColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantVatNumberVerifyHistoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantVatNumberVerifyHistoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantVatNumberVerifyHistoryDao) Columns() MerchantVatNumberVerifyHistoryColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantVatNumberVerifyHistoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantVatNumberVerifyHistoryDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantVatNumberVerifyHistoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
