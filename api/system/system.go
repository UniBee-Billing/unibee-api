// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package system

import (
	"context"
	
	"unibee/api/system/information"
	"unibee/api/system/invoice"
	"unibee/api/system/payment"
	"unibee/api/system/refund"
	"unibee/api/system/subscription"
)

type ISystemInformation interface {
	Get(ctx context.Context, req *information.GetReq) (res *information.GetRes, err error)
}

type ISystemInvoice interface {
	BulkChannelSync(ctx context.Context, req *invoice.BulkChannelSyncReq) (res *invoice.BulkChannelSyncRes, err error)
	ChannelSync(ctx context.Context, req *invoice.ChannelSyncReq) (res *invoice.ChannelSyncRes, err error)
}

type ISystemPayment interface {
	PaymentCallbackAgain(ctx context.Context, req *payment.PaymentCallbackAgainReq) (res *payment.PaymentCallbackAgainRes, err error)
	GatewayPaymentMethodList(ctx context.Context, req *payment.GatewayPaymentMethodListReq) (res *payment.GatewayPaymentMethodListRes, err error)
}

type ISystemRefund interface {
	BulkChannelSync(ctx context.Context, req *refund.BulkChannelSyncReq) (res *refund.BulkChannelSyncRes, err error)
}

type ISystemSubscription interface {
	TestClockWalk(ctx context.Context, req *subscription.TestClockWalkReq) (res *subscription.TestClockWalkRes, err error)
}


