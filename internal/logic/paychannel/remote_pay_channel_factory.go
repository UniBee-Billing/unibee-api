package paychannel

func GetPayChannelServiceProvider(channel int) (channelService RemotePayChannelService) {
	proxy := &PayChannelProxy{}
	proxy.channel = channel
	return proxy
}
