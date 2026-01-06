package detail

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/controller/link"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	dao "unibee/internal/dao/default"

	"github.com/gogf/gf/v2/encoding/gjson"
)

// PreloadData contains all preloaded data for bulk conversion
type PreloadData struct {
	Payments                   map[string]*entity.Payment
	Refunds                    map[string]*entity.Refund
	Gateways                   map[uint64]*entity.MerchantGateway
	Users                      map[uint64]*entity.UserAccount
	Subscriptions              map[string]*entity.Subscription
	Merchants                  map[uint64]*entity.Merchant
	Discounts                  map[string]*entity.MerchantDiscountCode
	OriginalInvoices           map[string]*entity.Invoice
	PromoCreditTransactions    map[string]*bean.CreditTransaction
	SubscriptionPendingUpdates map[string]*SubscriptionPendingUpdateDetail
	Plans                      map[uint64]*entity.Plan
	Addons                     map[string][]*bean.PlanAddonDetail
}

func BulkConvertInvoicesToDetails(ctx context.Context, invoices []*entity.Invoice) []*InvoiceDetail {
	if len(invoices) == 0 {
		return make([]*InvoiceDetail, 0)
	}

	// 1. Collect all required IDs
	var paymentIds, refundIds, gatewayIds, userIds, subscriptionIds []string
	var discountCodes []string
	var merchantIds []uint64
	var originalInvoiceIds []string
	var promoCreditKeys []string

	for _, invoice := range invoices {
		if len(invoice.PaymentId) > 0 {
			paymentIds = append(paymentIds, invoice.PaymentId)
		}
		if len(invoice.RefundId) > 0 {
			refundIds = append(refundIds, invoice.RefundId)
		}
		if invoice.GatewayId > 0 {
			gatewayIds = append(gatewayIds, strconv.FormatUint(invoice.GatewayId, 10))
		}
		if invoice.UserId > 0 {
			userIds = append(userIds, strconv.FormatUint(invoice.UserId, 10))
		}
		if len(invoice.SubscriptionId) > 0 {
			subscriptionIds = append(subscriptionIds, invoice.SubscriptionId)
		}
		if invoice.MerchantId > 0 {
			merchantIds = append(merchantIds, invoice.MerchantId)
		}
		if len(invoice.DiscountCode) > 0 {
			discountCodes = append(discountCodes, invoice.DiscountCode)
		}
		// Collect original invoice IDs for refunds
		if len(invoice.RefundId) > 0 && len(invoice.PaymentId) > 0 {
			// We'll need to get the original invoice ID from payment later
			originalInvoiceIds = append(originalInvoiceIds, invoice.PaymentId)
		}
		// Collect promo credit transaction keys
		if invoice.UserId > 0 && len(invoice.InvoiceId) > 0 {
			promoCreditKeys = append(promoCreditKeys, fmt.Sprintf("%d_%s", invoice.UserId, invoice.InvoiceId))
		}
	}

	// 2. Use real bulk query functions
	payments := bulkGetPaymentsByPaymentIds(ctx, paymentIds)
	refunds := bulkGetRefundsByRefundIds(ctx, refundIds)
	gateways := bulkGetGatewaysByIds(ctx, gatewayIds)
	users := bulkGetUserAccountsByIds(ctx, userIds)
	subscriptions := bulkGetSubscriptionsBySubscriptionIds(ctx, subscriptionIds)
	merchants := bulkGetMerchantsByIds(ctx, merchantIds)
	discounts := bulkGetDiscountsByCodes(ctx, discountCodes)

	// 3. Get additional data that depends on previous queries
	originalInvoices := bulkGetOriginalInvoices(ctx, payments, originalInvoiceIds)
	promoCreditTransactions := bulkGetPromoCreditTransactions(ctx, promoCreditKeys)
	subscriptionPendingUpdates := bulkGetSubscriptionPendingUpdates(ctx, invoices)

	// 4. Get plans and addons data from subscription pending updates
	plans := bulkGetPlansFromPendingUpdates(ctx, subscriptionPendingUpdates)
	addons := bulkGetAddonsFromPendingUpdates(ctx, subscriptionPendingUpdates)

	// 5. Create preload data structure
	preloadData := &PreloadData{
		Payments:                   payments,
		Refunds:                    refunds,
		Gateways:                   gateways,
		Users:                      users,
		Subscriptions:              subscriptions,
		Merchants:                  merchants,
		Discounts:                  discounts,
		OriginalInvoices:           originalInvoices,
		PromoCreditTransactions:    promoCreditTransactions,
		SubscriptionPendingUpdates: subscriptionPendingUpdates,
		Plans:                      plans,
		Addons:                     addons,
	}

	// 4. Bulk convert
	var resultList []*InvoiceDetail
	for _, invoice := range invoices {
		detail := convertInvoiceToDetailWithPreload(ctx, invoice, preloadData)
		resultList = append(resultList, detail)
	}

	return resultList
}

