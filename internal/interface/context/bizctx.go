package context

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee/internal/model"
	entity "unibee/internal/model/entity/default"
)

type IContext interface {
	Init(r *ghttp.Request, customCtx *model.Context)
	CopyFromCtx(sourceCtx context.Context) (targetCtx context.Context)
	Get(ctx context.Context) *model.Context
	SetUser(ctx context.Context, ctxUser *model.ContextUser)
	SetMerchantMember(ctx context.Context, ctxMerchantMember *model.ContextMerchantMember)
	SetData(ctx context.Context, data g.Map)
}

var singleTonContext IContext

func Context() IContext {
	if singleTonContext == nil {
		panic("implement not found for interface IContext, forgot register?")
	}
	return singleTonContext
}

const (
	SystemAssertPrefix = "system_assert: "
)

func GetMerchantId(ctx context.Context) uint64 {
	if Context().Get(ctx).MerchantId <= 0 {
		panic(SystemAssertPrefix + "Invalid Merchant")
	}
	return Context().Get(ctx).MerchantId
}

func GetBulkPreloadData(ctx context.Context) *model.PreloadData {
	if Context() != nil && Context().Get(ctx) != nil {
		if Context().Get(ctx).PreloadData == nil {
			Context().Get(ctx).PreloadData = &model.PreloadData{
				Plans:                   make(map[uint64]*entity.Plan),
				Products:                make(map[uint64]*entity.Product),
				Payments:                make(map[string]*entity.Payment),
				Refunds:                 make(map[string]*entity.Refund),
				Invoices:                make(map[string]*entity.Invoice),
				Gateways:                make(map[uint64]*entity.MerchantGateway),
				Users:                   make(map[uint64]*entity.UserAccount),
				Subscriptions:           make(map[string]*entity.Subscription),
				PromoCreditTransactions: make(map[string]*entity.CreditTransaction),
			}
		}
		return Context().Get(ctx).PreloadData
	}
	return nil
}

func GetPlanFromPreloadContext(ctx context.Context, planId uint64) (one *entity.Plan) {
	if GetBulkPreloadData(ctx) != nil && GetBulkPreloadData(ctx).Plans != nil {
		if _, ok := GetBulkPreloadData(ctx).Plans[planId]; ok {
			return GetBulkPreloadData(ctx).Plans[planId]
		}
	}
	return nil
}

func GetProductFromPreloadContext(ctx context.Context, productId uint64) (one *entity.Product) {
	if GetBulkPreloadData(ctx) != nil && GetBulkPreloadData(ctx).Products != nil {
		if _, ok := GetBulkPreloadData(ctx).Products[productId]; ok {
			return GetBulkPreloadData(ctx).Products[productId]
		}
	}
	return nil
}

func GetPaymentFromPreloadContext(ctx context.Context, paymentId string) (one *entity.Payment) {
	if GetBulkPreloadData(ctx) != nil && GetBulkPreloadData(ctx).Payments != nil {
		if _, ok := GetBulkPreloadData(ctx).Payments[paymentId]; ok {
			return GetBulkPreloadData(ctx).Payments[paymentId]
		}
	}
	return nil
}

func GetRefundFromPreloadContext(ctx context.Context, refundId string) (one *entity.Refund) {
	if GetBulkPreloadData(ctx) != nil && GetBulkPreloadData(ctx).Refunds != nil {
		if _, ok := GetBulkPreloadData(ctx).Refunds[refundId]; ok {
			return GetBulkPreloadData(ctx).Refunds[refundId]
		}
	}
	return nil
}

func GetGatewayFromPreloadContext(ctx context.Context, gatewayId uint64) (one *entity.MerchantGateway) {
	if GetBulkPreloadData(ctx) != nil && GetBulkPreloadData(ctx).Gateways != nil {
		if _, ok := GetBulkPreloadData(ctx).Gateways[gatewayId]; ok {
			return GetBulkPreloadData(ctx).Gateways[gatewayId]
		}
	}
	return nil
}

func GetInvoiceFromPreloadContext(ctx context.Context, invoiceId string) (one *entity.Invoice) {
	if GetBulkPreloadData(ctx) != nil && GetBulkPreloadData(ctx).Invoices != nil {
		if _, ok := GetBulkPreloadData(ctx).Invoices[invoiceId]; ok {
			return GetBulkPreloadData(ctx).Invoices[invoiceId]
		}
	}
	return nil
}

func GetUserFromPreloadContext(ctx context.Context, userId uint64) (one *entity.UserAccount) {
	if GetBulkPreloadData(ctx) != nil && GetBulkPreloadData(ctx).Users != nil {
		if _, ok := GetBulkPreloadData(ctx).Users[userId]; ok {
			return GetBulkPreloadData(ctx).Users[userId]
		}
	}
	return nil
}

func GetSubscriptionFromPreloadContext(ctx context.Context, subscriptionId string) (one *entity.Subscription) {
	if GetBulkPreloadData(ctx) != nil && GetBulkPreloadData(ctx).Subscriptions != nil {
		if _, ok := GetBulkPreloadData(ctx).Subscriptions[subscriptionId]; ok {
			return GetBulkPreloadData(ctx).Subscriptions[subscriptionId]
		}
	}
	return nil
}

func GetPromoCreditTransactionFromPreloadContext(ctx context.Context, invoiceId string) (one *entity.CreditTransaction) {
	if GetBulkPreloadData(ctx) != nil && GetBulkPreloadData(ctx).PromoCreditTransactions != nil {
		if _, ok := GetBulkPreloadData(ctx).PromoCreditTransactions[invoiceId]; ok {
			return GetBulkPreloadData(ctx).PromoCreditTransactions[invoiceId]
		}
	}
	return nil
}

func GetDiscountCodeFromPreloadContext(ctx context.Context, code string) (one *entity.MerchantDiscountCode) {
	if GetBulkPreloadData(ctx) != nil && GetBulkPreloadData(ctx).Discounts != nil {
		if _, ok := GetBulkPreloadData(ctx).Discounts[code]; ok {
			return GetBulkPreloadData(ctx).Discounts[code]
		}
	}
	return nil
}

func RegisterContext(i IContext) {
	singleTonContext = i
}
