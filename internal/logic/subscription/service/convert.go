package service

import (
	"go-oversea-pay/internal/logic/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
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
