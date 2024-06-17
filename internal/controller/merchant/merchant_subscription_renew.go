package merchant

import (
	"context"
	"fmt"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) Renew(ctx context.Context, req *subscription.RenewReq) (res *subscription.RenewRes, err error) {
	if len(req.SubscriptionId) == 0 {
		utility.Assert(req.UserId > 0, "one of SubscriptionId and UserId should provide")
		one := query.GetLatestActiveOrIncompleteSubscriptionByUserId(ctx, req.UserId, _interface.GetMerchantId(ctx))
		if one != nil {
			req.SubscriptionId = one.SubscriptionId
		} else {
			one = query.GetLatestSubscriptionByUserId(ctx, req.UserId, _interface.GetMerchantId(ctx))
			utility.Assert(one != nil, "no subscription found")
			req.SubscriptionId = one.SubscriptionId
		}
	}
	renewRes, err := service.SubscriptionRenew(ctx, &service.RenewInternalReq{
		MerchantId:     _interface.GetMerchantId(ctx),
		SubscriptionId: req.SubscriptionId,
		GatewayId:      req.GatewayId,
		TaxPercentage:  req.TaxPercentage,
		DiscountCode:   req.DiscountCode,
		Discount:       req.Discount,
		ManualPayment:  req.ManualPayment,
		ReturnUrl:      req.ReturnUrl,
		CancelUrl:      req.CancelUrl,
		ProductData:    req.ProductData,
		Metadata:       req.Metadata,
	})
	if err != nil {
		return nil, err
	}
	if err == nil {
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
	}
	return &subscription.RenewRes{
		Subscription: renewRes.Subscription,
		Paid:         renewRes.Paid,
		Link:         renewRes.Link,
	}, nil
}
