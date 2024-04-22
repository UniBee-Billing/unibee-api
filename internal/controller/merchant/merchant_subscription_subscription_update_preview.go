package merchant

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) UpdatePreview(ctx context.Context, req *subscription.UpdatePreviewReq) (res *subscription.UpdatePreviewRes, err error) {
	//Update 可以由 Admin 操作，service 层不做用户校验
	var merchantMemberId int64
	if !config.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.Context().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
		merchantMemberId = int64(_interface.Context().Get(ctx).MerchantMember.Id)
	}
	g.Log().Infof(ctx, "SubscriptionUpdatePreview merchantMemberId:%d", merchantMemberId)
	update, err := service.SubscriptionUpdatePreview(ctx, &service.UpdatePreviewInternalReq{
		SubscriptionId:  req.SubscriptionId,
		NewPlanId:       req.NewPlanId,
		Quantity:        req.Quantity,
		AddonParams:     req.AddonParams,
		EffectImmediate: req.EffectImmediate,
		DiscountCode:    req.DiscountCode,
	}, 0, merchantMemberId)
	if err != nil {
		return nil, err
	}
	return &subscription.UpdatePreviewRes{
		OriginAmount:      update.OriginAmount,
		TotalAmount:       update.TotalAmount,
		DiscountAmount:    update.DiscountAmount,
		Currency:          update.Currency,
		Invoice:           update.Invoice,
		NextPeriodInvoice: update.NextPeriodInvoice,
		ProrationDate:     update.ProrationDate,
		Discount:          update.Discount,
	}, nil
}
