// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProductDao is the data access object for table product.
type ProductDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns ProductColumns // columns contains all the column names of Table for convenient usage.
}

// ProductColumns defines and stores column names for table product.
type ProductColumns struct {
	Id          string //
	GmtCreate   string // create time
	GmtModify   string // update time
	CompanyId   string // company id
	MerchantId  string // merchant id
	ProductName string // ProductName
	Description string // description
	ImageUrl    string // image_url
	HomeUrl     string // home_url
	Status      string // status，1-active，2-inactive, default active
	IsDeleted   string // 0-UnDeleted，1-Deleted
	CreateTime  string // create utc time
	MetaData    string // meta_data(json)
}

// productColumns holds the columns for table product.
var productColumns = ProductColumns{
	Id:          "id",
	GmtCreate:   "gmt_create",
	GmtModify:   "gmt_modify",
	CompanyId:   "company_id",
	MerchantId:  "merchant_id",
	ProductName: "product_name",
	Description: "description",
	ImageUrl:    "image_url",
	HomeUrl:     "home_url",
	Status:      "status",
	IsDeleted:   "is_deleted",
	CreateTime:  "create_time",
	MetaData:    "meta_data",
}

// NewProductDao creates and returns a new DAO object for table data access.
func NewProductDao() *ProductDao {
	return &ProductDao{
		group:   "default",
		table:   "product",
		columns: productColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ProductDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ProductDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ProductDao) Columns() ProductColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ProductDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ProductDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ProductDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
