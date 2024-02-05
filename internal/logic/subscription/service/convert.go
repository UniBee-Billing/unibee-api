package service

import (
	"go-oversea-pay/internal/logic/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func ConvertChannelToRo(channel *entity.MerchantGateway) *ro.OutChannelRo {
	if channel == nil {
		return nil
	}
	return &ro.OutChannelRo{
		ChannelId:   channel.Id,
		ChannelName: channel.Name,
	}
}

func ConvertChannelsToRos(channels []*entity.MerchantGateway) []*ro.OutChannelRo {
	var outChannelRos []*ro.OutChannelRo
	for _, channel := range channels {
		outChannelRos = append(outChannelRos, &ro.OutChannelRo{
			ChannelId:   channel.Id,
			ChannelName: channel.Name,
		})
	}
	return outChannelRos
}
