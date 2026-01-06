package gateway_log

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"time"
	dao "unibee/internal/dao/default"
)

func TaskForDeleteChannelLogs(ctx context.Context) {
	g.Log().Infof(ctx, "TaskForDeleteChannelLogs start")
	time.Sleep(5 * time.Second)
	_, err := dao.GatewayHttpLog.Ctx(ctx).WhereLT(dao.GatewayHttpLog.Columns().GmtCreate, gtime.Now().AddDate(0, 0, -15)).Delete()
	if err != nil {
		g.Log().Errorf(ctx, "TaskForDeleteChannelLogs error:%s", err.Error())
	}
}

func TaskForDeleteWebhookMessage(ctx context.Context) {
	g.Log().Infof(ctx, "TaskForDeleteWebhookMessage start")
	time.Sleep(5 * time.Second)
	_, err := dao.MerchantWebhookMessage.Ctx(ctx).
		WhereLT(dao.MerchantWebhookMessage.Columns().WebsocketStatus, 50).
		WhereLT(dao.MerchantWebhookMessage.Columns().GmtCreate, gtime.Now().AddDate(0, 0, -15)).
		Delete()
	if err != nil {
		g.Log().Errorf(ctx, "TaskForDeleteWebhookMessage error:%s", err.Error())
	}
}

func TaskForDeleteWebhookLog(ctx context.Context) {
	g.Log().Infof(ctx, "TaskForDeleteWebhookLog start")
	time.Sleep(5 * time.Second)
	_, err := dao.MerchantWebhookLog.Ctx(ctx).WhereLT(dao.MerchantWebhookLog.Columns().GmtCreate, gtime.Now().AddDate(0, 0, -60)).Delete()
	if err != nil {
		g.Log().Errorf(ctx, "TaskForDeleteWebhookLog error:%s", err.Error())
	}
}

func TaskForDeleteOperationLog(ctx context.Context) {
	g.Log().Infof(ctx, "TaskForDeleteOperationLog start")
	time.Sleep(5 * time.Second)
	_, err := dao.MerchantOperationLog.Ctx(ctx).WhereLT(dao.MerchantOperationLog.Columns().GmtCreate, gtime.Now().AddDate(0, 0, -90)).Delete()
	if err != nil {
		g.Log().Errorf(ctx, "TaskForDeleteOperationLog error:%s", err.Error())
	}
}
