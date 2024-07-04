package dbupgrade

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/glog"
	"unibee/api/bean"
	"unibee/internal/cmd/config"
	"unibee/utility"
)

func StandAloneInit(ctx context.Context) {
	list := fetchColumnAppendListFromCloudApi()
	if len(list) > 0 {
		glog.Infof(ctx, "StandAloneInit DBUpgrade start")
		db, err := gdb.Instance()
		utility.AssertError(err, "StandAloneInit DBUpgrade Get Database Instance failure,%v")
		for _, one := range list {
			glog.Infof(ctx, "StandAloneInit DBUpgrade get,%v", utility.MarshalToJsonString(one))

			if db != nil && len(one.UpgradeSql) > 0 {
				// todo mark check or upgrade
				_, err = db.Exec(ctx, one.UpgradeSql)
				if err != nil {
					glog.Errorf(ctx, "StandAloneInit DBUpgrade for upgradeId:%v error:%v", one.Id, err.Error())
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
	// todo mark check databaseType
	response, err := utility.SendRequest(fmt.Sprintf("https://api.cloud.unibee.top/cloud/table/column_append?databaseType=mysql&env=%v", env), "GET", nil, nil)
	if err != nil {
		return list
	}
	data := gjson.New(response)
	if data.Contains("code") && data.Get("code").String() == "0" && data.Contains("data") && data.GetJson("data").Contains("tableUpgrades") {
		_ = gjson.Unmarshal([]byte(data.GetJson("data").Get("tableUpgrades").String()), &list)
	}
	return list
}
