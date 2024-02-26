package merchant

import (
	"context"
	"unibee/internal/logic/vat_gateway"

	"unibee/api/merchant/vat"
)

func (c *ControllerVat) CountryVatList(ctx context.Context, req *vat.CountryVatListReq) (res *vat.CountryVatListRes, err error) {
	list, err := vat_gateway.MerchantCountryRateList(ctx, req.MerchantId)
	if err != nil {
		return nil, err
	}
	return &vat.CountryVatListRes{VatCountryList: list}, nil
}
