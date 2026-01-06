// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package system

import (
	"context"

	"unibee/api/system/auth"
	"unibee/api/system/information"
	"unibee/api/system/invoice"
	"unibee/api/system/payment"
	"unibee/api/system/plan"
	"unibee/api/system/refund"
	"unibee/api/system/subscription"
	"unibee/api/system/user"
)

type ISystemAuth interface {
	TokenGenerator(ctx context.Context, req *auth.TokenGeneratorReq) (res *auth.TokenGeneratorRes, err error)
}

type ISystemInformation interface {
	Get(ctx context.Context, req *information.GetReq) (res *information.GetRes, err error)
	SendMockMQ(ctx context.Context, req *information.SendMockMQReq) (res *information.SendMockMQRes, err error)
}

type ISystemInvoice interface {
	BulkChannelSync(ctx context.Context, req *invoice.BulkChannelSyncReq) (res *invoice.BulkChannelSyncRes, err error)
	ChannelSync(ctx context.Context, req *invoice.ChannelSyncReq) (res *invoice.ChannelSyncRes, err error)
	InternalWebhookSync(ctx context.Context, req *invoice.InternalWebhookSyncReq) (res *invoice.InternalWebhookSyncRes, err error)
	QuickbooksSync(ctx context.Context, req *invoice.QuickbooksSyncReq) (res *invoice.QuickbooksSyncRes, err error)
	BatchSendInvoiceWebhookEvent(ctx context.Context, req *invoice.BatchSendInvoiceWebhookEventReq) (res *invoice.BatchSendInvoiceWebhookEventRes, err error)
}

type ISystemPayment interface {
	PaymentCallbackAgain(ctx context.Context, req *payment.PaymentCallbackAgainReq) (res *payment.PaymentCallbackAgainRes, err error)
	PaymentGatewayDetail(ctx context.Context, req *payment.PaymentGatewayDetailReq) (res *payment.PaymentGatewayDetailRes, err error)
	Detail(ctx context.Context, req *payment.DetailReq) (res *payment.DetailRes, err error)
	PaymentGatewayCheck(ctx context.Context, req *payment.PaymentGatewayCheckReq) (res *payment.PaymentGatewayCheckRes, err error)
	GetPaymentExchangeRate(ctx context.Context, req *payment.GetPaymentExchangeRateReq) (res *payment.GetPaymentExchangeRateRes, err error)
}

type ISystemPlan interface {
	Detail(ctx context.Context, req *plan.DetailReq) (res *plan.DetailRes, err error)
}

type ISystemRefund interface {
	BulkChannelSync(ctx context.Context, req *refund.BulkChannelSyncReq) (res *refund.BulkChannelSyncRes, err error)
	GatewayDetail(ctx context.Context, req *refund.GatewayDetailReq) (res *refund.GatewayDetailRes, err error)
}

type ISystemSubscription interface {
	TestClockWalk(ctx context.Context, req *subscription.TestClockWalkReq) (res *subscription.TestClockWalkRes, err error)
	InternalWebhookSync(ctx context.Context, req *subscription.InternalWebhookSyncReq) (res *subscription.InternalWebhookSyncRes, err error)
	BatchSendSubActivateWebhookEvent(ctx context.Context, req *subscription.BatchSendSubActivateWebhookEventReq) (res *subscription.BatchSendSubActivateWebhookEventRes, err error)
	BatchSendSubUpdateWebhookEvent(ctx context.Context, req *subscription.BatchSendSubUpdateWebhookEventReq) (res *subscription.BatchSendSubUpdateWebhookEventRes, err error)
}

type ISystemUser interface {
	InternalWebhookSync(ctx context.Context, req *user.InternalWebhookSyncReq) (res *user.InternalWebhookSyncRes, err error)
}
