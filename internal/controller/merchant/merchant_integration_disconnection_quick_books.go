package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/analysis/quickbooks"
	"unibee/internal/logic/merchant_config/update"
	"unibee/utility"

	"unibee/api/merchant/integration"
)

func (c *ControllerIntegration) DisconnectionQuickBooks(ctx context.Context, req *integration.DisconnectionQuickBooksReq) (res *integration.DisconnectionQuickBooksRes, err error) {
	_ = update.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), quickbooks.KeyMerchantQuickBooksConfig, utility.MarshalToJsonString(&quickbooks.MerchantQuickBooksConfig{}))
	return &integration.DisconnectionQuickBooksRes{}, nil
}
