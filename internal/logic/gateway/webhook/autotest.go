package webhook

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type AutoTestWebhook struct {
}

func (b AutoTestWebhook) DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.MerchantGateway) (err error) {
	//TODO implement me
	panic("implement me")
}

func (b AutoTestWebhook) DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.MerchantGateway) {
	//TODO implement me
	panic("implement me")
}

func (b AutoTestWebhook) DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.MerchantGateway) (res *ro.ChannelRedirectInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}