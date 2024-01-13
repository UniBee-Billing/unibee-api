// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CountryRateDao is the data access object for table country_rate.
type CountryRateDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns CountryRateColumns // columns contains all the column names of Table for convenient usage.
}

// CountryRateColumns defines and stores column names for table country_rate.
type CountryRateColumns struct {
	Id                    string // id
	VatName               string // vat channel name, em vatsense
	CountryCode           string // country_code
	CountryName           string // country_name
	Latitude              string // latitude
	Longitude             string // longitude
	Vat                   string // 是否有 vat，1-有，2-没有
	Eu                    string // 是否欧盟成员国, 1-是，2-不是
	StandardTaxPercentage string // Standard Tax百分比，10 表示 10%
	Other                 string // other rates(json格式)
	StandardDescription   string // standard_description
	StandardTypes         string // standard_typs 限定
	Provinces             string // Whether the country has provinces with provincial sales tax
	Mamo                  string // 备注
	GmtCreate             string // 创建时间
	GmtModify             string // 更新时间
	IsDeleted             string // 是否删除，0-未删除，1-删除
}

// countryRateColumns holds the columns for table country_rate.
var countryRateColumns = CountryRateColumns{
	Id:                    "id",
	VatName:               "vat_name",
	CountryCode:           "country_code",
	CountryName:           "country_name",
	Latitude:              "latitude",
	Longitude:             "longitude",
	Vat:                   "vat",
	Eu:                    "eu",
	StandardTaxPercentage: "standard_tax_percentage",
	Other:                 "other",
	StandardDescription:   "standard_description",
	StandardTypes:         "standard_types",
	Provinces:             "provinces",
	Mamo:                  "mamo",
	GmtCreate:             "gmt_create",
	GmtModify:             "gmt_modify",
	IsDeleted:             "is_deleted",
}

// NewCountryRateDao creates and returns a new DAO object for table data access.
func NewCountryRateDao() *CountryRateDao {
	return &CountryRateDao{
		group:   "oversea_pay",
		table:   "country_rate",
		columns: countryRateColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CountryRateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CountryRateDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CountryRateDao) Columns() CountryRateColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CountryRateDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CountryRateDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CountryRateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
