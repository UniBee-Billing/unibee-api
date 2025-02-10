package merchant

import (
	"context"
	_interface2 "unibee/internal/interface"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/gateway/service"
	"unibee/internal/logic/merchant_config"
	"unibee/internal/logic/merchant_config/update"
	"unibee/utility"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) EditSort(ctx context.Context, req *gateway.EditSortReq) (res *gateway.EditSortRes, err error) {
	utility.Assert(req.GatewaySorts != nil, "Invalid Sort")
	sortConfig := merchant_config.GetMerchantConfig(ctx, _interface.GetMerchantId(ctx), _interface2.KEY_MERCHANT_GATEWAY_SORT)
	if sortConfig == nil {
		var data map[string]int64
		for _, v := range req.GatewaySorts {
			if v.Id > 0 {
				service.UpdateGatewaySort(ctx, _interface.GetMerchantId(ctx), v.Id, v.Sort)
			} else {
				data[v.GatewayName] = v.Sort
			}
		}
		err = update.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), _interface2.KEY_MERCHANT_GATEWAY_SORT, utility.MarshalToJsonString(data))
		utility.AssertError(err, "Update sort failed")
	} else {
		var data map[string]int64
		err = utility.UnmarshalFromJsonString(sortConfig.ConfigValue, &data)
		utility.AssertError(err, "Update sort failed")
		for _, v := range req.GatewaySorts {
			if v.Id > 0 {
				service.UpdateGatewaySort(ctx, _interface.GetMerchantId(ctx), v.Id, v.Sort)
			} else {
				data[v.GatewayName] = v.Sort
			}
		}
		err = update.SetMerchantConfig(ctx, _interface.GetMerchantId(ctx), _interface2.KEY_MERCHANT_GATEWAY_SORT, utility.MarshalToJsonString(data))
		utility.AssertError(err, "Update sort failed")
	}
	listRes, err := c.SetupList(ctx, &gateway.SetupListReq{})
	if err != nil {
		return nil, err
	}
	return &gateway.EditSortRes{Gateways: listRes.Gateways}, nil
}
