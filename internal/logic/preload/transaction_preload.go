package preload

import (
	"context"
	_interface "unibee/internal/interface/context"
	entity "unibee/internal/model/entity/default"
)

// TransactionPreloadData holds all preloaded data for transactions
type TransactionPreloadData struct {
	Gateways map[uint64]*entity.MerchantGateway
	Users    map[uint64]*entity.UserAccount
	Payments map[string]*entity.Payment
	Refunds  map[string]*entity.Refund
}

func TransactionListPreloadForContext(ctx context.Context, paymentTimelines []*entity.PaymentTimeline) {
	if len(paymentTimelines) > 3 {
		preloadData := TransactionListPreload(ctx, paymentTimelines)
		if preloadData != nil && _interface.GetBulkPreloadData(ctx) != nil {
			_interface.Context().Get(ctx).PreloadData.Gateways = preloadData.Gateways
			_interface.Context().Get(ctx).PreloadData.Users = preloadData.Users
			_interface.Context().Get(ctx).PreloadData.Payments = preloadData.Payments
			_interface.Context().Get(ctx).PreloadData.Refunds = preloadData.Refunds
		}
	}
}

// TransactionListPreload preloads all related data to avoid N+1 queries
func TransactionListPreload(ctx context.Context, paymentTimelines []*entity.PaymentTimeline) *TransactionPreloadData {
	if len(paymentTimelines) == 0 {
		return &TransactionPreloadData{
			Gateways: make(map[uint64]*entity.MerchantGateway),
			Users:    make(map[uint64]*entity.UserAccount),
		}
	}

	// Use maps to deduplicate
	gatewayIdsMap := make(map[uint64]bool)
	userIdsMap := make(map[uint64]bool)
	paymentIdsMap := make(map[string]bool)
	refundIdsMap := make(map[string]bool)

	// Collect all required IDs
	var gatewayIds []uint64
	var userIds []uint64
	var paymentIds []string
	var refundIds []string

	for _, timeline := range paymentTimelines {
		if timeline.GatewayId > 0 && !gatewayIdsMap[timeline.GatewayId] {
			gatewayIds = append(gatewayIds, timeline.GatewayId)
			gatewayIdsMap[timeline.GatewayId] = true
		}
		if timeline.UserId > 0 && !userIdsMap[timeline.UserId] {
			userIds = append(userIds, timeline.UserId)
			userIdsMap[timeline.UserId] = true
		}
		if len(timeline.PaymentId) > 0 && !paymentIdsMap[timeline.PaymentId] {
			paymentIds = append(paymentIds, timeline.PaymentId)
			paymentIdsMap[timeline.PaymentId] = true
		}
		if len(timeline.RefundId) > 0 && !refundIdsMap[timeline.RefundId] {
			refundIds = append(refundIds, timeline.RefundId)
			refundIdsMap[timeline.RefundId] = true
		}
	}

	// Bulk query related data
	gateways := bulkGetGatewaysByIds(ctx, gatewayIds)
	users := bulkGetUserAccountsByIds(ctx, userIds)
	payments := bulkGetPaymentsByPaymentIds(ctx, paymentIds)
	refunds := bulkGetRefundsByRefundIds(ctx, refundIds)

	return &TransactionPreloadData{
		Gateways: gateways,
		Users:    users,
		Payments: payments,
		Refunds:  refunds,
	}
}
