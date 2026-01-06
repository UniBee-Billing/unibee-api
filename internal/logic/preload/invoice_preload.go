package preload

import (
	"context"
	"fmt"
	"unibee/api/bean"
	_interface "unibee/internal/interface/context"
	entity "unibee/internal/model/entity/default"
)

type InvoicePreloadData struct {
	Payments                map[string]*entity.Payment
	Refunds                 map[string]*entity.Refund
	Gateways                map[uint64]*entity.MerchantGateway
	Users                   map[uint64]*entity.UserAccount
	Subscriptions           map[string]*entity.Subscription
	PromoCreditTransactions map[string]*entity.CreditTransaction
	Discounts               map[string]*entity.MerchantDiscountCode
	VatCountryRates         map[string]*bean.VatCountryRate // key: "merchantId_countryCode"
}

func InvoiceListPreloadForContext(ctx context.Context, invoices []*entity.Invoice) {
	if len(invoices) > 3 {
		preloadData := InvoiceListPreload(ctx, invoices)
		if preloadData != nil && _interface.GetBulkPreloadData(ctx) != nil {
			_interface.Context().Get(ctx).PreloadData.Payments = preloadData.Payments
			_interface.Context().Get(ctx).PreloadData.Refunds = preloadData.Refunds
			_interface.Context().Get(ctx).PreloadData.Gateways = preloadData.Gateways
			_interface.Context().Get(ctx).PreloadData.Users = preloadData.Users
			_interface.Context().Get(ctx).PreloadData.Subscriptions = preloadData.Subscriptions
			_interface.Context().Get(ctx).PreloadData.Discounts = preloadData.Discounts
			_interface.Context().Get(ctx).PreloadData.PromoCreditTransactions = preloadData.PromoCreditTransactions
		}
	}
}

func InvoiceListPreload(ctx context.Context, invoices []*entity.Invoice) *InvoicePreloadData {
	if len(invoices) == 0 {
		return &InvoicePreloadData{
			Payments:                make(map[string]*entity.Payment),
			Refunds:                 make(map[string]*entity.Refund),
			Gateways:                make(map[uint64]*entity.MerchantGateway),
			Users:                   make(map[uint64]*entity.UserAccount),
			Subscriptions:           make(map[string]*entity.Subscription),
			PromoCreditTransactions: make(map[string]*entity.CreditTransaction),
			VatCountryRates:         make(map[string]*bean.VatCountryRate),
		}
	}
	// 1. Collect all required IDs
	var invoiceIds, paymentIds, refundIds, subscriptionIds []string
	var discountCodes []string
	var merchantIds, gatewayIds, userIds []uint64
	var originalInvoiceIds []string
	var vatCountryPairs []string // Add VAT country code pairs

	// Use maps to deduplicate
	paymentIdsMap := make(map[string]bool)
	refundIdsMap := make(map[string]bool)
	gatewayIdsMap := make(map[uint64]bool)
	userIdsMap := make(map[uint64]bool)
	subscriptionIdsMap := make(map[string]bool)
	merchantIdsMap := make(map[uint64]bool)
	discountCodesMap := make(map[string]bool)
	originalInvoiceIdsMap := make(map[string]bool)
	vatCountryPairsMap := make(map[string]bool)

	for _, invoice := range invoices {
		invoiceIds = append(invoiceIds, invoice.InvoiceId)

		if len(invoice.PaymentId) > 0 && !paymentIdsMap[invoice.PaymentId] {
			paymentIds = append(paymentIds, invoice.PaymentId)
			paymentIdsMap[invoice.PaymentId] = true
		}
		if len(invoice.RefundId) > 0 && !refundIdsMap[invoice.RefundId] {
			refundIds = append(refundIds, invoice.RefundId)
			refundIdsMap[invoice.RefundId] = true
		}
		if invoice.GatewayId > 0 && !gatewayIdsMap[invoice.GatewayId] {
			gatewayIds = append(gatewayIds, invoice.GatewayId)
			gatewayIdsMap[invoice.GatewayId] = true
		}
		if invoice.UserId > 0 && !userIdsMap[invoice.UserId] {
			userIds = append(userIds, invoice.UserId)
			userIdsMap[invoice.UserId] = true
		}
		if len(invoice.SubscriptionId) > 0 && !subscriptionIdsMap[invoice.SubscriptionId] {
			subscriptionIds = append(subscriptionIds, invoice.SubscriptionId)
			subscriptionIdsMap[invoice.SubscriptionId] = true
		}
		if invoice.MerchantId > 0 && !merchantIdsMap[invoice.MerchantId] {
			merchantIds = append(merchantIds, invoice.MerchantId)
			merchantIdsMap[invoice.MerchantId] = true
		}
		if len(invoice.DiscountCode) > 0 && !discountCodesMap[invoice.DiscountCode] {
			discountCodes = append(discountCodes, invoice.DiscountCode)
			discountCodesMap[invoice.DiscountCode] = true
		}
		// Collect original invoice IDs for refunds
		if len(invoice.RefundId) > 0 && len(invoice.PaymentId) > 0 && !originalInvoiceIdsMap[invoice.PaymentId] {
			// We'll need to get the original invoice ID from payment later
			originalInvoiceIds = append(originalInvoiceIds, invoice.PaymentId)
			originalInvoiceIdsMap[invoice.PaymentId] = true
		}

		// Collect VAT country code pairs for bulk query
		if invoice.MerchantId > 0 && len(invoice.CountryCode) > 0 {
			key := fmt.Sprintf("%d_%s", invoice.MerchantId, invoice.CountryCode)
			if !vatCountryPairsMap[key] {
				vatCountryPairs = append(vatCountryPairs, key)
				vatCountryPairsMap[key] = true
			}
		}
	}

	// 2. Use real bulk query functions
	payments := bulkGetPaymentsByPaymentIds(ctx, paymentIds)
	refunds := bulkGetRefundsByRefundIds(ctx, refundIds)
	gateways := bulkGetGatewaysByIds(ctx, gatewayIds)
	users := bulkGetUserAccountsByIds(ctx, userIds)
	subscriptions := bulkGetSubscriptionsBySubscriptionIds(ctx, subscriptionIds)
	promoCreditTransactions := bulkGetPromoCreditTransactions(ctx, invoiceIds)
	discounts := bulkGetDiscountsByCodes(ctx, discountCodes)

	// 3. Bulk query VAT related data
	vatCountryRates := bulkGetVatCountryRatesByPairs(ctx, vatCountryPairs)

	// 5. Create preload data structure
	return &InvoicePreloadData{
		Payments:                payments,
		Refunds:                 refunds,
		Gateways:                gateways,
		Users:                   users,
		Subscriptions:           subscriptions,
		PromoCreditTransactions: promoCreditTransactions,
		VatCountryRates:         vatCountryRates,
		Discounts:               discounts,
	}
}
