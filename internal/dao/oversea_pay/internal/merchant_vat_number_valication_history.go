// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantVatNumberValicationHistoryDao is the data access object for table merchant_vat_number_valication_history.
type MerchantVatNumberValicationHistoryDao struct {
	table   string                                    // table is the underlying table name of the DAO.
	group   string                                    // group is the database configuration group name of current DAO.
	columns MerchantVatNumberValicationHistoryColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantVatNumberValicationHistoryColumns defines and stores column names for table merchant_vat_number_valication_history.
type MerchantVatNumberValicationHistoryColumns struct {
	Id              string // ID
	MerchantId      string // merchantId
	VatNumber       string // vat_number
	Valid           string // 0-无效，1-有效
	ValidateChannel string // validate_channel
	CountryCode     string // country_code
	CompanyName     string // company_name
	CompanyAddress  string // company_address
	GmtCreate       string // 创建时间
	GmtModify       string // 修改时间
	IsDeleted       string // 0-UnDeleted，1-Deleted
	ValidateMessage string // validate_message
}

// merchantVatNumberValicationHistoryColumns holds the columns for table merchant_vat_number_valication_history.
var merchantVatNumberValicationHistoryColumns = MerchantVatNumberValicationHistoryColumns{
	Id:              "id",
	MerchantId:      "merchant_id",
	VatNumber:       "vat_number",
	Valid:           "valid",
	ValidateChannel: "validate_channel",
	CountryCode:     "country_code",
	CompanyName:     "company_name",
	CompanyAddress:  "company_address",
	GmtCreate:       "gmt_create",
	GmtModify:       "gmt_modify",
	IsDeleted:       "is_deleted",
	ValidateMessage: "validate_message",
}

// NewMerchantVatNumberValicationHistoryDao creates and returns a new DAO object for table data access.
func NewMerchantVatNumberValicationHistoryDao() *MerchantVatNumberValicationHistoryDao {
	return &MerchantVatNumberValicationHistoryDao{
		group:   "oversea_pay",
		table:   "merchant_vat_number_valication_history",
		columns: merchantVatNumberValicationHistoryColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantVatNumberValicationHistoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantVatNumberValicationHistoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantVatNumberValicationHistoryDao) Columns() MerchantVatNumberValicationHistoryColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantVatNumberValicationHistoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantVatNumberValicationHistoryDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantVatNumberValicationHistoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
