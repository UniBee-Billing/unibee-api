package preload

import (
	"context"
	_interface "unibee/internal/interface/context"
	entity "unibee/internal/model/entity/default"
)

// SubscriptionPreloadData holds all preloaded data for subscriptions
type SubscriptionPreloadData struct {
	Gateways map[uint64]*entity.MerchantGateway
	Users    map[uint64]*entity.UserAccount
	Plans    map[uint64]*entity.Plan
	Products map[uint64]*entity.Product
	Invoices map[string]*entity.Invoice
}

func SubscriptionListPreloadForContext(ctx context.Context, subscriptions []*entity.Subscription) {
	if len(subscriptions) > 3 {
		preloadData := SubscriptionListPreload(ctx, subscriptions)
		if preloadData != nil && _interface.GetBulkPreloadData(ctx) != nil {
			_interface.Context().Get(ctx).PreloadData.Gateways = preloadData.Gateways
			_interface.Context().Get(ctx).PreloadData.Users = preloadData.Users
			_interface.Context().Get(ctx).PreloadData.Plans = preloadData.Plans
			_interface.Context().Get(ctx).PreloadData.Products = preloadData.Products
			_interface.Context().Get(ctx).PreloadData.Invoices = preloadData.Invoices
		}
	}
}

// SubscriptionListPreload preloads all related data to avoid N+1 queries
func SubscriptionListPreload(ctx context.Context, subscriptions []*entity.Subscription) *SubscriptionPreloadData {
	if len(subscriptions) == 0 {
		return &SubscriptionPreloadData{
			Gateways: make(map[uint64]*entity.MerchantGateway),
			Users:    make(map[uint64]*entity.UserAccount),
			Plans:    make(map[uint64]*entity.Plan),
			Products: make(map[uint64]*entity.Product),
			Invoices: make(map[string]*entity.Invoice),
		}
	}

	// Use maps to deduplicate
	gatewayIdsMap := make(map[uint64]bool)
	userIdsMap := make(map[uint64]bool)
	planIdsMap := make(map[uint64]bool)
	productIdsMap := make(map[uint64]bool)
	invoiceIdsMap := make(map[string]bool)

	// Collect all required IDs
	var gatewayIds []uint64
	var userIds []uint64
	var planIds []uint64
	var productIds []uint64
	var invoiceIds []string

	for _, subscription := range subscriptions {
		if subscription.GatewayId > 0 && !gatewayIdsMap[subscription.GatewayId] {
			gatewayIds = append(gatewayIds, subscription.GatewayId)
			gatewayIdsMap[subscription.GatewayId] = true
		}
		if subscription.UserId > 0 && !userIdsMap[subscription.UserId] {
			userIds = append(userIds, subscription.UserId)
			userIdsMap[subscription.UserId] = true
		}
		if subscription.PlanId > 0 && !planIdsMap[subscription.PlanId] {
			planIds = append(planIds, subscription.PlanId)
			planIdsMap[subscription.PlanId] = true
		}
		if subscription.LatestInvoiceId != "" && !invoiceIdsMap[subscription.LatestInvoiceId] {
			invoiceIds = append(invoiceIds, subscription.LatestInvoiceId)
			invoiceIdsMap[subscription.LatestInvoiceId] = true
		}
	}

	// Get plan details to collect product IDs
	if len(planIds) > 0 {
		plans := bulkGetPlansByIds(ctx, planIds)
		for _, plan := range plans {
			if plan.ProductId > 0 && !productIdsMap[uint64(plan.ProductId)] {
				productIds = append(productIds, uint64(plan.ProductId))
				productIdsMap[uint64(plan.ProductId)] = true
			}
		}
	}

	// Bulk query related data
	gateways := bulkGetGatewaysByIds(ctx, gatewayIds)
	users := bulkGetUserAccountsByIds(ctx, userIds)
	plans := bulkGetPlansByIds(ctx, planIds)
	products := bulkGetProductsByIds(ctx, productIds)
	invoices := bulkGetInvoicesByInvoiceIds(ctx, invoiceIds)

	return &SubscriptionPreloadData{
		Gateways: gateways,
		Users:    users,
		Plans:    plans,
		Products: products,
		Invoices: invoices,
	}
}
