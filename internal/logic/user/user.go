package user

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "go-oversea-pay/internal/dao/oversea_pay"
)

func UpdateUserDefaultSubscription(ctx context.Context, userId int64, subscriptionId string) {
	_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
		dao.UserAccount.Columns().SubscriptionId: subscriptionId,
		dao.UserAccount.Columns().GmtModify:      gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, userId).OmitNil().Update()
	if err != nil {
		g.Log().Errorf(ctx, "UpdateUserDefaultSubscription err:%s", err.Error())
	}
}
