package gateway_log

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"time"
	dao "unibee/internal/dao/oversea_pay"
)

func TaskForDeleteChannelLogs(ctx context.Context) {
	g.Log().Infof(ctx, "TaskForDeleteChannelLogs start")
	time.Sleep(5 * time.Second)
	_, err := dao.GatewayHttpLog.Ctx(ctx).WhereLT(dao.GatewayHttpLog.Columns().GmtCreate, gtime.Now().AddDate(0, 0, -7)).Delete()
	if err != nil {
		g.Log().Errorf(ctx, "TaskForDeleteChannelLogs error:%s", err.Error())
	}
}
