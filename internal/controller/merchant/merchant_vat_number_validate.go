package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/vat_gateway"
	"unibee/utility"

	"unibee/api/merchant/vat"
)

func (c *ControllerVat) NumberValidate(ctx context.Context, req *vat.NumberValidateReq) (res *vat.NumberValidateRes, err error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(len(req.VatNumber) > 0, "vatNumber invalid")
	vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, _interface.GetMerchantId(ctx), 0, req.VatNumber, "")
	if err != nil {
		return nil, err
	}
	return &vat.NumberValidateRes{VatNumberValidate: vatNumberValidate}, nil
}
