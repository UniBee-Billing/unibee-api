package webhook

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type AutoTestWebhook struct {
}

func (b AutoTestWebhook) GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error) {
	//TODO implement me
	panic("implement me")
}

func (b AutoTestWebhook) GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway) {
	//TODO implement me
	panic("implement me")
}

func (b AutoTestWebhook) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *ro.GatewayRedirectInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}
