package _interface

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/channel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type RemotePayChannelInterface interface {
	//ChannelUser
	DoRemoteChannelUserCreate(ctx context.Context, payChannel *entity.OverseaPayChannel, user *entity.UserAccount) (res *ro.ChannelUserCreateInternalResp, err error)
	//Balance 余额接口
	DoRemoteChannelUserBalancesQuery(ctx context.Context, payChannel *entity.OverseaPayChannel, userId int64) (res *ro.ChannelUserBalanceQueryInternalResp, err error)
	DoRemoteChannelMerchantBalancesQuery(ctx context.Context, payChannel *entity.OverseaPayChannel) (res *ro.ChannelMerchantBalanceQueryInternalResp, err error)
	//Invoice 发票接口
	DoRemoteChannelInvoiceCreateAndPay(ctx context.Context, payChannel *entity.OverseaPayChannel, createInvoiceInternalReq *ro.ChannelCreateInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error)
	DoRemoteChannelInvoicePay(ctx context.Context, payChannel *entity.OverseaPayChannel, payInvoiceInternalReq *ro.ChannelPayInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error)
	DoRemoteChannelInvoiceCancel(ctx context.Context, payChannel *entity.OverseaPayChannel, cancelInvoiceInternalReq *ro.ChannelCancelInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error)
	DoRemoteChannelInvoiceDetails(ctx context.Context, payChannel *entity.OverseaPayChannel, channelInvoiceId string) (res *ro.ChannelDetailInvoiceInternalResp, err error)
	//Subscription Product And Plan 订阅接口
	DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreateProductInternalResp, err error)
	DoRemoteChannelPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreatePlanInternalResp, err error)
	DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error)
	DoRemoteChannelPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error)
	//Subscription 订阅接口
	DoRemoteChannelSubscriptionCreate(ctx context.Context, subscriptionCreateInternalReq *ro.ChannelCreateSubscriptionInternalReq) (res *ro.ChannelCreateSubscriptionInternalResp, err error)
	DoRemoteChannelSubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.ChannelCancelSubscriptionInternalReq) (res *ro.ChannelCancelSubscriptionInternalResp, err error)
	DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelCancelAtPeriodEndSubscriptionInternalResp, err error)
	DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error)
	DoRemoteChannelSubscriptionUpdateProrationPreview(ctx context.Context, subscriptionUpdateInternalReq *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionPreviewInternalResp, err error)
	DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionUpdateInternalReq *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionInternalResp, err error)
	DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error)
	DoRemoteChannelSubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription, newTrialEnd int64) (res *ro.ChannelDetailSubscriptionInternalResp, err error)
	DoRemoteChannelSubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error)
	//Payment 是支付接口
	DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.OverseaPayChannel) (err error)
	DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.OverseaPayChannel)
	DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.OverseaPayChannel) (res *ro.ChannelRedirectInternalResp, err error)
	DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error)
	DoRemoteChannelCapture(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCaptureRo, err error)
	DoRemoteChannelCancel(ctx context.Context, pay *entity.Payment) (res *ro.OutPayCancelRo, err error)
	DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.Payment) (res *ro.ChannelPaymentRo, err error)
	DoRemoteChannelPaymentList(ctx context.Context, payChannel *entity.OverseaPayChannel, listReq *ro.ChannelPaymentListReq) (res []*ro.ChannelPaymentRo, err error)
	DoRemoteChannelPaymentDetail(ctx context.Context, payChannel *entity.OverseaPayChannel, channelPaymentId string) (res *ro.ChannelPaymentRo, err error)
	DoRemoteChannelRefundStatusCheck(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error)
	DoRemoteChannelRefundList(ctx context.Context, payChannel *entity.OverseaPayChannel, channelPaymentId string) (res []*ro.OutPayRefundRo, err error)
	DoRemoteChannelRefundDetail(ctx context.Context, payChannel *entity.OverseaPayChannel, channelRefundId string) (res *ro.OutPayRefundRo, err error)
	DoRemoteChannelRefund(ctx context.Context, pay *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error)
}
