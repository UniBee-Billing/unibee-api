package outchannel

func GetPayChannelServiceProvider(channel int) (channelService RemotePayChannelInterface) {
	proxy := &PayChannelProxy{}
	proxy.channel = channel
	return proxy
}
