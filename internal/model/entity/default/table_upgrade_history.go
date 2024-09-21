// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TableUpgradeHistory is the golang structure for table table_upgrade_history.
type TableUpgradeHistory struct {
	Id            uint64      `json:"id"            description:"id"`                       // id
	DatabaseType  string      `json:"databaseType"  description:"type of database"`         // type of database
	Env           int         `json:"env"           description:"0-offline,1-stage,2-prod"` // 0-offline,1-stage,2-prod
	Action        string      `json:"action"        description:"action"`                   // action
	TableName     string      `json:"tableName"     description:"table_name"`               // table_name
	ColumnName    string      `json:"columnName"    description:"column_name"`              // column_name
	ServerVersion string      `json:"serverVersion" description:"server_version"`           // server_version
	UpgradeSql    string      `json:"upgradeSql"    description:"upgrade_sql"`              // upgrade_sql
	GmtCreate     *gtime.Time `json:"gmtCreate"     description:"create time"`              // create time
	GmtModify     *gtime.Time `json:"gmtModify"     description:"update time"`              // update time
	CreateTime    int64       `json:"createTime"    description:"create utc time"`          // create utc time
}
