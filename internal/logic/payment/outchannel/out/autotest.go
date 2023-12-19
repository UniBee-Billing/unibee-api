package out

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type AutoTest struct {
}

func (a AutoTest) DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelPlanInActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.CreateProductInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelPlanCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.CreatePlanInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelWebhook(r *ghttp.Request) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelRedirect(r *ghttp.Request) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelCapture(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelCancel(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelRefundStatusCheck(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelRefund(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}
