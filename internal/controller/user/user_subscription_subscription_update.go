package user

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"sync"
	"time"
	"unibee/api/user/subscription"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerSubscription) Update(ctx context.Context, req *subscription.UpdateReq) (res *subscription.UpdateRes, err error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.NewPlanId > 0, "PlanId invalid")
	utility.Assert(len(req.SubscriptionId) > 0, "SubscriptionId invalid")
	utility.Assert(req.EffectImmediate == 0, "EffectImmediate not support in user_portal")
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	if !config.GetConfigInstance().IsLocal() {
		utility.Assert(_interface.Context().Get(ctx).User != nil, "auth failure,not login")
		utility.Assert(_interface.Context().Get(ctx).User.Id == sub.UserId, "userId not match")
	}
	lockKey := fmt.Sprintf("SubscriptionUpdateProcess-%s", req.SubscriptionId)
	if !utility.TryLock(ctx, lockKey, 60) {
		utility.Assert(false, "Another subscription update is in process")
	}
	var wg sync.WaitGroup
	wg.Add(1)
	var update *service.UpdateInternalRes
	go func() {
		defer wg.Done()
		taskCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		var backgroundErr error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					backgroundErr = v
				} else {
					backgroundErr = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				g.Log().Errorf(taskCtx, "UserSubscriptionUpdateSubmit Panic Error:%s", backgroundErr.Error())
				err = backgroundErr
				return
			}
		}()
		update, err = service.SubscriptionUpdate(ctx, &service.UpdateInternalReq{
			SubscriptionId:         req.SubscriptionId,
			NewPlanId:              req.NewPlanId,
			Quantity:               req.Quantity,
			GatewayId:              req.GatewayId,
			GatewayPaymentType:     req.GatewayPaymentType,
			AddonParams:            req.AddonParams,
			ConfirmTotalAmount:     req.ConfirmTotalAmount,
			ConfirmCurrency:        req.ConfirmCurrency,
			ProrationDate:          req.ProrationDate,
			EffectImmediate:        req.EffectImmediate,
			Metadata:               req.Metadata,
			DiscountCode:           req.DiscountCode,
			ApplyPromoCredit:       req.ApplyPromoCredit,
			ApplyPromoCreditAmount: req.ApplyPromoCreditAmount,
			ReturnUrl:              req.ReturnUrl,
			CancelUrl:              req.CancelUrl,
		}, 0)
	}()
	wg.Wait()
	utility.ReleaseLock(context.Background(), lockKey)
	if err != nil {
		return nil, err
	}
	if update == nil {
		return nil, gerror.New("Server Error")
	}
	return &subscription.UpdateRes{
		SubscriptionPendingUpdate: update.SubscriptionPendingUpdate,
		Paid:                      update.Paid,
		Link:                      update.Link,
		Note:                      update.Note,
	}, err
}
