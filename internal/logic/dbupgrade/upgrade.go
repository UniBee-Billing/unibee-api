package dbupgrade

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"io"
	"net/http"
	"time"
	"unibee/api/bean"
	"unibee/internal/cmd/config"
	_ "unibee/internal/dao/default"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
	"unibee/utility/liberr"
)

func StandAloneInit(ctx context.Context) {
	list := fetchColumnAppendListFromCloudApi()
	if len(list) > 0 {
		glog.Infof(ctx, "StandAloneInit DBUpgrade start")
		historyList := GetUpgradeList(ctx)
		historyIds := make([]uint64, 0)
		for _, history := range historyList {
			historyIds = append(historyIds, history.Id)
		}
		database, err := gdb.Instance()
		tables, err := database.Tables(ctx, database.GetSchema())
		liberr.ErrIsNil(ctx, err, "DB Not Ready For Upgrade")
		utility.AssertError(err, "StandAloneInit DBUpgrade Get Database Instance failure,%v")
		for _, one := range list {
			if utility.IsUint64InArray(historyIds, one.Id) {
				continue
			}
			if database != nil && len(one.UpgradeSql) > 0 {
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
		glog.Infof(ctx, "StandAloneInit DBUpgrade end")
	}
}

func fetchColumnAppendListFromCloudApi() []*bean.TableUpgradeSimplify {
	var list = make([]*bean.TableUpgradeSimplify, 0)
	var env = 1
	if config.GetConfigInstance().IsProd() {
		env = 2
	}
	response, err := sendRequestInMainCtxStart(fmt.Sprintf("https://api.cloud.unibee.top/cloud/table/column_append?databaseType=%s&env=%v", g.DB("default").GetConfig().Type, env), "GET", nil, nil)
	if err != nil {
		return list
	}
	data := gjson.New(response)
	if data.Contains("code") && data.Get("code").String() == "0" && data.Contains("data") && data.GetJson("data").Contains("tableUpgrades") {
		_ = gjson.Unmarshal([]byte(data.GetJson("data").Get("tableUpgrades").String()), &list)
	}
	return list
}

func sendRequestInMainCtxStart(url string, method string, data []byte, headers map[string]string) ([]byte, error) {
	bodyReader := bytes.NewReader(data)
	request, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, gerror.NewCode(gcode.New(response.StatusCode, response.Status, response.Status+" "+string(responseBody)), response.Status+" "+string(responseBody))
	}
	return responseBody, nil
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

func SaveUpgradeHistory(ctx context.Context, one *bean.TableUpgradeSimplify) {
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
