package out

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type Invalid struct{}

func (i Invalid) DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPlanInActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.CreateProductInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPlanCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.CreatePlanInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelWebhook(r *ghttp.Request) {
	panic("invalid pay channel")
}

func (i Invalid) DoRemoteChannelRedirect(r *ghttp.Request) {
	panic("invalid pay channel")
}

func (i Invalid) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	panic("invalid pay channel")
}

func (i Invalid) DoRemoteChannelCapture(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCaptureRo, err error) {
	panic("invalid pay channel")
}

func (i Invalid) DoRemoteChannelCancel(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCancelRo, err error) {
	panic("invalid pay channel")
}

func (i Invalid) DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayRo, err error) {
	panic("invalid pay channel")
}

func (i Invalid) DoRemoteChannelRefundStatusCheck(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	panic("invalid pay channel")
}

func (i Invalid) DoRemoteChannelRefund(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	panic("invalid pay channel")
}
