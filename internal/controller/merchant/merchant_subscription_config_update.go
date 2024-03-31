package merchant

import (
	"context"
	"fmt"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/merchant_config"
	"unibee/internal/logic/subscription/config"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) ConfigUpdate(ctx context.Context, req *subscription.ConfigUpdateReq) (res *subscription.ConfigUpdateRes, err error) {
	if req.DowngradeEffectImmediately != nil {
		err = merchant_config.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), config.DowngradeEffectImmediately, fmt.Sprintf("%v", *req.DowngradeEffectImmediately))
		if err != nil {
			return nil, err
		}
	}
	if req.UpgradeProration != nil {
		err = merchant_config.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), config.UpdateProration, fmt.Sprintf("%v", *req.UpgradeProration))
		if err != nil {
			return nil, err
		}
	}
	if req.IncompleteExpireTime != nil {
		err = merchant_config.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), config.IncompleteExpireTime, fmt.Sprintf("%v", *req.IncompleteExpireTime))
		if err != nil {
			return nil, err
		}
	}
	if req.InvoiceEmail != nil {
		err = merchant_config.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), config.InvoiceEmail, fmt.Sprintf("%v", *req.InvoiceEmail))
		if err != nil {
			return nil, err
		}
	}
	if req.TryAutomaticPaymentBeforePeriodEnd != nil {
		err = merchant_config.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), config.TryAutomaticPaymentBeforePeriodEnd, fmt.Sprintf("%v", *req.TryAutomaticPaymentBeforePeriodEnd))
		if err != nil {
			return nil, err
		}
	}
	return &subscription.ConfigUpdateRes{Config: config.GetMerchantSubscriptionConfig(ctx, _interface.GetMerchantId(ctx))}, nil
}