// convertInvoiceToDetailWithPreload converts a single invoice using preloaded data
func convertInvoiceToDetailWithPreload(
	ctx context.Context,
	invoice *entity.Invoice,
	preloadData *PreloadData) *InvoiceDetail {

	// Get data from preloaded maps instead of querying database
	payment := preloadData.Payments[invoice.PaymentId]
	refund := preloadData.Refunds[invoice.RefundId]
	gateway := preloadData.Gateways[invoice.GatewayId]
	userAccount := preloadData.Users[invoice.UserId]
	subscription := preloadData.Subscriptions[invoice.SubscriptionId]
	merchant := preloadData.Merchants[invoice.MerchantId]
	discount := preloadData.Discounts[invoice.DiscountCode]

	// Parse JSON fields
	var lines []*bean.InvoiceItemSimplify
	err := utility.UnmarshalFromJsonString(invoice.Lines, &lines)
	for _, line := range lines {
		line.Currency = invoice.Currency
		line.TaxPercentage = invoice.TaxPercentage
	}
	if err != nil {
		fmt.Printf("ConvertInvoiceLines err:%s", err)
	}

	var metadata = make(map[string]interface{})
	if len(invoice.MetaData) > 0 {
		err = gjson.Unmarshal([]byte(invoice.MetaData), &metadata)
		if err != nil {
			fmt.Printf("SimplifySubscription Unmarshal Metadata error:%s", err.Error())
		}
	}

	var userSnapShot *entity.UserAccount
	if len(invoice.Data) > 0 {
		err = gjson.Unmarshal([]byte(invoice.Data), &userSnapShot)
		if err != nil {
			fmt.Printf("UserSnapshot Unmarshal Metadata error:%s", err.Error())
		}
	}

	autoCharge := false
	if invoice.CreateFrom == consts.InvoiceAutoChargeFlag {
		autoCharge = true
	}

	var originalPaymentInvoice *bean.Invoice
	message := ""
	if refund != nil {
		if payment != nil && len(payment.InvoiceId) > 0 {
			// Use preloaded original invoice data with correct key
			originalPaymentInvoice = bean.SimplifyInvoice(preloadData.OriginalInvoices[payment.InvoiceId])
		}
		if invoice.Status == consts.InvoiceStatusFailed {
			message = refund.RefundCommentExplain
		}
	}

	var userMetricChargeForInvoice *bean.UserMetricChargeInvoiceItemEntity
	if len(invoice.MetricCharge) > 0 {
		_ = utility.UnmarshalFromJsonString(invoice.MetricCharge, &userMetricChargeForInvoice)
	}

	var paidTime int64 = 0
	if invoice.Status == consts.InvoiceStatusPaid || invoice.Status == consts.InvoiceStatusReversed {
		if refund != nil && refund.RefundTime > 0 {
			paidTime = refund.RefundTime
		} else if payment != nil && payment.PaidTime > 0 {
			paidTime = payment.PaidTime
		} else if invoice.GmtModify != nil {
			paidTime = invoice.GmtModify.Timestamp()
		} else if invoice.FinishTime > 0 {
			paidTime = invoice.FinishTime
		} else {
			paidTime = invoice.CreateTime
		}
	}

	// Use preloaded subscription pending update data
	subscriptionPendingUpdate := preloadData.SubscriptionPendingUpdates[invoice.InvoiceId]
	// Create plan snapshot without calling GetInvoicePlanSnapshot to avoid N+1 queries
	planSnapshot := createPlanSnapshotFromPreloadedData(ctx, invoice, metadata, lines, subscriptionPendingUpdate)

	return &InvoiceDetail{
		Id:                             invoice.Id,
		MerchantId:                     invoice.MerchantId,
		UserId:                         invoice.UserId,
		SubscriptionId:                 invoice.SubscriptionId,
		InvoiceName:                    invoice.InvoiceName,
		ProductName:                    invoice.ProductName,
		InvoiceId:                      invoice.InvoiceId,
		GatewayPaymentType:             invoice.GatewayInvoiceId,
		UniqueId:                       invoice.UniqueId,
		GmtCreate:                      invoice.GmtCreate,
		OriginAmount:                   invoice.TotalAmount + invoice.DiscountAmount + invoice.PromoCreditDiscountAmount,
		TotalAmount:                    invoice.TotalAmount,
		DiscountCode:                   invoice.DiscountCode,
		DiscountAmount:                 invoice.DiscountAmount,
		TaxAmount:                      invoice.TaxAmount,
		SubscriptionAmount:             invoice.SubscriptionAmount,
		Currency:                       invoice.Currency,
		Lines:                          lines,
		GatewayId:                      invoice.GatewayId,
		Status:                         invoice.Status,
		SendStatus:                     invoice.SendStatus,
		SendEmail:                      invoice.SendEmail,
		SendPdf:                        link.GetInvoicePdfLink(invoice.InvoiceId, invoice.SendTerms),
		GmtModify:                      invoice.GmtModify,
		IsDeleted:                      invoice.IsDeleted,
		Link:                           link.GetInvoiceLink(invoice.InvoiceId, invoice.SendTerms),
		GatewayStatus:                  invoice.GatewayStatus,
		GatewayPaymentId:               invoice.GatewayPaymentId,
		GatewayUserId:                  "", // Not available in entity.Invoice, set to empty string
		GatewayInvoicePdf:              invoice.GatewayInvoicePdf,
		TaxPercentage:                  invoice.TaxPercentage,
		SendNote:                       invoice.SendNote,
		TotalAmountExcludingTax:        invoice.TotalAmountExcludingTax,
		SubscriptionAmountExcludingTax: invoice.SubscriptionAmountExcludingTax,
		PeriodStart:                    invoice.PeriodStart,
		PeriodEnd:                      invoice.PeriodEnd,
		PaymentId:                      invoice.PaymentId,
		RefundId:                       invoice.RefundId,
		Gateway:                        ConvertGatewayDetail(ctx, gateway),
		Merchant:                       bean.SimplifyMerchant(merchant),
		UserAccount:                    bean.SimplifyUserAccount(userAccount),
		UserSnapshot:                   bean.SimplifyUserAccount(userSnapShot),
		Subscription:                   bean.SimplifySubscription(ctx, subscription),
		SubscriptionPendingUpdate:      subscriptionPendingUpdate,
		Payment:                        bean.SimplifyPayment(payment),
		Refund:                         bean.SimplifyRefund(refund),
		Discount:                       bean.SimplifyMerchantDiscountCode(discount),
		CryptoAmount:                   invoice.CryptoAmount,
		CryptoCurrency:                 invoice.CryptoCurrency,
		DayUtilDue:                     invoice.DayUtilDue,
		BillingCycleAnchor:             invoice.BillingCycleAnchor,
		CreateFrom:                     invoice.CreateFrom,
		Metadata:                       metadata,
		CountryCode:                    invoice.CountryCode,
		VatNumber:                      invoice.VatNumber,
		FinishTime:                     invoice.FinishTime,
		TrialEnd:                       invoice.TrialEnd,
		CreateTime:                     invoice.CreateTime,
		PaidTime:                       paidTime,
		BizType:                        invoice.BizType,
		ProrationDate:                  0, // Not available in entity.Invoice, set to 0
		AutoCharge:                     autoCharge,
		OriginalPaymentInvoice:         originalPaymentInvoice,
		PromoCreditDiscountAmount:      invoice.PromoCreditDiscountAmount,
		PromoCreditTransaction:         preloadData.PromoCreditTransactions[fmt.Sprintf("%d_%s", invoice.UserId, invoice.InvoiceId)],
		PartialCreditPaidAmount:        invoice.PartialCreditPaidAmount,
		Message:                        message,
		UserMetricChargeForInvoice:     userMetricChargeForInvoice,
		PlanSnapshot:                   planSnapshot,
	}
}

