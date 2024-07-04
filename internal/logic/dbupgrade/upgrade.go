package dbupgrade

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"unibee/api/bean"
	"unibee/internal/cmd/config"
	_ "unibee/internal/dao/oversea_pay"
	"unibee/utility"
)

func StandAloneInit(ctx context.Context) {
	list := fetchColumnAppendListFromCloudApi()
	if len(list) > 0 {
		glog.Infof(ctx, "StandAloneInit DBUpgrade start")
		database := g.DB("default")
		db, err := gdb.Instance()
		tables, err := database.Tables(ctx, database.GetSchema())
		if err != nil {
			glog.Errorf(ctx, "StandAloneInit DBUpgrade Get Tables error:%v", err.Error())
			return
		}
		utility.AssertError(err, "StandAloneInit DBUpgrade Get Database Instance failure,%v")
		for _, one := range list {
			if db != nil && len(one.UpgradeSql) > 0 {
				if len(one.Action) == 0 || len(one.TableName) == 0 {
					glog.Infof(ctx, "StandAloneInit DBUpgrade upgradeId:%v skip by empty action or tableName", one.Id)
					continue
				}
				if one.Action == "table_creation" {
					if !utility.IsStringInArray(tables, one.TableName) {
						_, err = db.Exec(ctx, one.UpgradeSql)
						if err != nil {
							glog.Errorf(ctx, "StandAloneInit DBUpgrade Create Table for upgradeId:%v error:%v", one.Id, err.Error())
						}
					}
				} else if one.Action == "column_add" {
					if utility.IsStringInArray(tables, one.TableName) {
						if len(one.ColumnName) == 0 {
							glog.Infof(ctx, "StandAloneInit DBUpgrade upgradeId:%v skip by empty columnName", one.Id)
							continue
						}
						fields, err := database.TableFields(ctx, one.TableName, database.GetSchema())
						if err != nil {
							glog.Errorf(ctx, "StandAloneInit DBUpgrade Get Table: %s Fields error:%v", one.TableName, err.Error())
							continue
						}
						if _, ok := fields[one.ColumnName]; !ok {
							_, err = db.Exec(ctx, one.UpgradeSql)
							if err != nil {
								glog.Errorf(ctx, "StandAloneInit DBUpgrade Append Table Column %s for upgradeId:%v error:%v", one.ColumnName, one.Id, err.Error())
							}
						}
					}
				} else if one.Action == "column_alter" {
					if utility.IsStringInArray(tables, one.TableName) {
						if len(one.ColumnName) == 0 {
							glog.Infof(ctx, "StandAloneInit DBUpgrade upgradeId:%v skip by empty columnName", one.Id)
							continue
						}
						fields, err := database.TableFields(ctx, one.TableName, database.GetSchema())
						if err != nil {
							glog.Errorf(ctx, "StandAloneInit DBUpgrade Get Table: %s Fields error:%v", one.TableName, err.Error())
							continue
						}
						if _, ok := fields[one.ColumnName]; ok {
							_, err = db.Exec(ctx, one.UpgradeSql)
							if err != nil {
								glog.Errorf(ctx, "StandAloneInit DBUpgrade Edit Table Column %s for upgradeId:%v error:%v", one.ColumnName, one.Id, err.Error())
							}
						}
					}
				}
			}
		}
		glog.Infof(ctx, "StandAloneInit DBUpgrade end")
	}
}

func fetchColumnAppendListFromCloudApi() []*bean.TableUpgradeSimplify {
	var list = make([]*bean.TableUpgradeSimplify, 0)
	var env = 1
	if config.GetConfigInstance().IsProd() {
		env = 2
	}
	response, err := utility.SendRequest(fmt.Sprintf("https://api.cloud.unibee.top/cloud/table/column_append?databaseType=%s&env=%v", g.DB("default").GetConfig().Type, env), "GET", nil, nil)
	if err != nil {
		return list
	}
	data := gjson.New(response)
	if data.Contains("code") && data.Get("code").String() == "0" && data.Contains("data") && data.GetJson("data").Contains("tableUpgrades") {
		_ = gjson.Unmarshal([]byte(data.GetJson("data").Get("tableUpgrades").String()), &list)
	}
	return list
}
