package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/vat_gateway"
	"unibee/utility"

	"unibee/api/merchant/vat"
)

func (c *ControllerVat) SetupVatGateway(ctx context.Context, req *vat.SetupVatGatewayReq) (res *vat.SetupVatGatewayRes, err error) {
	//Admin 操作，service 层不做用户校验
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	err = vat_gateway.SetupMerchantVatConfig(ctx, _interface.GetMerchantId(ctx), req.VatName, req.VatData, req.IsDefault)
	if err != nil {
		return nil, err
	}
	if req.IsDefault {
		err := vat_gateway.InitMerchantDefaultVatGateway(ctx, _interface.GetMerchantId(ctx))
		if err != nil {
			return nil, err
		}
	}
	return &vat.SetupVatGatewayRes{}, nil
}
