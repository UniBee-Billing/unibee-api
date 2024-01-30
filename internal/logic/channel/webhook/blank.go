package webhook

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/channel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type BlankWebhook struct {
}

func (b BlankWebhook) DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.MerchantChannelConfig) (err error) {
	//TODO implement me
	panic("implement me")
}

func (b BlankWebhook) DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.MerchantChannelConfig) {
	//TODO implement me
	panic("implement me")
}

func (b BlankWebhook) DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.MerchantChannelConfig) (res *ro.ChannelRedirectInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}
