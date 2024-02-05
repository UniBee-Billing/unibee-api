package api

import (
	"context"
	"go-oversea-pay/internal/logic/gateway/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type Invalid struct{}

func (i Invalid) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *ro.GatewayUserPaymentMethodListInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *ro.GatewayUserCreateInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewaySubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *ro.GatewayPaymentListReq) (res []*ro.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, channelPaymentId string) (res []*ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, channelPaymentId string) (res *ro.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, channelRefundId string) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *ro.GatewayMerchantBalanceQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayInvoiceCancel(ctx context.Context, gateway *entity.MerchantGateway, cancelInvoiceInternalReq *ro.GatewayCancelInvoiceInternalReq) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *ro.GatewayUserDetailQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayInvoiceCreateAndPay(ctx context.Context, gateway *entity.MerchantGateway, createInvoiceInternalReq *ro.GatewayCreateInvoiceInternalReq) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayInvoicePay(ctx context.Context, gateway *entity.MerchantGateway, payInvoiceInternalReq *ro.GatewayPayInvoiceInternalReq) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (res *ro.GatewayCreateProductInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (res *ro.GatewayCreatePlanInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewaySubscriptionCreate(ctx context.Context, subscriptionRo *ro.GatewayCreateSubscriptionInternalReq) (res *ro.GatewayCreateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewaySubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.GatewayCancelSubscriptionInternalReq) (res *ro.GatewayCancelSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewaySubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewaySubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewaySubscriptionUpdateProrationPreview(ctx context.Context, subscriptionRo *ro.GatewayUpdateSubscriptionInternalReq) (res *ro.GatewayUpdateSubscriptionPreviewInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewaySubscriptionUpdate(ctx context.Context, subscriptionRo *ro.GatewayUpdateSubscriptionInternalReq) (res *ro.GatewayUpdateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewaySubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewaySubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription, newTrialEnd int64) (res *ro.GatewayDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayCapture(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayCancel(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayPayStatusCheck(ctx context.Context, payment *entity.Payment) (res *ro.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayRefundStatusCheck(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (i Invalid) GatewayInvoiceDetails(ctx context.Context, gateway *entity.MerchantGateway, channelInvoiceId string) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}
