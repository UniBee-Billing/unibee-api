package service

import (
	"go-oversea-pay/internal/logic/channel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func ConvertChannelToRo(channel *entity.MerchantChannelConfig) *ro.OutChannelRo {
	if channel == nil {
		return nil
	}
	return &ro.OutChannelRo{
		ChannelId:   channel.Id,
		ChannelName: channel.Name,
	}
}

func ConvertChannelsToRos(channels []*entity.MerchantChannelConfig) []*ro.OutChannelRo {
	var outChannelRos []*ro.OutChannelRo
	for _, channel := range channels {
		outChannelRos = append(outChannelRos, &ro.OutChannelRo{
			ChannelId:   channel.Id,
			ChannelName: channel.Name,
		})
	}
	return outChannelRos
}
