// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package system

import (
	"context"

	"unibee/api/system/information"
	"unibee/api/system/invoice"
	"unibee/api/system/payment"
	"unibee/api/system/plan"
	"unibee/api/system/refund"
	"unibee/api/system/subscription"
)

type ISystemInformation interface {
	Get(ctx context.Context, req *information.GetReq) (res *information.GetRes, err error)
}

type ISystemInvoice interface {
	BulkChannelSync(ctx context.Context, req *invoice.BulkChannelSyncReq) (res *invoice.BulkChannelSyncRes, err error)
	ChannelSync(ctx context.Context, req *invoice.ChannelSyncReq) (res *invoice.ChannelSyncRes, err error)
	InternalWebhookSync(ctx context.Context, req *invoice.InternalWebhookSyncReq) (res *invoice.InternalWebhookSyncRes, err error)
}

type ISystemPayment interface {
	PaymentCallbackAgain(ctx context.Context, req *payment.PaymentCallbackAgainReq) (res *payment.PaymentCallbackAgainRes, err error)
	PaymentGatewayDetail(ctx context.Context, req *payment.PaymentGatewayDetailReq) (res *payment.PaymentGatewayDetailRes, err error)
}

type ISystemPlan interface {
	Detail(ctx context.Context, req *plan.DetailReq) (res *plan.DetailRes, err error)
}

type ISystemRefund interface {
	BulkChannelSync(ctx context.Context, req *refund.BulkChannelSyncReq) (res *refund.BulkChannelSyncRes, err error)
}

type ISystemSubscription interface {
	TestClockWalk(ctx context.Context, req *subscription.TestClockWalkReq) (res *subscription.TestClockWalkRes, err error)
}