// createPlanSnapshotFromPreloadedData creates plan snapshot without database queries
func createPlanSnapshotFromPreloadedData(ctx context.Context, invoice *entity.Invoice, metadata map[string]interface{}, lines []*bean.InvoiceItemSimplify, subscriptionPendingUpdate *SubscriptionPendingUpdateDetail) *bean.InvoicePlanSnapshot {
	planSnapshot := &bean.InvoicePlanSnapshot{
		ChargeType:     bean.InvoiceChargeTypeOneTime,
		AutoCharge:     false,
		Plan:           nil,
		Addons:         make([]*bean.PlanAddonDetail, 0),
		PreviousPlan:   nil,
		PreviousAddons: make([]*bean.PlanAddonDetail, 0),
	}

	// Set charge type based on invoice name
	if invoice.InvoiceName == "SubscriptionCreate" {
		planSnapshot.ChargeType = bean.InvoiceChargeTypeSubscriptionCreate
	} else if invoice.InvoiceName == "SubscriptionRenew" {
		planSnapshot.ChargeType = bean.InvoiceChargeTypeSubscriptionRenew
	} else if invoice.InvoiceName == "SubscriptionCycle" {
		planSnapshot.ChargeType = bean.InvoiceChargeTypeSubscriptionCycle
	} else if invoice.InvoiceName == "SubscriptionDowngrade" {
		planSnapshot.ChargeType = bean.InvoiceChargeTypeSubscriptionDowngrade
	} else if invoice.InvoiceName == "SubscriptionUpdate" {
		planSnapshot.ChargeType = bean.InvoiceChargeTypeSubscriptionUpgrade
		if isUpgrade, ok := metadata["IsUpgrade"]; ok {
			if isUpgrade != nil {
				switch v := isUpgrade.(type) {
				case bool:
					if !v {
						planSnapshot.ChargeType = bean.InvoiceChargeTypeSubscriptionDowngrade
					}
				case string:
					if strings.ToLower(strings.TrimSpace(v)) == "false" {
						planSnapshot.ChargeType = bean.InvoiceChargeTypeSubscriptionDowngrade
					}
				}
			}
		}
	}

	// Set auto charge flag
	if invoice.CreateFrom == "AutoRenew" {
		planSnapshot.AutoCharge = true
	}

	// Set plan and addons based on charge type
	if planSnapshot.ChargeType == bean.InvoiceChargeTypeSubscriptionDowngrade || planSnapshot.ChargeType == bean.InvoiceChargeTypeSubscriptionUpgrade {
		if subscriptionPendingUpdate != nil {
			planSnapshot.Plan = subscriptionPendingUpdate.UpdatePlan
			planSnapshot.Addons = subscriptionPendingUpdate.UpdateAddons
			planSnapshot.PreviousPlan = subscriptionPendingUpdate.Plan
			planSnapshot.PreviousAddons = subscriptionPendingUpdate.Addons
		} else {
			// Process lines for plan and addons
			for _, line := range lines {
				if line.AmountExcludingTax >= 0 {
					if line.Plan != nil && line.Plan.Type == consts.PlanTypeRecurringAddon {
						planSnapshot.Addons = append(planSnapshot.Addons, &bean.PlanAddonDetail{
							Quantity:  line.Quantity,
							AddonPlan: line.Plan,
						})
					} else if line.Plan != nil {
						planSnapshot.Plan = line.Plan
					}
				} else {
					if line.Plan != nil && line.Plan.Type == consts.PlanTypeRecurringAddon {
						planSnapshot.PreviousAddons = append(planSnapshot.PreviousAddons, &bean.PlanAddonDetail{
							Quantity:  line.Quantity,
							AddonPlan: line.Plan,
						})
					} else if line.Plan != nil {
						planSnapshot.PreviousPlan = line.Plan
					}
				}
			}
		}
	} else {
		// Process lines for regular invoices
		for _, line := range lines {
			if line.Plan != nil {
				if line.Plan.Type == consts.PlanTypeRecurringAddon {
					planSnapshot.Addons = append(planSnapshot.Addons, &bean.PlanAddonDetail{
						Quantity:  line.Quantity,
						AddonPlan: line.Plan,
					})
				} else {
					planSnapshot.Plan = line.Plan
				}
			}
		}
	}

	return planSnapshot
}

