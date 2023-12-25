package out

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type Blank struct {
}

func (b Blank) DoRemoteChannelSubscriptionCreate(ctx context.Context, subscriptionRo *ro.CreateSubscriptionRo) (res *ro.CreateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelSubscriptionCancel(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.CancelSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionRo *ro.UpdateSubscriptionRo) (res *ro.UpdateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ListSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelSubscriptionWebhook(r *ghttp.Request) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.CreateProductInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.CreatePlanInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelWebhook(r *ghttp.Request) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelRedirect(r *ghttp.Request) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelCapture(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelCancel(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelRefundStatusCheck(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelRefund(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}
