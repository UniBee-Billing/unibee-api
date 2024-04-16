package user

import (
	"context"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/invoice/service"
	"unibee/utility"

	"unibee/api/user/invoice"
)

func (c *ControllerInvoice) List(ctx context.Context, req *invoice.ListReq) (res *invoice.ListRes, err error) {
	if !config.GetConfigInstance().IsLocal() {
		utility.Assert(_interface.Context().Get(ctx).User != nil, "user auth failure,not login")
		utility.Assert(_interface.Context().Get(ctx).User.Id > 0, "userId invalid")
	}

	internalResult, err := service.InvoiceList(ctx, &service.InvoiceListInternalReq{
		MerchantId: _interface.GetMerchantId(ctx),
		UserId:     _interface.Context().Get(ctx).User.Id,
		SortField:  req.SortField,
		SortType:   req.SortType,
		Page:       req.Page,
		Count:      req.Count,
	})
	if err != nil {
		return nil, err
	}
	return &invoice.ListRes{Invoices: internalResult.Invoices}, nil
}
