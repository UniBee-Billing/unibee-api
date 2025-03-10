package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/vat_gateway"

	"unibee/api/merchant/vat"
)

func (c *ControllerVat) CountryList(ctx context.Context, req *vat.CountryListReq) (res *vat.CountryListRes, err error) {
	list, _ := vat_gateway.MerchantCountryRateList(ctx, _interface.GetMerchantId(ctx))
	return &vat.CountryListRes{VatCountryList: list}, nil
}
