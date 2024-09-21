// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TableUpgrade is the golang structure of table table_upgrade for DAO operations like Where/Data.
type TableUpgrade struct {
	g.Meta       `orm:"table:table_upgrade, do:true"`
	Id           interface{} // id
	DatabaseType interface{} // type of database
	Env          interface{} // 0-offline,1-stage,2-prod
	Action       interface{} // action
	TableName    interface{} // table_name
	ColumnName   interface{} // column_name
	ColumnType   interface{} // column_type
	UpgradeSql   interface{} // upgrade_sql
	GmtCreate    *gtime.Time // create time
	GmtModify    *gtime.Time // update time
}
