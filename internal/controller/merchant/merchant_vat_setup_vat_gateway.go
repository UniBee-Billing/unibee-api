package merchant

import (
	"context"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/vat_gateway"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/merchant/vat"
)

func (c *ControllerVat) SetupVatGateway(ctx context.Context, req *vat.SetupVatGatewayReq) (res *vat.SetupVatGatewayRes, err error) {
	//Admin 操作，service 层不做用户校验
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).Merchant != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).Merchant.Id > 0, "merchantUserId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).Merchant.MerchantId == uint64(req.MerchantId), "token not match MerchantId invalid")
	}
	err = vat_gateway.SetupMerchantVatConfig(ctx, req.MerchantId, req.VatName, req.VatData, req.IsDefault)
	if err != nil {
		return nil, err
	}
	if req.IsDefault {
		vat_gateway.InitMerchantDefaultVatGateway(ctx, req.MerchantId)
	}
	return &vat.SetupVatGatewayRes{}, nil
}
