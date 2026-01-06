package preload

import (
	"context"
	"strconv"
	"strings"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/default"
)

func bulkGetPaymentsByPaymentIds(ctx context.Context, paymentIds []string) map[string]*entity.Payment {
	if len(paymentIds) == 0 {
		return make(map[string]*entity.Payment)
	}

	var payments []*entity.Payment
	err := dao.Payment.Ctx(ctx).WhereIn(dao.Payment.Columns().PaymentId, paymentIds).Scan(&payments)
	if err != nil {
		return make(map[string]*entity.Payment)
	}

	result := make(map[string]*entity.Payment)
	for _, payment := range payments {
		result[payment.PaymentId] = payment
	}
	return result
}

// bulkGetRefundsByRefundIds bulk query refunds
func bulkGetRefundsByRefundIds(ctx context.Context, refundIds []string) map[string]*entity.Refund {
	if len(refundIds) == 0 {
		return make(map[string]*entity.Refund)
	}

	var refunds []*entity.Refund
	err := dao.Refund.Ctx(ctx).WhereIn(dao.Refund.Columns().RefundId, refundIds).Scan(&refunds)
	if err != nil {
		return make(map[string]*entity.Refund)
	}

	result := make(map[string]*entity.Refund)
	for _, refund := range refunds {
		result[refund.RefundId] = refund
	}
	return result
}

// bulkGetGatewaysByIds bulk query gateways
func bulkGetGatewaysByIds(ctx context.Context, gatewayIds []uint64) map[uint64]*entity.MerchantGateway {
	if len(gatewayIds) == 0 {
		return make(map[uint64]*entity.MerchantGateway)
	}

	var gateways []*entity.MerchantGateway
	err := dao.MerchantGateway.Ctx(ctx).WhereIn(dao.MerchantGateway.Columns().Id, gatewayIds).Scan(&gateways)
	if err != nil {
		return make(map[uint64]*entity.MerchantGateway)
	}

	result := make(map[uint64]*entity.MerchantGateway)
	for _, gateway := range gateways {
		result[gateway.Id] = gateway
	}
	return result
}

// bulkGetUserAccountsByIds bulk query users
func bulkGetUserAccountsByIds(ctx context.Context, userIds []uint64) map[uint64]*entity.UserAccount {
	if len(userIds) == 0 {
		return make(map[uint64]*entity.UserAccount)
	}

	var users []*entity.UserAccount
	err := dao.UserAccount.Ctx(ctx).WhereIn(dao.UserAccount.Columns().Id, userIds).Scan(&users)
	if err != nil {
		return make(map[uint64]*entity.UserAccount)
	}

	result := make(map[uint64]*entity.UserAccount)
	for _, user := range users {
		result[user.Id] = user
	}
	return result
}

// bulkGetSubscriptionsBySubscriptionIds bulk query subscriptions
func bulkGetSubscriptionsBySubscriptionIds(ctx context.Context, subscriptionIds []string) map[string]*entity.Subscription {
	if len(subscriptionIds) == 0 {
		return make(map[string]*entity.Subscription)
	}

	var subscriptions []*entity.Subscription
	err := dao.Subscription.Ctx(ctx).WhereIn(dao.Subscription.Columns().SubscriptionId, subscriptionIds).Scan(&subscriptions)
	if err != nil {
		return make(map[string]*entity.Subscription)
	}

	result := make(map[string]*entity.Subscription)
	for _, subscription := range subscriptions {
		result[subscription.SubscriptionId] = subscription
	}
	return result
}

func bulkGetPromoCreditTransactions(ctx context.Context, invoiceIds []string) map[string]*entity.CreditTransaction {
	if len(invoiceIds) == 0 {
		return make(map[string]*entity.CreditTransaction)
	}
	var creditTransactions []*entity.CreditTransaction
	err := dao.CreditTransaction.Ctx(ctx).WhereIn(dao.CreditTransaction.Columns().InvoiceId, invoiceIds).Scan(&creditTransactions)
	if err != nil {
		return make(map[string]*entity.CreditTransaction)
	}
	result := make(map[string]*entity.CreditTransaction)
	for _, one := range creditTransactions {
		result[one.InvoiceId] = one
	}

	return result
}

