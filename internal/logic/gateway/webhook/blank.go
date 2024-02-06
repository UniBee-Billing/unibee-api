package webhook

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee-api/internal/logic/gateway/ro"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

type BlankWebhook struct {
}

func (b BlankWebhook) GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error) {
	//TODO implement me
	panic("implement me")
}

func (b BlankWebhook) GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway) {
	//TODO implement me
	panic("implement me")
}

func (b BlankWebhook) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *ro.GatewayRedirectInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}
