// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package system

import (
	"context"
	
	"go-oversea-pay/api/system/invoice"
	"go-oversea-pay/api/system/payment"
	"go-oversea-pay/api/system/refund"
	"go-oversea-pay/api/system/subscription"
)

type ISystemInvoice interface {
	BulkChannelSync(ctx context.Context, req *invoice.BulkChannelSyncReq) (res *invoice.BulkChannelSyncRes, err error)
	ChannelSync(ctx context.Context, req *invoice.ChannelSyncReq) (res *invoice.ChannelSyncRes, err error)
}

type ISystemPayment interface {
	BulkChannelSync(ctx context.Context, req *payment.BulkChannelSyncReq) (res *payment.BulkChannelSyncRes, err error)
}

type ISystemRefund interface {
	BulkChannelSync(ctx context.Context, req *refund.BulkChannelSyncReq) (res *refund.BulkChannelSyncRes, err error)
}

type ISystemSubscription interface {
	BulkChannelSync(ctx context.Context, req *subscription.BulkChannelSyncReq) (res *subscription.BulkChannelSyncRes, err error)
	SubscriptionEndTrial(ctx context.Context, req *subscription.SubscriptionEndTrialReq) (res *subscription.SubscriptionEndTrialRes, err error)
}


