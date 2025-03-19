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

func (c *ControllerSubscription) Renew(ctx context.Context, req *subscription.RenewReq) (res *subscription.RenewRes, err error) {
	g.Log().Debugf(ctx, "Subscription Renew called by payload:%s\n", utility.MarshalToJsonString(req))
	if len(req.SubscriptionId) == 0 {
		utility.Assert(req.UserId > 0, "one of SubscriptionId and UserId should provide")
		one := query.GetLatestActiveOrIncompleteSubscriptionByUserId(ctx, req.UserId, _interface.GetMerchantId(ctx), req.ProductId)
		if one != nil {
			req.SubscriptionId = one.SubscriptionId
		} else {
			one = query.GetLatestSubscriptionByUserId(ctx, req.UserId, _interface.GetMerchantId(ctx), req.ProductId)
			utility.Assert(one != nil, "no subscription found")
			req.SubscriptionId = one.SubscriptionId
		}
	}
	if req.Discount != nil {
		utility.Assert(_interface.Context().Get(ctx).IsOpenApiCall, "Discount only available for api call")
	}
	lockKey := fmt.Sprintf("SubscriptionRenewProcess-%s", req.SubscriptionId)
	if !utility.TryLock(ctx, lockKey, 60) {
		utility.Assert(false, "Another subscription renew is in process")
	}
	var wg sync.WaitGroup
	wg.Add(1)
	var renewRes *service.CreateInternalRes
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
				g.Log().Errorf(taskCtx, "MerchantSubscriptionRenew Panic Error:%s", backgroundErr.Error())
				err = backgroundErr
				return
			}
		}()
		renewRes, err = service.SubscriptionRenew(ctx, &service.RenewInternalReq{
			MerchantId:             _interface.GetMerchantId(ctx),
			SubscriptionId:         req.SubscriptionId,
			GatewayId:              req.GatewayId,
			GatewayPaymentType:     req.GatewayPaymentType,
			TaxPercentage:          req.TaxPercentage,
			DiscountCode:           req.DiscountCode,
			Discount:               req.Discount,
			ManualPayment:          req.ManualPayment,
			ReturnUrl:              req.ReturnUrl,
			CancelUrl:              req.CancelUrl,
			ProductData:            req.ProductData,
			Metadata:               req.Metadata,
			ApplyPromoCredit:       req.ApplyPromoCredit,
			ApplyPromoCreditAmount: req.ApplyPromoCreditAmount,
		})
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     renewRes.Subscription.MerchantId,
			Target:         fmt.Sprintf("Subscription(%v)", renewRes.Subscription.SubscriptionId),
			Content:        "RenewUserSubscription",
			UserId:         renewRes.Subscription.UserId,
			SubscriptionId: renewRes.Subscription.SubscriptionId,
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
	if renewRes == nil {
		return nil, gerror.New("Server Error")
	}
	return &subscription.RenewRes{
		Subscription: renewRes.Subscription,
		Paid:         renewRes.Paid,
		Link:         renewRes.Link,
	}, nil
}
