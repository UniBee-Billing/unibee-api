// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package checkout

import (
	"context"

	"unibee/api/checkout/checkout"
	"unibee/api/checkout/gateway"
	"unibee/api/checkout/ip"
	"unibee/api/checkout/payment"
	"unibee/api/checkout/plan"
	"unibee/api/checkout/subscription"
	"unibee/api/checkout/translater"
	"unibee/api/checkout/vat"
)

type ICheckoutCheckout interface {
	Get(ctx context.Context, req *checkout.GetReq) (res *checkout.GetRes, err error)
}

type ICheckoutGateway interface {
	List(ctx context.Context, req *gateway.ListReq) (res *gateway.ListRes, err error)
}

type ICheckoutIp interface {
	Resolve(ctx context.Context, req *ip.ResolveReq) (res *ip.ResolveRes, err error)
}

type ICheckoutPayment interface {
	Detail(ctx context.Context, req *payment.DetailReq) (res *payment.DetailRes, err error)
}

type ICheckoutPlan interface {
	Detail(ctx context.Context, req *plan.DetailReq) (res *plan.DetailRes, err error)
}

type ICheckoutSubscription interface {
	CreatePreview(ctx context.Context, req *subscription.CreatePreviewReq) (res *subscription.CreatePreviewRes, err error)
	Create(ctx context.Context, req *subscription.CreateReq) (res *subscription.CreateRes, err error)
}

type ICheckoutTranslater interface {
	Translate(ctx context.Context, req *translater.TranslateReq) (res *translater.TranslateRes, err error)
}

type ICheckoutVat interface {
	CountryList(ctx context.Context, req *vat.CountryListReq) (res *vat.CountryListRes, err error)
	NumberValidate(ctx context.Context, req *vat.NumberValidateReq) (res *vat.NumberValidateRes, err error)
}
