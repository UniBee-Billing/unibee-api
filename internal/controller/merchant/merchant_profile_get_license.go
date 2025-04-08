package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/profile"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/middleware"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerProfile) GetLicense(ctx context.Context, req *profile.GetLicenseReq) (res *profile.GetLicenseRes, err error) {
	merchant := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	utility.Assert(merchant != nil, "merchant not found")
	return &profile.GetLicenseRes{
		Merchant:     bean.SimplifyMerchant(merchant),
		License:      middleware.GetMerchantLicense(ctx, merchant.Id),
		APIRateLimit: middleware.GetMerchantAPIRateLimit(ctx, merchant.Id),
	}, nil
}