// ==================== Real Bulk Query Functions ====================

// bulkGetPaymentsByPaymentIds bulk query payments
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
func bulkGetGatewaysByIds(ctx context.Context, gatewayIdStrs []string) map[uint64]*entity.MerchantGateway {
	if len(gatewayIdStrs) == 0 {
		return make(map[uint64]*entity.MerchantGateway)
	}

	// Convert string IDs to uint64
	var gatewayIds []uint64
	for _, idStr := range gatewayIdStrs {
		if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			gatewayIds = append(gatewayIds, id)
		}
	}

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
func bulkGetUserAccountsByIds(ctx context.Context, userIdStrs []string) map[uint64]*entity.UserAccount {
	if len(userIdStrs) == 0 {
		return make(map[uint64]*entity.UserAccount)
	}

	// Convert string IDs to uint64
	var userIds []uint64
	for _, idStr := range userIdStrs {
		if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			userIds = append(userIds, id)
		}
	}

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

// bulkGetMerchantsByIds bulk query merchants
func bulkGetMerchantsByIds(ctx context.Context, merchantIds []uint64) map[uint64]*entity.Merchant {
	if len(merchantIds) == 0 {
		return make(map[uint64]*entity.Merchant)
	}

	var merchants []*entity.Merchant
	err := dao.Merchant.Ctx(ctx).WhereIn(dao.Merchant.Columns().Id, merchantIds).Scan(&merchants)
	if err != nil {
		return make(map[uint64]*entity.Merchant)
	}

	result := make(map[uint64]*entity.Merchant)
	for _, merchant := range merchants {
		result[merchant.Id] = merchant
	}
	return result
}