// bulkGetVatCountryRatesByPairs bulk query VAT country rates by merchant-country pairs
func bulkGetVatCountryRatesByPairs(ctx context.Context, vatCountryPairs []string) map[string]*bean.VatCountryRate {
	if len(vatCountryPairs) == 0 {
		return make(map[string]*bean.VatCountryRate)
	}

	// Deduplicate
	uniquePairs := make(map[string]bool)
	for _, pair := range vatCountryPairs {
		uniquePairs[pair] = true
	}

	result := make(map[string]*bean.VatCountryRate)

	// Query each unique merchantId_countryCode pair
	for pair := range uniquePairs {
		parts := strings.Split(pair, "_")
		if len(parts) == 2 {
			merchantIdStr := parts[0]
			countryCode := parts[1]

			merchantId, err := strconv.ParseUint(merchantIdStr, 10, 64)
			if err != nil {
				continue
			}

			// Query VAT country rate
			vatCountryRate, _ := vat_gateway.QueryVatCountryRateByMerchant(ctx, merchantId, countryCode)
			if vatCountryRate != nil {
				result[pair] = vatCountryRate
			}
		}
	}

	return result
}

//
//// bulkGetGatewaysByIds bulk query gateways
//func bulkGetGatewaysByIds(ctx context.Context, gatewayIds []uint64) map[uint64]*entity.MerchantGateway {
//	if len(gatewayIds) == 0 {
//		return make(map[uint64]*entity.MerchantGateway)
//	}
//
//	list := query.GetGatewaysByIds(ctx, gatewayIds)
//	result := make(map[uint64]*entity.MerchantGateway)
//	if len(list) == 0 {
//		return result
//	}
//	for _, gateway := range list {
//		result[gateway.Id] = gateway
//	}
//	return result
//}
//
//// bulkGetUserAccountsByIds bulk query users
//func bulkGetUserAccountsByIds(ctx context.Context, userIds []uint64) map[uint64]*entity.UserAccount {
//	if len(userIds) == 0 {
//		return make(map[uint64]*entity.UserAccount)
//	}
//
//	list := query.GetUserAccountsByIds(ctx, userIds)
//	result := make(map[uint64]*entity.UserAccount)
//	if len(list) == 0 {
//		return result
//	}
//	for _, user := range list {
//		result[user.Id] = user
//	}
//	return result
//}

// bulkGetPlansByIds bulk query plans
func bulkGetPlansByIds(ctx context.Context, planIds []uint64) map[uint64]*entity.Plan {
	if len(planIds) == 0 {
		return make(map[uint64]*entity.Plan)
	}

	var plans []*entity.Plan
	err := dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, planIds).Where(dao.Plan.Columns().IsDeleted, 0).Scan(&plans)
	if err != nil {
		return make(map[uint64]*entity.Plan)
	}

	result := make(map[uint64]*entity.Plan)
	for _, plan := range plans {
		result[plan.Id] = plan
	}
	return result
}

// bulkGetProductsByIds bulk query products
func bulkGetProductsByIds(ctx context.Context, productIds []uint64) map[uint64]*entity.Product {
	if len(productIds) == 0 {
		return make(map[uint64]*entity.Product)
	}

	var products []*entity.Product
	err := dao.Product.Ctx(ctx).WhereIn(dao.Product.Columns().Id, productIds).Where(dao.Product.Columns().IsDeleted, 0).Scan(&products)
	if err != nil {
		return make(map[uint64]*entity.Product)
	}

	result := make(map[uint64]*entity.Product)
	for _, product := range products {
		result[product.Id] = product
	}
	return result
}

// bulkGetProductsByIds bulk query products
func bulkGetInvoicesByInvoiceIds(ctx context.Context, invoiceIds []string) map[string]*entity.Invoice {
	if len(invoiceIds) == 0 {
		return make(map[string]*entity.Invoice)
	}

	var invoices []*entity.Invoice
	err := dao.Invoice.Ctx(ctx).WhereIn(dao.Invoice.Columns().InvoiceId, invoiceIds).Where(dao.Invoice.Columns().IsDeleted, 0).Scan(&invoices)
	if err != nil {
		return make(map[string]*entity.Invoice)
	}

	result := make(map[string]*entity.Invoice)
	for _, invoice := range invoices {
		result[invoice.InvoiceId] = invoice
	}
	return result
}

func bulkGetDiscountsByCodes(ctx context.Context, codes []string) map[string]*entity.MerchantDiscountCode {
	if len(codes) == 0 {
		return make(map[string]*entity.MerchantDiscountCode)
	}

	var discounts []*entity.MerchantDiscountCode
	err := dao.MerchantDiscountCode.Ctx(ctx).WhereIn(dao.MerchantDiscountCode.Columns().Code, codes).Scan(&discounts)
	if err != nil {
		return make(map[string]*entity.MerchantDiscountCode)
	}

	result := make(map[string]*entity.MerchantDiscountCode)
	for _, discount := range discounts {
		result[discount.Code] = discount
	}
	return result
}
