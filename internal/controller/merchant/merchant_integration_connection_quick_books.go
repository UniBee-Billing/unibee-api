package merchant

import (
	"context"
	"fmt"
	"unibee/internal/cmd/config"
	"unibee/internal/controller/link"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/analysis/quickbooks"
	"unibee/internal/logic/analysis/quickbooks/quickbooksdk"
	"unibee/internal/logic/merchant_config/update"
	"unibee/utility"

	"unibee/api/merchant/integration"
)

func (c *ControllerIntegration) ConnectionQuickBooks(ctx context.Context, req *integration.ConnectionQuickBooksReq) (res *integration.ConnectionQuickBooksRes, err error) {
	apiKeys, err := quickbooks.GetCloudQuickBooksPartnerAPIKeys(ctx)
	utility.AssertError(err, "Get QuickBooks APIKeys Error")
	qbClient, err := quickbooksdk.NewClient(apiKeys.ClientId, apiKeys.ClientSecret, "", config.GetConfigInstance().IsProd(), "", nil)
	utility.AssertError(err, "Connection QuickBooks Client Error")
	url, err := qbClient.FindAuthorizationUrl("com.intuit.quickbooks.accounting", fmt.Sprintf("%d", _interface.GetMerchantId(ctx)), link.GetQuickbooksAuthorizationLink())
	utility.AssertError(err, "Get QuickBooks Authorization URL Error")
	if len(req.ReturnUrl) > 0 {
		_ = update.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), quickbooks.KeyMerchantQuickBooksConfig, utility.MarshalToJsonString(&quickbooks.MerchantQuickBooksConfig{
			SetupReturnUrl: req.ReturnUrl,
		}))
	}
	return &integration.ConnectionQuickBooksRes{AuthorizationURL: url}, nil
}
