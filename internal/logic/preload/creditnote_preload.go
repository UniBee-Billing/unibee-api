package preload

import (
	"context"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

// CreditNotePreloadData holds all preloaded data for credit notes
type CreditNotePreloadData struct {
	Gateways      map[uint64]*entity.MerchantGateway
	Users         map[uint64]*entity.UserAccount
	Payments      map[string]*entity.Payment
	Refunds       map[string]*entity.Refund
	Subscriptions map[string]*entity.Subscription
	Discounts     map[string]*entity.MerchantDiscountCode
}

func CreditNoteListPreloadForContext(ctx context.Context, creditNotes []*entity.Invoice) {
	if len(creditNotes) > 3 {
		preloadData := CreditNoteListPreload(ctx, creditNotes)
		if preloadData != nil && _interface.GetBulkPreloadData(ctx) != nil {
			_interface.Context().Get(ctx).PreloadData.Gateways = preloadData.Gateways
			_interface.Context().Get(ctx).PreloadData.Users = preloadData.Users
			_interface.Context().Get(ctx).PreloadData.Payments = preloadData.Payments
			_interface.Context().Get(ctx).PreloadData.Refunds = preloadData.Refunds
			_interface.Context().Get(ctx).PreloadData.Subscriptions = preloadData.Subscriptions
			_interface.Context().Get(ctx).PreloadData.Discounts = preloadData.Discounts
		}
	}
}

// CreditNoteListPreload preloads all related data for credit notes in bulk
func CreditNoteListPreload(ctx context.Context, creditNotes []*entity.Invoice) *CreditNotePreloadData {
	preload := &CreditNotePreloadData{
		Gateways:      make(map[uint64]*entity.MerchantGateway),
		Users:         make(map[uint64]*entity.UserAccount),
		Payments:      make(map[string]*entity.Payment),
		Refunds:       make(map[string]*entity.Refund),
		Subscriptions: make(map[string]*entity.Subscription),
		Discounts:     make(map[string]*entity.MerchantDiscountCode),
	}

	if len(creditNotes) == 0 {
		return preload
	}

	// Collect unique IDs
	gatewayIds := make(map[uint64]bool)
	userIds := make(map[uint64]bool)
	paymentIds := make(map[string]bool)
	refundIds := make(map[string]bool)
	subscriptionIds := make(map[string]bool)
	discountCodes := make(map[string]bool)

	for _, creditNote := range creditNotes {
		if creditNote.GatewayId > 0 {
			gatewayIds[creditNote.GatewayId] = true
		}
		if creditNote.UserId > 0 {
			userIds[creditNote.UserId] = true
		}
		if len(creditNote.PaymentId) > 0 {
			paymentIds[creditNote.PaymentId] = true
		}
		if len(creditNote.RefundId) > 0 {
			refundIds[creditNote.RefundId] = true
		}
		if len(creditNote.SubscriptionId) > 0 {
			subscriptionIds[creditNote.SubscriptionId] = true
		}
		if len(creditNote.DiscountCode) > 0 {
			discountCodes[creditNote.DiscountCode] = true
		}
	}

	// Bulk query gateways
	if len(gatewayIds) > 0 {
		gatewayIdList := make([]uint64, 0, len(gatewayIds))
		for id := range gatewayIds {
			gatewayIdList = append(gatewayIdList, id)
		}

		var gateways []*entity.MerchantGateway
		err := dao.MerchantGateway.Ctx(ctx).WhereIn(dao.MerchantGateway.Columns().Id, gatewayIdList).Scan(&gateways)
		if err == nil {
			for _, gateway := range gateways {
				preload.Gateways[gateway.Id] = gateway
			}
		}
	}

	// Bulk query users
	if len(userIds) > 0 {
		userIdList := make([]uint64, 0, len(userIds))
		for id := range userIds {
			userIdList = append(userIdList, id)
		}

		var users []*entity.UserAccount
		err := dao.UserAccount.Ctx(ctx).WhereIn(dao.UserAccount.Columns().Id, userIdList).Scan(&users)
		if err == nil {
			for _, user := range users {
				preload.Users[user.Id] = user
			}
		}
	}

	// Bulk query payments
	if len(paymentIds) > 0 {
		paymentIdList := make([]string, 0, len(paymentIds))
		for id := range paymentIds {
			paymentIdList = append(paymentIdList, id)
		}

		var payments []*entity.Payment
		err := dao.Payment.Ctx(ctx).WhereIn(dao.Payment.Columns().PaymentId, paymentIdList).Scan(&payments)
		if err == nil {
			for _, payment := range payments {
				preload.Payments[payment.PaymentId] = payment
			}
		}
	}

	// Bulk query refunds
	if len(refundIds) > 0 {
		refundIdList := make([]string, 0, len(refundIds))
		for id := range refundIds {
			refundIdList = append(refundIdList, id)
		}

		var refunds []*entity.Refund
		err := dao.Refund.Ctx(ctx).WhereIn(dao.Refund.Columns().RefundId, refundIdList).Scan(&refunds)
		if err == nil {
			for _, refund := range refunds {
				preload.Refunds[refund.RefundId] = refund
			}
		}
	}

	// Bulk query subscriptions
	if len(subscriptionIds) > 0 {
		subscriptionIdList := make([]string, 0, len(subscriptionIds))
		for id := range subscriptionIds {
			subscriptionIdList = append(subscriptionIdList, id)
		}

		var subscriptions []*entity.Subscription
		err := dao.Subscription.Ctx(ctx).WhereIn(dao.Subscription.Columns().SubscriptionId, subscriptionIdList).Scan(&subscriptions)
		if err == nil {
			for _, subscription := range subscriptions {
				preload.Subscriptions[subscription.SubscriptionId] = subscription
			}
		}
	}

	// Bulk query discounts
	if len(discountCodes) > 0 {
		discountCodeList := make([]string, 0, len(discountCodes))
		for code := range discountCodes {
			discountCodeList = append(discountCodeList, code)
		}

		// Note: This might need to be adjusted based on the actual discount query structure
		for _, code := range discountCodeList {
			if discount := query.GetDiscountByCode(ctx, creditNotes[0].MerchantId, code); discount != nil {
				preload.Discounts[code] = discount
			}
		}
	}

	return preload
}
