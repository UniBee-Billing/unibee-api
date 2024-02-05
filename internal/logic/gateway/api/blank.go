package api

import (
	"context"
	"go-oversea-pay/internal/logic/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type Blank struct {
}

func (b Blank) DoRemoteChannelUserPaymentMethodListQuery(ctx context.Context, payChannel *entity.MerchantGateway, userId int64) (res *ro.ChannelUserPaymentMethodListInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelUserCreate(ctx context.Context, payChannel *entity.MerchantGateway, user *entity.UserAccount) (res *ro.ChannelUserCreateInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelSubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelPaymentList(ctx context.Context, payChannel *entity.MerchantGateway, listReq *ro.ChannelPaymentListReq) (res []*ro.ChannelPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelRefundList(ctx context.Context, payChannel *entity.MerchantGateway, channelPaymentId string) (res []*ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelPaymentDetail(ctx context.Context, payChannel *entity.MerchantGateway, channelPaymentId string) (res *ro.ChannelPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelRefundDetail(ctx context.Context, payChannel *entity.MerchantGateway, channelRefundId string) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelMerchantBalancesQuery(ctx context.Context, payChannel *entity.MerchantGateway) (res *ro.ChannelMerchantBalanceQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelInvoiceCancel(ctx context.Context, payChannel *entity.MerchantGateway, cancelInvoiceInternalReq *ro.ChannelCancelInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelUserDetailQuery(ctx context.Context, payChannel *entity.MerchantGateway, userId int64) (res *ro.ChannelUserDetailQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelInvoiceCreateAndPay(ctx context.Context, payChannel *entity.MerchantGateway, createInvoiceInternalReq *ro.ChannelCreateInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelInvoicePay(ctx context.Context, payChannel *entity.MerchantGateway, payInvoiceInternalReq *ro.ChannelPayInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelInvoiceDetails(ctx context.Context, payChannel *entity.MerchantGateway, channelInvoiceId string) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelSubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription, newTrialEnd int64) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelSubscriptionUpdateProrationPreview(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionPreviewInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelSubscriptionCreate(ctx context.Context, subscriptionRo *ro.ChannelCreateSubscriptionInternalReq) (res *ro.ChannelCreateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelSubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.ChannelCancelSubscriptionInternalReq) (res *ro.ChannelCancelSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.ChannelCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (res *ro.ChannelCreateProductInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (res *ro.ChannelCreatePlanInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelCapture(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelCancel(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.Payment) (res *ro.ChannelPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelRefundStatusCheck(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (b Blank) DoRemoteChannelRefund(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}
