package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/invoice/service"
	"unibee/utility"

	merchantInvoice "unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) List(ctx context.Context, req *merchantInvoice.ListReq) (res *merchantInvoice.ListRes, err error) {
	//Merchant 权限检查

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}

	internalResult, err := service.SubscriptionInvoiceList(ctx, &service.SubscriptionInvoiceListInternalReq{
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Currency:      req.Currency,
		Status:        req.Status,
		AmountStart:   req.AmountStart,
		AmountEnd:     req.AmountEnd,
		MerchantId:    _interface.GetMerchantId(ctx),
		UserId:        req.UserId,
		SendEmail:     req.SendEmail,
		SortField:     req.SortField,
		SortType:      req.SortType,
		DeleteInclude: req.DeleteInclude,
		Page:          req.Page,
		Count:         req.Count,
	})
	if err != nil {
		return nil, err
	}
	return &merchantInvoice.ListRes{Invoices: internalResult.Invoices}, nil
}
