package api

import (
	"context"
	"unibee-api/internal/logic/gateway/ro"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

type AutoTest struct {
}

func (a AutoTest) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *ro.GatewayUserPaymentMethodListInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *ro.GatewayUserCreateInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewaySubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *ro.GatewayPaymentListReq) (res []*ro.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res *ro.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *ro.GatewayMerchantBalanceQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayInvoiceCancel(ctx context.Context, gateway *entity.MerchantGateway, cancelInvoiceInternalReq *ro.GatewayCancelInvoiceInternalReq) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *ro.GatewayUserDetailQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayInvoiceCreateAndPay(ctx context.Context, gateway *entity.MerchantGateway, createInvoiceInternalReq *ro.GatewayCreateInvoiceInternalReq) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayInvoicePay(ctx context.Context, gateway *entity.MerchantGateway, payInvoiceInternalReq *ro.GatewayPayInvoiceInternalReq) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayInvoiceDetails(ctx context.Context, gateway *entity.MerchantGateway, gatewayInvoiceId string) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewaySubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription, newTrialEnd int64) (res *ro.GatewayDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewaySubscriptionUpdateProrationPreview(ctx context.Context, subscriptionRo *ro.GatewayUpdateSubscriptionInternalReq) (res *ro.GatewayUpdateSubscriptionPreviewInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewaySubscriptionCreate(ctx context.Context, subscriptionRo *ro.GatewayCreateSubscriptionInternalReq) (res *ro.GatewayCreateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewaySubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.GatewayCancelSubscriptionInternalReq) (res *ro.GatewayCancelSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewaySubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewaySubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewaySubscriptionUpdate(ctx context.Context, subscriptionRo *ro.GatewayUpdateSubscriptionInternalReq) (res *ro.GatewayUpdateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewaySubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (res *ro.GatewayCreateProductInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (res *ro.GatewayCreatePlanInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayCapture(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayCancel(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayPayStatusCheck(ctx context.Context, pay *entity.Payment) (res *ro.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayRefundStatusCheck(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (a AutoTest) GatewayRefund(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}
