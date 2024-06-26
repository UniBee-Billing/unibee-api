package merchant

import (
	"context"
	"fmt"
	"unibee/api/merchant/subscription"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerSubscription) Update(ctx context.Context, req *subscription.UpdateReq) (res *subscription.UpdateRes, err error) {
	if len(req.SubscriptionId) == 0 {
		utility.Assert(req.UserId > 0, "one of SubscriptionId and UserId should provide")
		one := query.GetLatestActiveOrIncompleteSubscriptionByUserId(ctx, req.UserId, _interface.GetMerchantId(ctx))
		utility.Assert(one != nil, "no active or incomplete subscription found")
		req.SubscriptionId = one.SubscriptionId
	}
	var memberMemberId int64 = 0
	if _interface.Context().Get(ctx).MerchantMember != nil {
		memberMemberId = int64(_interface.Context().Get(ctx).MerchantMember.Id)
	}
	update, err := service.SubscriptionUpdate(ctx, &service.UpdateInternalReq{
		SubscriptionId:     req.SubscriptionId,
		NewPlanId:          req.NewPlanId,
		Quantity:           req.Quantity,
		AddonParams:        req.AddonParams,
		EffectImmediate:    req.EffectImmediate,
		GatewayId:          req.GatewayId,
		ConfirmTotalAmount: req.ConfirmTotalAmount,
		ConfirmCurrency:    req.ConfirmCurrency,
		ProrationDate:      req.ProrationDate,
		Metadata:           req.Metadata,
		DiscountCode:       req.DiscountCode,
		TaxPercentage:      req.TaxPercentage,
		Discount:           req.Discount,
		ManualPayment:      req.ManualPayment,
		ReturnUrl:          req.ReturnUrl,
		CancelUrl:          req.CancelUrl,
		ProductData:        req.ProductData,
	}, memberMemberId)
	if err != nil {
		return nil, err
	}
	if err == nil {
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     update.SubscriptionPendingUpdate.MerchantId,
			Target:         fmt.Sprintf("Subscription(%v)", update.SubscriptionPendingUpdate.SubscriptionId),
			Content:        "Update",
			UserId:         update.SubscriptionPendingUpdate.UserId,
			SubscriptionId: update.SubscriptionPendingUpdate.SubscriptionId,
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)
	}
	return &subscription.UpdateRes{
		SubscriptionPendingUpdate: update.SubscriptionPendingUpdate,
		Paid:                      update.Paid,
		Link:                      update.Link,
	}, nil
}
