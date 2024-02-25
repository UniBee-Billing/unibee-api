// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package system

import (
	"context"
	
	"unibee-api/api/system/information"
	"unibee-api/api/system/invoice"
	"unibee-api/api/system/payment"
	"unibee-api/api/system/refund"
	"unibee-api/api/system/subscription"
)

type ISystemInformation interface {
	MerchantInformation(ctx context.Context, req *information.MerchantInformationReq) (res *information.MerchantInformationRes, err error)
}

type ISystemInvoice interface {
	BulkChannelSync(ctx context.Context, req *invoice.BulkChannelSyncReq) (res *invoice.BulkChannelSyncRes, err error)
	ChannelSync(ctx context.Context, req *invoice.ChannelSyncReq) (res *invoice.ChannelSyncRes, err error)
}

type ISystemPayment interface {
	BulkChannelSync(ctx context.Context, req *payment.BulkChannelSyncReq) (res *payment.BulkChannelSyncRes, err error)
	PaymentCallbackAgain(ctx context.Context, req *payment.PaymentCallbackAgainReq) (res *payment.PaymentCallbackAgainRes, err error)
	GatewayPaymentMethodList(ctx context.Context, req *payment.GatewayPaymentMethodListReq) (res *payment.GatewayPaymentMethodListRes, err error)
}

type ISystemRefund interface {
	BulkChannelSync(ctx context.Context, req *refund.BulkChannelSyncReq) (res *refund.BulkChannelSyncRes, err error)
}

type ISystemSubscription interface {
	BulkChannelSync(ctx context.Context, req *subscription.BulkChannelSyncReq) (res *subscription.BulkChannelSyncRes, err error)
	SubscriptionEndTrial(ctx context.Context, req *subscription.SubscriptionEndTrialReq) (res *subscription.SubscriptionEndTrialRes, err error)
	SubscriptionExpire(ctx context.Context, req *subscription.SubscriptionExpireReq) (res *subscription.SubscriptionExpireRes, err error)
	SubscriptionWalkTestClock(ctx context.Context, req *subscription.SubscriptionWalkTestClockReq) (res *subscription.SubscriptionWalkTestClockRes, err error)
}


