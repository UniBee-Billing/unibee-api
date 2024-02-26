package user

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/vat_gateway"

	"unibee/api/user/vat"
)

func (c *ControllerVat) CountryVatList(ctx context.Context, req *vat.CountryVatListReq) (res *vat.CountryVatListRes, err error) {
	list, err := vat_gateway.MerchantCountryRateList(ctx, _interface.GetMerchantId(ctx))
	if err != nil {
		return nil, err
	}
	return &vat.CountryVatListRes{VatCountryList: list}, nil
}
