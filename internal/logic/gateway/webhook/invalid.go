package webhook

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
)

type InvalidWebhook struct {
}

func (b InvalidWebhook) GatewayNewPaymentMethodRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (err error) {
	//TODO implement me
	panic("implement me")
}

func (b InvalidWebhook) GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error) {
	//TODO implement me
	panic("implement me")
}

func (b InvalidWebhook) GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway) {
	//TODO implement me
	panic("implement me")
}

func (b InvalidWebhook) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayRedirectResp, err error) {
	//TODO implement me
	panic("implement me")
}
