package checkout

import (
	"context"
	"unibee/internal/logic/vat_gateway"
	"unibee/utility"

	"unibee/api/checkout/vat"
)

func (c *ControllerVat) NumberValidate(ctx context.Context, req *vat.NumberValidateReq) (res *vat.NumberValidateRes, err error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	utility.Assert(len(req.VatNumber) > 0, "vatNumber invalid")
	vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, req.MerchantId, 0, req.VatNumber, "")
	if err != nil {
		return nil, err
	}
	return &vat.NumberValidateRes{VatNumberValidate: vatNumberValidate}, nil
}
