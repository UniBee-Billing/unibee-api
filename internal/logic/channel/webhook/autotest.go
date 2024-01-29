package webhook

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/channel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type AutoTestWebhook struct {
}

func (b AutoTestWebhook) DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.OverseaPayChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (b AutoTestWebhook) DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.OverseaPayChannel) {
	//TODO implement me
	panic("implement me")
}

func (b AutoTestWebhook) DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.OverseaPayChannel) (res *ro.ChannelRedirectInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}
