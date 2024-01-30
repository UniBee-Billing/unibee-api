package out

import (
	"context"
	"go-oversea-pay/internal/logic/channel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type AutoTest struct {
}

func (a AutoTest) DoRemoteChannelUserPaymentMethodListQuery(ctx context.Context, payChannel *entity.OverseaPayChannel, userId int64) (res *ro.ChannelUserPaymentMethodListInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelUserCreate(ctx context.Context, payChannel *entity.OverseaPayChannel, user *entity.UserAccount) (res *ro.ChannelUserCreateInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelSubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelPaymentList(ctx context.Context, payChannel *entity.OverseaPayChannel, listReq *ro.ChannelPaymentListReq) (res []*ro.ChannelPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelRefundList(ctx context.Context, payChannel *entity.OverseaPayChannel, channelPaymentId string) (res []*ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelPaymentDetail(ctx context.Context, payChannel *entity.OverseaPayChannel, channelPaymentId string) (res *ro.ChannelPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelRefundDetail(ctx context.Context, payChannel *entity.OverseaPayChannel, channelRefundId string) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelMerchantBalancesQuery(ctx context.Context, payChannel *entity.OverseaPayChannel) (res *ro.ChannelMerchantBalanceQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelInvoiceCancel(ctx context.Context, payChannel *entity.OverseaPayChannel, cancelInvoiceInternalReq *ro.ChannelCancelInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelUserDetailQuery(ctx context.Context, payChannel *entity.OverseaPayChannel, userId int64) (res *ro.ChannelUserDetailQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelInvoiceCreateAndPay(ctx context.Context, payChannel *entity.OverseaPayChannel, createInvoiceInternalReq *ro.ChannelCreateInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelInvoicePay(ctx context.Context, payChannel *entity.OverseaPayChannel, payInvoiceInternalReq *ro.ChannelPayInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelInvoiceDetails(ctx context.Context, payChannel *entity.OverseaPayChannel, channelInvoiceId string) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelSubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription, newTrialEnd int64) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelSubscriptionUpdateProrationPreview(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionPreviewInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelSubscriptionCreate(ctx context.Context, subscriptionRo *ro.ChannelCreateSubscriptionInternalReq) (res *ro.ChannelCreateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelSubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.ChannelCancelSubscriptionInternalReq) (res *ro.ChannelCancelSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreateProductInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreatePlanInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelCapture(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelCancel(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.Payment) (res *ro.ChannelPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelRefundStatusCheck(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) DoRemoteChannelRefund(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}
