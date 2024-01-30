package _interface

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/channel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type RemotePayChannelInterface interface {
	// User
	DoRemoteChannelUserCreate(ctx context.Context, payChannel *entity.MerchantChannelConfig, user *entity.UserAccount) (res *ro.ChannelUserCreateInternalResp, err error)
	// Balance
	DoRemoteChannelUserDetailQuery(ctx context.Context, payChannel *entity.MerchantChannelConfig, userId int64) (res *ro.ChannelUserDetailQueryInternalResp, err error)
	DoRemoteChannelUserPaymentMethodListQuery(ctx context.Context, payChannel *entity.MerchantChannelConfig, userId int64) (res *ro.ChannelUserPaymentMethodListInternalResp, err error)
	DoRemoteChannelMerchantBalancesQuery(ctx context.Context, payChannel *entity.MerchantChannelConfig) (res *ro.ChannelMerchantBalanceQueryInternalResp, err error)
	// Invoice
	DoRemoteChannelInvoiceCreateAndPay(ctx context.Context, payChannel *entity.MerchantChannelConfig, createInvoiceInternalReq *ro.ChannelCreateInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error)
	DoRemoteChannelInvoicePay(ctx context.Context, payChannel *entity.MerchantChannelConfig, payInvoiceInternalReq *ro.ChannelPayInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error)
	DoRemoteChannelInvoiceCancel(ctx context.Context, payChannel *entity.MerchantChannelConfig, cancelInvoiceInternalReq *ro.ChannelCancelInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error)
	DoRemoteChannelInvoiceDetails(ctx context.Context, payChannel *entity.MerchantChannelConfig, channelInvoiceId string) (res *ro.ChannelDetailInvoiceInternalResp, err error)
	// Subscription Product And Plan - Deprecated
	DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (res *ro.ChannelCreateProductInternalResp, err error)
	DoRemoteChannelPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (res *ro.ChannelCreatePlanInternalResp, err error)
	DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (err error)
	DoRemoteChannelPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (err error)
	// Subscription - Deprecated
	DoRemoteChannelSubscriptionCreate(ctx context.Context, subscriptionCreateInternalReq *ro.ChannelCreateSubscriptionInternalReq) (res *ro.ChannelCreateSubscriptionInternalResp, err error)
	DoRemoteChannelSubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.ChannelCancelSubscriptionInternalReq) (res *ro.ChannelCancelSubscriptionInternalResp, err error)
	DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelCancelAtPeriodEndSubscriptionInternalResp, err error)
	DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error)
	DoRemoteChannelSubscriptionUpdateProrationPreview(ctx context.Context, subscriptionUpdateInternalReq *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionPreviewInternalResp, err error)
	DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionUpdateInternalReq *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionInternalResp, err error)
	DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error)
	DoRemoteChannelSubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription, newTrialEnd int64) (res *ro.ChannelDetailSubscriptionInternalResp, err error)
	DoRemoteChannelSubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error)
	// Payment
	DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error)
	DoRemoteChannelCapture(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCaptureRo, err error)
	DoRemoteChannelCancel(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCancelRo, err error)
	DoRemoteChannelPayStatusCheck(ctx context.Context, payment *entity.Payment) (res *ro.ChannelPaymentRo, err error)
	DoRemoteChannelPaymentList(ctx context.Context, payChannel *entity.MerchantChannelConfig, listReq *ro.ChannelPaymentListReq) (res []*ro.ChannelPaymentRo, err error)
	DoRemoteChannelPaymentDetail(ctx context.Context, payChannel *entity.MerchantChannelConfig, channelPaymentId string) (res *ro.ChannelPaymentRo, err error)
	DoRemoteChannelRefundStatusCheck(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error)
	DoRemoteChannelRefundList(ctx context.Context, payChannel *entity.MerchantChannelConfig, channelPaymentId string) (res []*ro.OutPayRefundRo, err error)
	DoRemoteChannelRefundDetail(ctx context.Context, payChannel *entity.MerchantChannelConfig, channelRefundId string) (res *ro.OutPayRefundRo, err error)
	DoRemoteChannelRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error)
}

type RemotePaymentChannelWebhookInterface interface {
	DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.MerchantChannelConfig) (err error)
	DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.MerchantChannelConfig)
	DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.MerchantChannelConfig) (res *ro.ChannelRedirectInternalResp, err error)
}
