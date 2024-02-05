package api

import (
	"context"
	"go-oversea-pay/internal/logic/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type Invalid struct{}

func (i Invalid) DoRemoteChannelUserPaymentMethodListQuery(ctx context.Context, payChannel *entity.MerchantGateway, userId int64) (res *ro.ChannelUserPaymentMethodListInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelUserCreate(ctx context.Context, payChannel *entity.MerchantGateway, user *entity.UserAccount) (res *ro.ChannelUserCreateInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelSubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPaymentList(ctx context.Context, payChannel *entity.MerchantGateway, listReq *ro.ChannelPaymentListReq) (res []*ro.ChannelPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelRefundList(ctx context.Context, payChannel *entity.MerchantGateway, channelPaymentId string) (res []*ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPaymentDetail(ctx context.Context, payChannel *entity.MerchantGateway, channelPaymentId string) (res *ro.ChannelPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelRefundDetail(ctx context.Context, payChannel *entity.MerchantGateway, channelRefundId string) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelMerchantBalancesQuery(ctx context.Context, payChannel *entity.MerchantGateway) (res *ro.ChannelMerchantBalanceQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelInvoiceCancel(ctx context.Context, payChannel *entity.MerchantGateway, cancelInvoiceInternalReq *ro.ChannelCancelInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelUserDetailQuery(ctx context.Context, payChannel *entity.MerchantGateway, userId int64) (res *ro.ChannelUserDetailQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelInvoiceCreateAndPay(ctx context.Context, payChannel *entity.MerchantGateway, createInvoiceInternalReq *ro.ChannelCreateInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelInvoicePay(ctx context.Context, payChannel *entity.MerchantGateway, payInvoiceInternalReq *ro.ChannelPayInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (res *ro.ChannelCreateProductInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (res *ro.ChannelCreatePlanInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (err error) {
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

func (i Invalid) DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.ChannelCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
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

func (i Invalid) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelSubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription, newTrialEnd int64) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelCapture(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelCancel(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelPayStatusCheck(ctx context.Context, payment *entity.Payment) (res *ro.ChannelPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelRefundStatusCheck(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) DoRemoteChannelInvoiceDetails(ctx context.Context, payChannel *entity.MerchantGateway, channelInvoiceId string) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}
