package service

import (
	"go-oversea-pay/internal/logic/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func ConvertGatewayToRo(channel *entity.MerchantGateway) *ro.OutGatewayRo {
	if channel == nil {
		return nil
	}
	return &ro.OutGatewayRo{
		GatewayId:   channel.Id,
		GatewayName: channel.Name,
	}
}

func ConvertChannelsToRos(channels []*entity.MerchantGateway) []*ro.OutGatewayRo {
	var outChannelRos []*ro.OutGatewayRo
	for _, channel := range channels {
		outChannelRos = append(outChannelRos, &ro.OutGatewayRo{
			GatewayId:   channel.Id,
			GatewayName: channel.Name,
		})
	}
	return outChannelRos
}
