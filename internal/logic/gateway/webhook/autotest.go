package webhook

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/oversea_pay"
)

type AutoTestWebhook struct {
}

func (b AutoTestWebhook) GatewayCheckAndSetupWebhook(ctx context.Context, gateway *entity.MerchantGateway) (err error) {
	return nil
}

func (b AutoTestWebhook) GatewayWebhook(r *ghttp.Request, gateway *entity.MerchantGateway) {
	//TODO implement me
	panic("implement me")
}

func (b AutoTestWebhook) GatewayRedirect(r *ghttp.Request, gateway *entity.MerchantGateway) (res *gateway_bean.GatewayRedirectResp, err error) {
	//TODO implement me
	panic("implement me")
}
