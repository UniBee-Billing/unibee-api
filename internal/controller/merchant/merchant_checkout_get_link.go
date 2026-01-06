package merchant

import (
	"context"
	"fmt"
	"unibee/internal/cmd/config"

	"unibee/api/merchant/checkout"
)

func (c *ControllerCheckout) GetLink(ctx context.Context, req *checkout.GetLinkReq) (res *checkout.GetLinkRes, err error) {
	link := fmt.Sprintf("%s/checkout?checkoutId=%d&planId=%d&env=%s", config.GetConfigInstance().Server.GetHostedPath(), req.CheckoutId, req.PlanId, config.GetConfigInstance().Env)
	return &checkout.GetLinkRes{Link: link}, nil
}
