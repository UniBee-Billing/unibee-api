// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CreditConfigDao is the data access object for table credit_config.
type CreditConfigDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns CreditConfigColumns // columns contains all the column names of Table for convenient usage.
}

// CreditConfigColumns defines and stores column names for table credit_config.
type CreditConfigColumns struct {
	Id                    string // Id
	Type                  string // type of credit account, 1-main account, 2-promo credit account
	Currency              string // currency
	ExchangeRate          string // keep two decimal places，multiply by 100 saved, 1 currency = 1 credit * (exchange_rate/100), main account fixed rate to 100
	GmtCreate             string // create time
	GmtModify             string // update time
	CreateTime            string // create utc time
	MerchantId            string // merchant id
	Recurring             string // apply to recurring, default no, 0-no,1-yes
	DiscountCodeExclusive string // discount code exclusive when purchase, default no, 0-no, 1-yes
	Logo                  string // logo image base64, show when user purchase
	Name                  string // name
	Description           string // description
	LogoUrl               string // logo url, show when user purchase
	MetaData              string // meta_data(json)
	RechargeEnable        string // 0-yes, 1-no
	PayoutEnable          string // 0-yes, 1-no
	PreviewDefaultUsed    string // is default used when in purchase preview,0-no, 1-yes
}

// creditConfigColumns holds the columns for table credit_config.
var creditConfigColumns = CreditConfigColumns{
	Id:                    "id",
	Type:                  "type",
	Currency:              "currency",
	ExchangeRate:          "exchange_rate",
	GmtCreate:             "gmt_create",
	GmtModify:             "gmt_modify",
	CreateTime:            "create_time",
	MerchantId:            "merchant_id",
	Recurring:             "recurring",
	DiscountCodeExclusive: "discount_code_exclusive",
	Logo:                  "logo",
	Name:                  "name",
	Description:           "description",
	LogoUrl:               "logo_url",
	MetaData:              "meta_data",
	RechargeEnable:        "recharge_enable",
	PayoutEnable:          "payout_enable",
	PreviewDefaultUsed:    "preview_default_used",
}

// NewCreditConfigDao creates and returns a new DAO object for table data access.
func NewCreditConfigDao() *CreditConfigDao {
	return &CreditConfigDao{
		group:   "default",
		table:   "credit_config",
		columns: creditConfigColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CreditConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CreditConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CreditConfigDao) Columns() CreditConfigColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CreditConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CreditConfigDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CreditConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
