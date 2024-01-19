package user

import (
	"context"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/subscription/service"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/user/invoice"
)

func (c *ControllerInvoice) SubscriptionInvoiceList(ctx context.Context, req *invoice.SubscriptionInvoiceListReq) (res *invoice.SubscriptionInvoiceListRes, err error) {
	//Merchant 权限检查

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "user auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).User.Id > 0, "userId invalid")
	}

	internalResult, err := service.SubscriptionInvoiceList(ctx, &service.SubscriptionInvoiceListInternalReq{
		MerchantId:    req.MerchantId,
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
