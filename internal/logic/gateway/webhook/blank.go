package webhook

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/oversea_pay"
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

func (b BlankWebhook) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayRedirectResp, err error) {
	//TODO implement me
	panic("implement me")
}
