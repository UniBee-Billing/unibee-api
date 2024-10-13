package merchant

import (
	"context"
	"fmt"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/merchant_config/update"
	"unibee/internal/logic/subscription/config"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) ConfigUpdate(ctx context.Context, req *subscription.ConfigUpdateReq) (res *subscription.ConfigUpdateRes, err error) {
	if req.DowngradeEffectImmediately != nil {
		err = update.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), config.DowngradeEffectImmediately, fmt.Sprintf("%v", *req.DowngradeEffectImmediately))
		if err != nil {
			return nil, err
		}
	}
	if req.UpgradeProration != nil {
		err = update.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), config.UpdateProration, fmt.Sprintf("%v", *req.UpgradeProration))
		if err != nil {
			return nil, err
		}
	}
	if req.IncompleteExpireTime != nil {
		utility.Assert(*req.IncompleteExpireTime >= 0, "Value should be greater than 0")
		err = update.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), config.IncompleteExpireTime, fmt.Sprintf("%v", *req.IncompleteExpireTime))
		if err != nil {
			return nil, err
		}
	}
	if req.InvoiceEmail != nil {
		err = update.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), config.InvoiceEmail, fmt.Sprintf("%v", *req.InvoiceEmail))
		if err != nil {
			return nil, err
		}
	}
	if req.TryAutomaticPaymentBeforePeriodEnd != nil {
		utility.Assert(*req.TryAutomaticPaymentBeforePeriodEnd >= 5*60, "Value should be greater than 300")
		err = update.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), config.TryAutomaticPaymentBeforePeriodEnd, fmt.Sprintf("%v", *req.TryAutomaticPaymentBeforePeriodEnd))
		if err != nil {
			return nil, err
		}
	}
	if len(req.GatewayVATRule) > 0 {
		err = update.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), config.GatewayVATRule, fmt.Sprintf("%s", utility.MarshalToJsonString(req.GatewayVATRule)))
		if err != nil {
			return nil, err
		}
	}
	if req.ShowZeroInvoice != nil {
		err = update.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), config.ShowZeroInvoice, fmt.Sprintf("%v", *req.ShowZeroInvoice))
		if err != nil {
			return nil, err
		}
	}
	return &subscription.ConfigUpdateRes{Config: config.GetMerchantSubscriptionConfig(ctx, _interface.GetMerchantId(ctx))}, nil
}
