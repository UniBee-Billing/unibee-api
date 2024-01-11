package merchant

import (
	"context"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/subscription/service"
	"go-oversea-pay/utility"

	merchantInvoice "go-oversea-pay/api/merchant/invoice"
)

func (c *ControllerInvoice) SubscriptionInvoiceList(ctx context.Context, req *merchantInvoice.SubscriptionInvoiceListReq) (res *merchantInvoice.SubscriptionInvoiceListRes, err error) {
	//Merchant 权限检查

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).Merchant != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).Merchant.Id > 0, "merchantUserId invalid")
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
	return &merchantInvoice.SubscriptionInvoiceListRes{Invoices: internalResult.Invoices}, nil
}
