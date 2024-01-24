package out

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/channel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type Invalid struct{}

func (i Invalid) DoRemoteChannelSubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPaymentList(ctx context.Context, payChannel *entity.OverseaPayChannel, listReq *ro.ChannelPaymentListReq) (res []*ro.OutPayRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelRefundList(ctx context.Context, payChannel *entity.OverseaPayChannel, channelPaymentId string) (res []*ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPaymentDetail(ctx context.Context, payChannel *entity.OverseaPayChannel, channelPaymentId string) (res *ro.OutPayRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelRefundDetail(ctx context.Context, payChannel *entity.OverseaPayChannel, channelRefundId string) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelMerchantBalancesQuery(ctx context.Context, payChannel *entity.OverseaPayChannel) (res *ro.ChannelMerchantBalanceQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelInvoiceCancel(ctx context.Context, payChannel *entity.OverseaPayChannel, cancelInvoiceInternalReq *ro.ChannelCancelInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelUserBalancesQuery(ctx context.Context, payChannel *entity.OverseaPayChannel, customerId string) (res *ro.ChannelUserBalanceQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelInvoiceCreateAndPay(ctx context.Context, payChannel *entity.OverseaPayChannel, createInvoiceInternalReq *ro.ChannelCreateInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelInvoicePay(ctx context.Context, payChannel *entity.OverseaPayChannel, payInvoiceInternalReq *ro.ChannelPayInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreateProductInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreatePlanInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelSubscriptionCreate(ctx context.Context, subscriptionRo *ro.ChannelCreateSubscriptionInternalReq) (res *ro.ChannelCreateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelSubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.ChannelCancelSubscriptionInternalReq) (res *ro.ChannelCancelSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelSubscriptionUpdateProrationPreview(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionPreviewInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelSubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription, newTrialEnd int64) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.OverseaPayChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.OverseaPayChannel) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.OverseaPayChannel) (res *ro.ChannelRedirectInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelCapture(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelCancel(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.Payment) (res *ro.OutPayRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelRefundStatusCheck(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelRefund(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelInvoiceDetails(ctx context.Context, payChannel *entity.OverseaPayChannel, channelInvoiceId string) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}
