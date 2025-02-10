package merchant

import (
	"context"
	gateway3 "unibee/api/bean/detail"
	"unibee/api/merchant/gateway"
	_interface "unibee/internal/interface/context"
	gateway2 "unibee/internal/logic/gateway/service"
)

func (c *ControllerGateway) Setup(ctx context.Context, req *gateway.SetupReq) (res *gateway.SetupRes, err error) {
	return &gateway.SetupRes{Gateway: gateway3.ConvertGatewayDetail(ctx, gateway2.SetupGateway(ctx, _interface.GetMerchantId(ctx), req.GatewayName, req.GatewayKey, req.GatewaySecret, req.SubGateway, req.DisplayName, req.GatewayIcons, req.Sort, req.CurrencyExchange))}, nil
}
