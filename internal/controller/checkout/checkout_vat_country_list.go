package checkout

import (
	"context"
	"unibee/internal/logic/vat_gateway"
	"unibee/utility"

	"unibee/api/checkout/vat"
)

func (c *ControllerVat) CountryList(ctx context.Context, req *vat.CountryListReq) (res *vat.CountryListRes, err error) {
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	list, _ := vat_gateway.MerchantCountryRateList(ctx, req.MerchantId)
	return &vat.CountryListRes{VatCountryList: list}, nil
}
