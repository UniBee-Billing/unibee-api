// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TableUpgradeHistoryDao is the data access object for table table_upgrade_history.
type TableUpgradeHistoryDao struct {
	table   string                     // table is the underlying table name of the DAO.
	group   string                     // group is the database configuration group name of current DAO.
	columns TableUpgradeHistoryColumns // columns contains all the column names of Table for convenient usage.
}

// TableUpgradeHistoryColumns defines and stores column names for table table_upgrade_history.
type TableUpgradeHistoryColumns struct {
	Id            string // id
	DatabaseType  string // type of database
	Env           string // 0-offline,1-stage,2-prod
	Action        string // action
	TableName     string // table_name
	ColumnName    string // column_name
	ServerVersion string // server_version
	UpgradeSql    string // upgrade_sql
	GmtCreate     string // create time
	GmtModify     string // update time
	CreateTime    string // create utc time
}

// tableUpgradeHistoryColumns holds the columns for table table_upgrade_history.
var tableUpgradeHistoryColumns = TableUpgradeHistoryColumns{
	Id:            "id",
	DatabaseType:  "database_type",
	Env:           "env",
	Action:        "action",
	TableName:     "table_name",
	ColumnName:    "column_name",
	ServerVersion: "server_version",
	UpgradeSql:    "upgrade_sql",
	GmtCreate:     "gmt_create",
	GmtModify:     "gmt_modify",
	CreateTime:    "create_time",
}

// NewTableUpgradeHistoryDao creates and returns a new DAO object for table data access.
func NewTableUpgradeHistoryDao() *TableUpgradeHistoryDao {
	return &TableUpgradeHistoryDao{
		group:   "default",
		table:   "table_upgrade_history",
		columns: tableUpgradeHistoryColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TableUpgradeHistoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TableUpgradeHistoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TableUpgradeHistoryDao) Columns() TableUpgradeHistoryColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *TableUpgradeHistoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TableUpgradeHistoryDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TableUpgradeHistoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
