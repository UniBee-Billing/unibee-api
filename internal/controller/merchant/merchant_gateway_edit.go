package merchant

import (
	"context"
	gateway3 "unibee/api/bean/detail"
	"unibee/api/merchant/gateway"
	_interface "unibee/internal/interface/context"
	gateway2 "unibee/internal/logic/gateway/service"
)

func (c *ControllerGateway) Edit(ctx context.Context, req *gateway.EditReq) (res *gateway.EditRes, err error) {
	return &gateway.EditRes{Gateway: gateway3.ConvertGatewayDetail(ctx, gateway2.EditGateway(ctx, _interface.GetMerchantId(ctx), req.GatewayId, req.GatewayKey, req.GatewaySecret, req.SubGateway, req.DisplayName, req.GatewayLogo, req.Sort, req.CurrencyExchange))}, nil
}
