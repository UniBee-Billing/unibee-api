package merchant

import (
	"context"
	subscription2 "go-oversea-pay/api/user/subscription"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/subscription/service"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/merchant/subscription"
)

func (c *ControllerSubscription) SubscriptionUpdatePreview(ctx context.Context, req *subscription.SubscriptionUpdatePreviewReq) (res *subscription.SubscriptionUpdatePreviewRes, err error) {
	//Update 可以由 Admin 操作，service 层不做用户校验
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).Merchant != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).Merchant.Id > 0, "merchantUserId invalid")
	}
	update, err := service.SubscriptionUpdatePreview(ctx, &subscription2.SubscriptionUpdatePreviewReq{
		SubscriptionId: req.SubscriptionId,
		NewPlanId:      req.NewPlanId,
		Quantity:       req.Quantity,
		AddonParams:    req.AddonParams,
	}, 0, int64(_interface.BizCtx().Get(ctx).Merchant.Id))
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionUpdatePreviewRes{
		TotalAmount:   update.TotalAmount,
		Currency:      update.Currency,
		Invoice:       update.Invoice,
		ProrationDate: update.ProrationDate,
	}, nil
}
