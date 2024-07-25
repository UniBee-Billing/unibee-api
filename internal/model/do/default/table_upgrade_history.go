// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TableUpgradeHistory is the golang structure of table table_upgrade_history for DAO operations like Where/Data.
type TableUpgradeHistory struct {
	g.Meta        `orm:"table:table_upgrade_history, do:true"`
	Id            interface{} // id
	DatabaseType  interface{} // type of database
	Env           interface{} // 0-offline,1-stage,2-prod
	Action        interface{} // action
	TableName     interface{} // table_name
	ColumnName    interface{} // column_name
	ServerVersion interface{} // server_version
	UpgradeSql    interface{} // upgrade_sql
	GmtCreate     *gtime.Time // create time
	GmtModify     *gtime.Time // update time
	CreateTime    interface{} // create utc time
}