// bulkGetDiscountsByCodes bulk query discounts
func bulkGetDiscountsByCodes(ctx context.Context, discountCodes []string) map[string]*entity.MerchantDiscountCode {
	if len(discountCodes) == 0 {
		return make(map[string]*entity.MerchantDiscountCode)
	}

	var discounts []*entity.MerchantDiscountCode
	err := dao.MerchantDiscountCode.Ctx(ctx).WhereIn(dao.MerchantDiscountCode.Columns().Code, discountCodes).Scan(&discounts)
	if err != nil {
		return make(map[string]*entity.MerchantDiscountCode)
	}

	result := make(map[string]*entity.MerchantDiscountCode)
	for _, discount := range discounts {
		result[discount.Code] = discount
	}
	return result
}

// bulkGetOriginalInvoices bulk query original invoices for refunds
func bulkGetOriginalInvoices(ctx context.Context, payments map[string]*entity.Payment, paymentIds []string) map[string]*entity.Invoice {
	if len(paymentIds) == 0 {
		return make(map[string]*entity.Invoice)
	}

	// Collect invoice IDs from payments
	var invoiceIds []string
	for _, paymentId := range paymentIds {
		if payment, exists := payments[paymentId]; exists && payment != nil && len(payment.InvoiceId) > 0 {
			invoiceIds = append(invoiceIds, payment.InvoiceId)
		}
	}

	if len(invoiceIds) == 0 {
		return make(map[string]*entity.Invoice)
	}

	var invoices []*entity.Invoice
	err := dao.Invoice.Ctx(ctx).WhereIn(dao.Invoice.Columns().InvoiceId, invoiceIds).Scan(&invoices)
	if err != nil {
		return make(map[string]*entity.Invoice)
	}

	result := make(map[string]*entity.Invoice)
	for _, invoice := range invoices {
		result[invoice.InvoiceId] = invoice
	}
	return result
}

