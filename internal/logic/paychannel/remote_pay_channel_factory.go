package paychannel

import "go-oversea-pay/internal/logic/paychannel/impl"

func GetPayChannelServiceProvider(channel int) (channelService RemotePayChannelService) {
	return &impl.Evonet{}
}
