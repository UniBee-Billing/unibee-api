package user

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/vat_gateway"

	"unibee/api/user/vat"
)

func (c *ControllerVat) CountryList(ctx context.Context, req *vat.CountryListReq) (res *vat.CountryListRes, err error) {
	list, err := vat_gateway.MerchantCountryRateList(ctx, _interface.GetMerchantId(ctx))
	if err != nil {
		return nil, err
	}
	return &vat.CountryListRes{VatCountryList: list}, nil
}