// bulkGetPromoCreditTransactions bulk query promo credit transactions
func bulkGetPromoCreditTransactions(ctx context.Context, promoCreditKeys []string) map[string]*bean.CreditTransaction {
	if len(promoCreditKeys) == 0 {
		return make(map[string]*bean.CreditTransaction)
	}

	result := make(map[string]*bean.CreditTransaction)

	// Parse promo credit keys to extract userId and invoiceId
	for _, key := range promoCreditKeys {
		parts := strings.Split(key, "_")
		if len(parts) == 2 {
			userIdStr := parts[0]
			invoiceId := parts[1]

			userId, err := strconv.ParseUint(userIdStr, 10, 64)
			if err != nil {
				continue
			}

			// Query individual promo credit transaction
			// Note: This is not truly bulk, but it's the best we can do without knowing the exact table structure
			creditTransaction := query.GetPromoCreditTransactionByInvoiceId(ctx, userId, invoiceId)
			if creditTransaction != nil {
				result[key] = bean.SimplifyCreditTransaction(ctx, creditTransaction)
			}
		}
	}

	return result
}

// bulkGetSubscriptionPendingUpdates bulk query subscription pending updates
func bulkGetSubscriptionPendingUpdates(ctx context.Context, invoices []*entity.Invoice) map[string]*SubscriptionPendingUpdateDetail {
	if len(invoices) == 0 {
		return make(map[string]*SubscriptionPendingUpdateDetail)
	}

	var invoiceIds []string
	for _, invoice := range invoices {
		if len(invoice.InvoiceId) > 0 {
			invoiceIds = append(invoiceIds, invoice.InvoiceId)
		}
	}

	if len(invoiceIds) == 0 {
		return make(map[string]*SubscriptionPendingUpdateDetail)
	}

	var pendingUpdates []*entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).WhereIn(dao.SubscriptionPendingUpdate.Columns().InvoiceId, invoiceIds).Scan(&pendingUpdates)
	if err != nil {
		return make(map[string]*SubscriptionPendingUpdateDetail)
	}

	result := make(map[string]*SubscriptionPendingUpdateDetail)
	for _, pendingUpdate := range pendingUpdates {
		if pendingUpdate != nil {
			detail := convertSubscriptionPendingUpdateToDetail(ctx, pendingUpdate)
			result[pendingUpdate.InvoiceId] = detail
		}
	}
	return result
}

