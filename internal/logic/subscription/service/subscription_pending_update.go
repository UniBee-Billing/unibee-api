package service

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/payment/service"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
)

func SubscriptionPendingUpdateCancel(ctx context.Context, pendingUpdateId string, reason string) error {
	one := query.GetSubscriptionPendingUpdateByPendingUpdateId(ctx, pendingUpdateId)
	if one != nil {
		if one.Status == consts.PendingSubStatusFinished {
			return nil
		}
		if one.Status == consts.PendingSubStatusCancelled {
			return nil
		}

		if len(one.ChannelUpdateId) > 0 {
			payment := query.GetPaymentByPaymentId(ctx, one.ChannelUpdateId)
			if payment != nil {
				err := service.DoChannelCancel(ctx, payment)
				if err != nil {
					g.Log().Errorf(ctx, "DoChannelCancel Error:%s", err.Error())
				}
			}
		}
		_, err := dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
			dao.SubscriptionPendingUpdate.Columns().Status:    consts.PendingSubStatusCancelled,
			dao.SubscriptionPendingUpdate.Columns().GmtModify: gtime.Now(),
		}).Where(dao.SubscriptionPendingUpdate.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return err
		}
	}
	return nil
}

func CancelOtherUnfinishedPendingUpdatesBackground(subscriptionId string, pendingUpdateId string, reason string) {
	go func() {
		var err error
		ctx := context.Background()
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				g.Log().Errorf(ctx, "CancelOtherUnfinishedPendingUpdatesBackground Panic Error:%s", err.Error())
				return
			}
		}()
		var mainList []*entity.SubscriptionPendingUpdate
		err = dao.SubscriptionPendingUpdate.Ctx(ctx).
			Where(dao.SubscriptionPendingUpdate.Columns().SubscriptionId, subscriptionId).
			WhereNot(dao.SubscriptionPendingUpdate.Columns().UpdateSubscriptionId, pendingUpdateId).
			WhereLT(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusFinished).
			Limit(0, 100).
			OmitEmpty().Scan(&mainList)
		if err != nil {
			g.Log().Errorf(ctx, "CancelOtherUnfinishedPendingUpdatesBackground Search List Error:%s", err.Error())
		}
		for _, one := range mainList {
			err = SubscriptionPendingUpdateCancel(ctx, one.UpdateSubscriptionId, reason)
			if err != nil {
				g.Log().Errorf(ctx, "CancelOtherUnfinishedPendingUpdatesBackground SubscriptionPendingUpdateCancel Error:%s", err.Error())
			}
		}
	}()
}