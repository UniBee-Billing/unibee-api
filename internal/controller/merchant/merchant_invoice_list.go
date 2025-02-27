package merchant

import (
	"context"
	merchantInvoice "unibee/api/merchant/invoice"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/invoice/service"
)

func (c *ControllerInvoice) List(ctx context.Context, req *merchantInvoice.ListReq) (res *merchantInvoice.ListRes, err error) {
	internalResult, err := service.InvoiceList(ctx, &service.InvoiceListInternalReq{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Currency:        req.Currency,
		Status:          req.Status,
		AmountStart:     req.AmountStart,
		AmountEnd:       req.AmountEnd,
		MerchantId:      _interface.GetMerchantId(ctx),
		UserId:          req.UserId,
		SendEmail:       req.SendEmail,
		SortField:       req.SortField,
		SortType:        req.SortType,
		DeleteInclude:   req.DeleteInclude,
		Type:            req.Type,
		Page:            req.Page,
		Count:           req.Count,
		CreateTimeStart: req.CreateTimeStart,
		CreateTimeEnd:   req.CreateTimeEnd,
	})
	if err != nil {
		return nil, err
	}
	return &merchantInvoice.ListRes{Invoices: internalResult.Invoices, Total: internalResult.Total}, nil
}
