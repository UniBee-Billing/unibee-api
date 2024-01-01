package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/subscription/service"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/subscription/v1"
)

func (c *ControllerV1) SubscriptionPlanChannelTransferAndActivate(ctx context.Context, req *v1.SubscriptionPlanChannelTransferAndActivateReq) (res *v1.SubscriptionPlanChannelTransferAndActivateRes, err error) {
	utility.Assert(req.PlanId > 0, "plan should > 0")
	//utility.Assert(req.ChannelId > 0, "ConfirmChannelId should > 0")
	//多个渠道Plan 创建并激活
	list := query.GetListSubscriptionTypePayChannels(ctx) // todo mark 需改造成获取 merchantId 相关的 Channel
	utility.Assert(len(list) > 0, "no channel found, need at least one")
	for _, channel := range list {
		err = service.SubscriptionPlanChannelTransferAndActivate(ctx, req.PlanId, int64(channel.Id))
		if err != nil {
			utility.FailureJsonExit(g.RequestFromCtx(ctx), fmt.Sprintf("%s", err))
			return
		}
	}
	//发布 Plan
	err = service.SubscriptionPlanActivate(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}

	utility.SuccessJsonExit(g.RequestFromCtx(ctx), nil)
	return
}
