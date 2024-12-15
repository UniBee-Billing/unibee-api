package dbupgrade

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	_ "unibee/internal/dao/default"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/platform"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
	"unibee/utility/liberr"
)

func StandAloneInit(ctx context.Context) {
	list := platform.FetchColumnAppendListFromPlatformApi()
	if len(list) > 0 {
		glog.Infof(ctx, "StandAloneInit DBUpgrade start")
		historyList := GetUpgradeList(ctx)
		historyIds := make([]uint64, 0)
		for _, history := range historyList {
			historyIds = append(historyIds, history.Id)
		}
		database, err := gdb.Instance()
		utility.AssertError(err, "StandAloneInit DBUpgrade Get Database Instance failure,%v")
		tables, err := database.Tables(ctx, database.GetSchema())
		liberr.ErrIsNil(ctx, err, "DB Not Ready For Upgrade")
		utility.AssertError(err, "StandAloneInit DBUpgrade Get Database Instance failure,%v")
		var needClearTableCache = false
		for _, one := range list {
			if utility.IsUint64InArray(historyIds, one.Id) {
				continue
			}
			if len(one.UpgradeSql) > 0 {
				if len(one.Action) == 0 || len(one.TableName) == 0 {
					glog.Infof(ctx, "StandAloneInit DBUpgrade upgradeId:%v skip by empty action or tableName", one.Id)
					continue
				}
				if one.Action == "table_creation" {
					if !utility.IsStringInArray(tables, one.TableName) {
						_, err = database.Exec(ctx, one.UpgradeSql)
						if err != nil {
							glog.Errorf(ctx, "StandAloneInit DBUpgrade Create Table for upgradeId:%v error:%v", one.Id, err.Error())
						} else {
							SaveUpgradeHistory(ctx, one)
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
							_, err = database.Exec(ctx, one.UpgradeSql)
							if err != nil {
								glog.Errorf(ctx, "StandAloneInit DBUpgrade Append Table Column %s for upgradeId:%v error:%v", one.ColumnName, one.Id, err.Error())
							} else {
								SaveUpgradeHistory(ctx, one)
								needClearTableCache = true
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
							_, err = database.Exec(ctx, one.UpgradeSql)
							if err != nil {
								glog.Errorf(ctx, "StandAloneInit DBUpgrade Edit Table Column %s for upgradeId:%v error:%v", one.ColumnName, one.Id, err.Error())
							} else {
								SaveUpgradeHistory(ctx, one)
								needClearTableCache = true
							}
						}
					}
				} else if one.Action == "index_add" {
					if utility.IsStringInArray(tables, one.TableName) {
						fields, err := database.TableFields(ctx, one.TableName, database.GetSchema())
						if err != nil {
							glog.Errorf(ctx, "StandAloneInit DBUpgrade Get Table: %s Fields error:%v", one.TableName, err.Error())
							continue
						}
						if _, ok := fields[one.ColumnName]; ok {
							_, err = database.Exec(ctx, one.UpgradeSql)
							if err != nil {
								glog.Errorf(ctx, "StandAloneInit DBUpgrade Add Table Key %s for upgradeId:%v error:%v", one.ColumnName, one.Id, err.Error())
							} else {
								SaveUpgradeHistory(ctx, one)
							}
						}
					}
				}
			}
		}
		if needClearTableCache {
			err = g.DB().GetCore().ClearTableFieldsAll(ctx)
			if err != nil {
				glog.Errorf(ctx, "StandAloneInit ClearTableFieldsAll error:%v", err.Error())
			}
			err = g.DB().GetCore().ClearCacheAll(ctx)
			if err != nil {
				glog.Errorf(ctx, "StandAloneInit ClearCacheAll error:%v", err.Error())
			}
		}
		glog.Infof(ctx, "StandAloneInit DBUpgrade end")
	}
}

func GetUpgradeList(ctx context.Context) (list []*entity.TableUpgradeHistory) {
	var data = make([]*entity.TableUpgradeHistory, 0)
	err := dao.TableUpgradeHistory.Ctx(ctx).Scan(&data)
	if err != nil {
		g.Log().Errorf(ctx, "GetUpgradeList error:%s", err)
		return nil
	}
	return data
}

func SaveUpgradeHistory(ctx context.Context, one *bean.TableUpgrade) {
	g.Log().Info(ctx, "StandAloneInit DBUpgrade success and save upgradeId:%v", one.Id)
	_, err := dao.TableUpgradeHistory.Ctx(ctx).Data(&entity.TableUpgradeHistory{
		Id:            one.Id,
		DatabaseType:  one.DatabaseType,
		Env:           one.Env,
		Action:        one.Action,
		TableName:     one.TableName,
		ColumnName:    one.ColumnName,
		ServerVersion: g.Server().GetOpenApi().Info.Version,
		UpgradeSql:    one.UpgradeSql,
		GmtCreate:     gtime.Now(),
		GmtModify:     gtime.Now(),
		CreateTime:    gtime.Now().Timestamp(),
	}).Insert()
	if err != nil {
		g.Log().Info(ctx, "StandAloneInit DBUpgrade save upgradeId:%v error:%v", one.Id, err.Error())
	}
}
