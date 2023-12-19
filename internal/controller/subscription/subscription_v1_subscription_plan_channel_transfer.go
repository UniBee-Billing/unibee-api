package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/subscription/service"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/subscription/v1"
)

func (c *ControllerV1) SubscriptionPlanChannelTransfer(ctx context.Context, req *v1.SubscriptionPlanChannelTransferReq) (res *v1.SubscriptionPlanChannelTransferRes, err error) {
	utility.Assert(req.PlanId > 0, "plan should > 0")
	utility.Assert(req.ChannelId > 0, "ChannelId should > 0")
	err = service.SubscriptionPlanChannelTransfer(ctx, req.PlanId, req.ChannelId)
	if err != nil {
		utility.FailureJsonExit(g.RequestFromCtx(ctx), fmt.Sprintf("%s", err))
		return
	}
	utility.SuccessJsonExit(g.RequestFromCtx(ctx), nil)
	return
}
