package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/email"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/ro"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func HandlePendingUpdatePaymentFailure(ctx context.Context, pendingUpdateId string) (bool, error) {
	one := query.GetSubscriptionPendingUpdateByPendingUpdateId(ctx, pendingUpdateId)
	if one == nil {
		return false, gerror.New("FinishPendingUpdateForSubscription PendingUpdate Not Found:" + one.UpdateSubscriptionId)
	}
	if one.Status == consts.PendingSubStatusFinished {
		return true, nil
	}
	if one.Status == consts.PendingSubStatusCancelled {
		return true, nil
	}
	_, err := dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPendingUpdate.Columns().Status:    consts.PendingSubStatusCancelled,
		dao.SubscriptionPendingUpdate.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionPendingUpdate.Columns().Id, one.Id).Where(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusCreate).OmitNil().Update()
	if err != nil {
		return false, err
	}
	return true, nil
}

func FinishPendingUpdateForSubscription(ctx context.Context, sub *entity.Subscription, pendingUpdateId string) (bool, error) {
	one := query.GetSubscriptionPendingUpdateByPendingUpdateId(ctx, pendingUpdateId)
	utility.Assert(one != nil, "FinishPendingUpdateForSubscription PendingUpdate Not Found:"+pendingUpdateId)
	if one.Status == consts.PendingSubStatusFinished {
		return true, nil
	}
	utility.Assert(one.Status == consts.PendingSubStatusCreate, "pendingUpdate not status create")
	if consts.ProrationUsingUniBeeCompute && one.EffectImmediate == 1 && sub.Type == consts.SubTypeDefault {
		var addonParams []*ro.SubscriptionPlanAddonParamRo
		err := utility.UnmarshalFromJsonString(one.UpdateAddonData, &addonParams)
		if err != nil {
			return false, err
		}
		_, err = api.GetGatewayServiceProvider(ctx, one.GatewayId).GatewaySubscriptionUpdate(ctx, &ro.GatewayUpdateSubscriptionInternalReq{
			Plan:            query.GetPlanById(ctx, one.UpdatePlanId),
			Quantity:        one.UpdateQuantity,
			AddonPlans:      checkAndListAddonsFromParams(ctx, addonParams, one.GatewayId),
			GatewayPlan:     query.GetGatewayPlan(ctx, one.UpdatePlanId, one.GatewayId),
			Subscription:    query.GetSubscriptionBySubscriptionId(ctx, one.SubscriptionId),
			ProrationDate:   one.ProrationDate,
			EffectImmediate: false,
		})
		if err != nil {
			return false, err
		}
	}

	// 先创建 SubscriptionTimeLine 在做 Sub 更新
	err := CreateOrUpdateSubscriptionTimeline(ctx, sub, fmt.Sprintf("pendingUpdateFinish-%s", one.UpdateSubscriptionId))
	if err != nil {
		g.Log().Errorf(ctx, "CreateOrUpdateSubscriptionTimeline error:%s", err.Error())
	}
	// todo mark use transaction
	_, err = dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPendingUpdate.Columns().Status:    consts.PendingSubStatusFinished,
		dao.SubscriptionPendingUpdate.Columns().GmtModify: gtime.Now(),
	}).Where(dao.SubscriptionPendingUpdate.Columns().Id, one.Id).Where(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusCreate).OmitNil().Update()
	if err != nil {
		return false, err
	}

	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().PlanId:          one.UpdatePlanId,
		dao.Subscription.Columns().Quantity:        one.UpdateQuantity,
		dao.Subscription.Columns().AddonData:       one.UpdateAddonData,
		dao.Subscription.Columns().Amount:          one.UpdateAmount,
		dao.Subscription.Columns().Currency:        one.UpdateCurrency,
		dao.Subscription.Columns().LastUpdateTime:  gtime.Now().Timestamp(),
		dao.Subscription.Columns().GmtModify:       gtime.Now(),
		dao.Subscription.Columns().PendingUpdateId: "", //clear PendingUpdateId
	}).Where(dao.Subscription.Columns().SubscriptionId, one.SubscriptionId).OmitNil().Update()
	if err != nil {
		return false, err
	}

	user := query.GetUserAccountById(ctx, uint64(sub.UserId))
	merchant := query.GetMerchantInfoById(ctx, sub.MerchantId)
	err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, email.TemplateSubscriptionUpdate, "", &email.TemplateVariable{
		UserName:            user.FirstName + " " + user.LastName,
		MerchantProductName: query.GetPlanById(ctx, one.UpdatePlanId).GatewayProductName,
		MerchantCustomEmail: merchant.Email,
		MerchantName:        merchant.Name,
		PeriodEnd:           gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
	})
	if err != nil {
		fmt.Printf("SendTemplateEmail err:%s", err.Error())
	}
	return true, nil
}
