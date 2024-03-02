package merchant

import (
	"context"
	"unibee/internal/logic/vat_gateway"

	"unibee/api/merchant/vat"
)

func (c *ControllerVat) CountryList(ctx context.Context, req *vat.CountryListReq) (res *vat.CountryListRes, err error) {
	list, err := vat_gateway.MerchantCountryRateList(ctx, req.MerchantId)
	if err != nil {
		return nil, err
	}
	return &vat.CountryListRes{VatCountryList: list}, nil
}
