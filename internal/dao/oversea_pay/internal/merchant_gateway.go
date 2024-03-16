// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantGatewayDao is the data access object for table merchant_gateway.
type MerchantGatewayDao struct {
	table   string                 // table is the underlying table name of the DAO.
	group   string                 // group is the database configuration group name of current DAO.
	columns MerchantGatewayColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantGatewayColumns defines and stores column names for table merchant_gateway.
type MerchantGatewayColumns struct {
	Id                    string // gateway_id
	MerchantId            string // merchant_id
	EnumKey               string // enum key , match in gateway implementation
	GatewayType           string // gateway type，1-Default｜ 2-Crypto
	GatewayName           string // gateway name
	Name                  string // name
	SubGateway            string // sub_gateway_enum
	BrandData             string //
	Logo                  string // gateway logo
	Host                  string // pay host
	GatewayAccountId      string // gateway account id
	GatewayKey            string //
	GatewaySecret         string // secret
	Custom                string // custom
	GmtCreate             string // create time
	GmtModify             string // update time
	Description           string // description
	WebhookKey            string // webhook_key
	WebhookSecret         string // webhook_secret
	UniqueProductId       string // unique  gateway productId, only stripe need
	CreateTime            string // create utc time
	IsDeleted             string // 0-UnDeleted，1-Deleted
	CryptoReceiveCurrency string //
	CountryConfig         string //
}

// merchantGatewayColumns holds the columns for table merchant_gateway.
var merchantGatewayColumns = MerchantGatewayColumns{
	Id:                    "id",
	MerchantId:            "merchant_id",
	EnumKey:               "enum_key",
	GatewayType:           "gateway_type",
	GatewayName:           "gateway_name",
	Name:                  "name",
	SubGateway:            "sub_gateway",
	BrandData:             "brand_data",
	Logo:                  "logo",
	Host:                  "host",
	GatewayAccountId:      "gateway_account_id",
	GatewayKey:            "gateway_key",
	GatewaySecret:         "gateway_secret",
	Custom:                "custom",
	GmtCreate:             "gmt_create",
	GmtModify:             "gmt_modify",
	Description:           "description",
	WebhookKey:            "webhook_key",
	WebhookSecret:         "webhook_secret",
	UniqueProductId:       "unique_product_id",
	CreateTime:            "create_time",
	IsDeleted:             "is_deleted",
	CryptoReceiveCurrency: "crypto_receive_currency",
	CountryConfig:         "country_config",
}

// NewMerchantGatewayDao creates and returns a new DAO object for table data access.
func NewMerchantGatewayDao() *MerchantGatewayDao {
	return &MerchantGatewayDao{
		group:   "default",
		table:   "merchant_gateway",
		columns: merchantGatewayColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantGatewayDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantGatewayDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantGatewayDao) Columns() MerchantGatewayColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantGatewayDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantGatewayDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantGatewayDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
