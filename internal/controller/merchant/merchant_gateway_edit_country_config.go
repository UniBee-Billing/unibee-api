package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/service"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) EditCountryConfig(ctx context.Context, req *gateway.EditCountryConfigReq) (res *gateway.EditCountryConfigRes, err error) {
	err = service.EditGatewayCountryConfig(ctx, _interface.GetMerchantId(ctx), req.GatewayId, req.CountryConfig)
	if err != nil {
		return nil, err
	}
	return &gateway.EditCountryConfigRes{}, nil
}
