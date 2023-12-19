package paypal

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

// 接口文档：https://developer.paypal.com/docs/api/payments/v1/#payment_create
// clientId ATaWQ8G9oJNFyle9YCt59
// Secret EHUy5GALkYr1Qp0n6MepJY8LnUwYCBIWElG4Iv_DO3mdYcbB2l6zwJxk99OrPhbdNRLk7GkHEqb5RHEA
type Paypal struct {
}

func (p Paypal) DoRemoteChannelSubscriptionCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.CreateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelSubscriptionCancel(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.CancelSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelSubscriptionUpdate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.UpdateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelSubscriptionLists(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ListSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelSubscriptionWebhook(r *ghttp.Request) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelPlanInActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.CreateProductInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelPlanCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.CreatePlanInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelWebhook(r *ghttp.Request) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelRedirect(r *ghttp.Request) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelCapture(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelCancel(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelRefundStatusCheck(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelRefund(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}
