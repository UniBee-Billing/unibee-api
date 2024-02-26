package service

import (
	"unibee/internal/logic/gateway/ro"
	entity "unibee/internal/model/entity/oversea_pay"
)

func ConvertGatewayToRo(gateway *entity.MerchantGateway) *ro.GatewaySimplify {
	if gateway == nil {
		return nil
	}
	return &ro.GatewaySimplify{
		Id:          gateway.Id,
		GatewayName: gateway.Name,
	}
}

func ConvertChannelsToRos(gateways []*entity.MerchantGateway) []*ro.GatewaySimplify {
	var outChannelRos []*ro.GatewaySimplify
	for _, gateway := range gateways {
		outChannelRos = append(outChannelRos, &ro.GatewaySimplify{
			Id:          gateway.Id,
			GatewayName: gateway.Name,
		})
	}
	return outChannelRos
}
