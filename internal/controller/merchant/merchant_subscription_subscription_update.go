package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"sync"
	"time"
	"unibee/api/merchant/subscription"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerSubscription) Update(ctx context.Context, req *subscription.UpdateReq) (res *subscription.UpdateRes, err error) {
	if len(req.SubscriptionId) == 0 {
		utility.Assert(req.UserId > 0, "one of SubscriptionId and UserId should provide")
		utility.Assert(req.NewPlanId > 0, "newPlanId should provide while SubscriptionId is blank")
		plan := query.GetPlanById(ctx, req.NewPlanId)
		utility.Assert(plan != nil, fmt.Sprintf("plan not found:%v", req.NewPlanId))
		one := query.GetLatestActiveOrIncompleteSubscriptionByUserId(ctx, req.UserId, _interface.GetMerchantId(ctx), plan.ProductId)
		utility.Assert(one != nil, "no active or incomplete subscription found")
		req.SubscriptionId = one.SubscriptionId
	} else {
		one := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
		utility.Assert(one != nil, "sub not found")
		utility.Assert(one.MerchantId == _interface.GetMerchantId(ctx), "merchantId not match")
	}
	var memberMemberId int64 = 0
	if _interface.Context().Get(ctx) != nil && _interface.Context().Get(ctx).MerchantMember != nil {
		memberMemberId = int64(_interface.Context().Get(ctx).MerchantMember.Id)
	}

	if req.Discount != nil {
		utility.Assert(_interface.Context().Get(ctx).IsOpenApiCall, "Discount only available for api call")
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
				g.Log().Errorf(taskCtx, "MerchantSubscriptionUpdateSubmit Panic Error:%s", backgroundErr.Error())
				err = backgroundErr
				return
			}
		}()

		update, err = service.SubscriptionUpdate(taskCtx, &service.UpdateInternalReq{
			SubscriptionId:         req.SubscriptionId,
			NewPlanId:              req.NewPlanId,
			Quantity:               req.Quantity,
			AddonParams:            req.AddonParams,
			EffectImmediate:        req.EffectImmediate,
			GatewayId:              req.GatewayId,
			GatewayPaymentType:     req.GatewayPaymentType,
			ConfirmTotalAmount:     req.ConfirmTotalAmount,
			ConfirmCurrency:        req.ConfirmCurrency,
			ProrationDate:          req.ProrationDate,
			Metadata:               req.Metadata,
			DiscountCode:           req.DiscountCode,
			TaxPercentage:          req.TaxPercentage,
			Discount:               req.Discount,
			ManualPayment:          req.ManualPayment,
			ReturnUrl:              req.ReturnUrl,
			CancelUrl:              req.CancelUrl,
			ProductData:            req.ProductData,
			ApplyPromoCredit:       req.ApplyPromoCredit,
			ApplyPromoCreditAmount: req.ApplyPromoCreditAmount,
		}, memberMemberId)
		operation_log.AppendOptLog(taskCtx, &operation_log.OptLogRequest{
			MerchantId:     update.SubscriptionPendingUpdate.MerchantId,
			Target:         fmt.Sprintf("Subscription(%v)", update.SubscriptionPendingUpdate.SubscriptionId),
			Content:        "Update",
			UserId:         update.SubscriptionPendingUpdate.UserId,
			SubscriptionId: update.SubscriptionPendingUpdate.SubscriptionId,
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)
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
	}, nil
}
