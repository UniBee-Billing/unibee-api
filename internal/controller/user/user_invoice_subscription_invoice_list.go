package user

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/invoice/service"
	"unibee/utility"

	"unibee/api/user/invoice"
)

func (c *ControllerInvoice) SubscriptionInvoiceList(ctx context.Context, req *invoice.SubscriptionInvoiceListReq) (res *invoice.SubscriptionInvoiceListRes, err error) {
	//Merchant 权限检查

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "user auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).User.Id > 0, "userId invalid")
	}

	internalResult, err := service.SubscriptionInvoiceList(ctx, &service.SubscriptionInvoiceListInternalReq{
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
	return &invoice.SubscriptionInvoiceListRes{Invoices: internalResult.Invoices}, nil
}
