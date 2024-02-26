package system

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/system/subscription"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/subscription/handler"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
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
				WhereNotNull(dao.Subscription.Columns().GatewaySubscriptionId).
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
				gatewayPlan := query.GetGatewayPlan(backgroundCtx, one.PlanId, one.GatewayId)
				utility.Assert(gatewayPlan != nil, "invalid gatewayPlan")
				details, err := api.GetGatewayServiceProvider(backgroundCtx, one.GatewayId).GatewaySubscriptionDetails(backgroundCtx, plan, gatewayPlan, one)
				if err == nil {
					err := handler.UpdateSubWithGatewayDetailBack(backgroundCtx, one, details)
					if err != nil {
						fmt.Printf("BulkChannelSync Background UpdateSubWithGatewayDetailBack SubscriptionId:%s error%s\n", one.SubscriptionId, err.Error())
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
