// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package checkout

import (
	"context"

	"unibee/api/checkout/gateway"
	"unibee/api/checkout/ip"
	"unibee/api/checkout/subscription"
	"unibee/api/checkout/vat"
)

type ICheckoutGateway interface {
	List(ctx context.Context, req *gateway.ListReq) (res *gateway.ListRes, err error)
}

type ICheckoutIp interface {
	Resolve(ctx context.Context, req *ip.ResolveReq) (res *ip.ResolveRes, err error)
}

type ICheckoutSubscription interface {
	CreatePreview(ctx context.Context, req *subscription.CreatePreviewReq) (res *subscription.CreatePreviewRes, err error)
	Create(ctx context.Context, req *subscription.CreateReq) (res *subscription.CreateRes, err error)
}

type ICheckoutVat interface {
	CountryList(ctx context.Context, req *vat.CountryListReq) (res *vat.CountryListRes, err error)
	NumberValidate(ctx context.Context, req *vat.NumberValidateReq) (res *vat.NumberValidateRes, err error)
}
