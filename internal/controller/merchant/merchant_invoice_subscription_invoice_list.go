package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/invoice/service"
	"unibee-api/utility"

	merchantInvoice "unibee-api/api/merchant/invoice"
)

func (c *ControllerInvoice) SubscriptionInvoiceList(ctx context.Context, req *merchantInvoice.SubscriptionInvoiceListReq) (res *merchantInvoice.SubscriptionInvoiceListRes, err error) {
	//Merchant 权限检查

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
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