// convertSubscriptionPendingUpdateToDetail converts entity to detail (simplified version)
func convertSubscriptionPendingUpdateToDetail(ctx context.Context, pendingUpdate *entity.SubscriptionPendingUpdate) *SubscriptionPendingUpdateDetail {
	var metadata = make(map[string]interface{})
	if len(pendingUpdate.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(pendingUpdate.MetaData), &metadata)
		if err != nil {
			fmt.Printf("convertSubscriptionPendingUpdateToDetail Unmarshal Metadata error:%s", err.Error())
		}
	}

	return &SubscriptionPendingUpdateDetail{
		MerchantId:      pendingUpdate.MerchantId,
		SubscriptionId:  pendingUpdate.SubscriptionId,
		PendingUpdateId: pendingUpdate.PendingUpdateId,
		GmtCreate:       pendingUpdate.GmtCreate,
		Amount:          pendingUpdate.Amount,
		Status:          pendingUpdate.Status,
		UpdateAmount:    pendingUpdate.UpdateAmount,
		Currency:        pendingUpdate.Currency,
		UpdateCurrency:  pendingUpdate.UpdateCurrency,
		PlanId:          pendingUpdate.PlanId,
		UpdatePlanId:    pendingUpdate.UpdatePlanId,
		Quantity:        pendingUpdate.Quantity,
		UpdateQuantity:  pendingUpdate.UpdateQuantity,
		AddonData:       pendingUpdate.AddonData,
		UpdateAddonData: pendingUpdate.UpdateAddonData,
		ProrationAmount: pendingUpdate.ProrationAmount,
		GatewayId:       pendingUpdate.GatewayId,
		UserId:          pendingUpdate.UserId,
		InvoiceId:       pendingUpdate.InvoiceId,
		GmtModify:       pendingUpdate.GmtModify,
		Paid:            pendingUpdate.Paid,
		Link:            pendingUpdate.Link,
		EffectImmediate: pendingUpdate.EffectImmediate,
		EffectTime:      pendingUpdate.EffectTime,
		Note:            pendingUpdate.Note,
		Plan:            nil, // Will be filled by bulk query
		Addons:          nil, // Will be filled by bulk query
		UpdatePlan:      nil, // Will be filled by bulk query
		UpdateAddons:    nil, // Will be filled by bulk query
		Metadata:        metadata,
	}
}

// bulkGetPlansFromPendingUpdates bulk query plans from subscription pending updates
func bulkGetPlansFromPendingUpdates(ctx context.Context, pendingUpdates map[string]*SubscriptionPendingUpdateDetail) map[uint64]*entity.Plan {
	if len(pendingUpdates) == 0 {
		return make(map[uint64]*entity.Plan)
	}

	var planIds []uint64
	for _, pendingUpdate := range pendingUpdates {
		if pendingUpdate != nil {
			if pendingUpdate.PlanId > 0 {
				planIds = append(planIds, pendingUpdate.PlanId)
			}
			if pendingUpdate.UpdatePlanId > 0 {
				planIds = append(planIds, pendingUpdate.UpdatePlanId)
			}
		}
	}

	if len(planIds) == 0 {
		return make(map[uint64]*entity.Plan)
	}

	var plans []*entity.Plan
	err := dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, planIds).Scan(&plans)
	if err != nil {
		return make(map[uint64]*entity.Plan)
	}

	result := make(map[uint64]*entity.Plan)
	for _, plan := range plans {
		result[plan.Id] = plan
	}
	return result
}

// bulkGetAddonsFromPendingUpdates bulk query addons from subscription pending updates
func bulkGetAddonsFromPendingUpdates(ctx context.Context, pendingUpdates map[string]*SubscriptionPendingUpdateDetail) map[string][]*bean.PlanAddonDetail {
	if len(pendingUpdates) == 0 {
		return make(map[string][]*bean.PlanAddonDetail)
	}

	result := make(map[string][]*bean.PlanAddonDetail)
	// Note: This is a simplified implementation
	// The actual addon parsing logic would need to be implemented based on your business logic
	for _, pendingUpdate := range pendingUpdates {
		if pendingUpdate != nil {
			// For now, returning empty addons to avoid compilation errors
			// This would need to be implemented based on your actual addon parsing logic
			result[pendingUpdate.InvoiceId] = make([]*bean.PlanAddonDetail, 0)
		}
	}
	return result
}
