package service

import (
	"unibee-api/internal/logic/gateway/ro"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

func ConvertGatewayToRo(gateway *entity.MerchantGateway) *ro.OutGatewayRo {
	if gateway == nil {
		return nil
	}
	return &ro.OutGatewayRo{
		GatewayId:   gateway.Id,
		GatewayName: gateway.Name,
	}
}

func ConvertChannelsToRos(gateways []*entity.MerchantGateway) []*ro.OutGatewayRo {
	var outChannelRos []*ro.OutGatewayRo
	for _, gateway := range gateways {
		outChannelRos = append(outChannelRos, &ro.OutGatewayRo{
			GatewayId:   gateway.Id,
			GatewayName: gateway.Name,
		})
	}
	return outChannelRos
}
