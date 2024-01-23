package system

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/api/system/subscription"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/gateway"
	"go-oversea-pay/internal/logic/subscription/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func (c *ControllerSubscription) BulkChannelSync(ctx context.Context, req *subscription.BulkChannelSyncReq) (*subscription.BulkChannelSyncRes, error) {
	utility.Assert(len(req.MerchantId) > 0, "merchantId invalid")
	go func() {
		defer func() {
			if exception := recover(); exception != nil {
				var err error
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				g.Log().Errorf(context.Background(), "BulkChannelSync Background panic error:%s\n", err.Error())
				return
			}
		}()
		var page = 0
		var count = 100
		for {
			backgroundCtx := context.Background()
			var mainList []*entity.Subscription
			err := dao.Subscription.Ctx(backgroundCtx).
				Where(dao.Subscription.Columns().MerchantId, req.MerchantId).
				WhereNotNull(dao.Subscription.Columns().ChannelSubscriptionId).
				OrderDesc("id").
				Limit(page*count, count).
				OmitEmpty().Scan(&mainList)
			if err != nil {
				fmt.Printf("BulkChannelSync Background List error%s\n", err.Error())
				return
			}
			for _, one := range mainList {
				plan := query.GetPlanById(backgroundCtx, one.PlanId)
				utility.Assert(plan != nil, "invalid planId")
				utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
				planChannel := query.GetPlanChannel(backgroundCtx, one.PlanId, one.ChannelId)
				utility.Assert(planChannel != nil, "invalid planChannel")
				details, err := gateway.GetPayChannelServiceProvider(backgroundCtx, one.ChannelId).DoRemoteChannelSubscriptionDetails(backgroundCtx, plan, planChannel, one)
				if err == nil {
					err := handler.UpdateSubWithChannelDetailBack(backgroundCtx, one, details)
					if err != nil {
						fmt.Printf("BulkChannelSync Background UpdateSubWithChannelDetailBack SubscriptionId:%s error%s\n", one.SubscriptionId, err.Error())
						return
					}
					fmt.Printf("BulkChannelSync Background Fetch SubscriptionId:%s success\n", one.SubscriptionId)
				} else {
					fmt.Printf("BulkChannelSync Background Fetch SubscriptionId:%s error%s\n", one.SubscriptionId, err.Error())
				}
			}
			if len(mainList) == 0 {
				break
			}
			clear(mainList)
			page = page + 1
		}
	}()
	return nil, nil
}
