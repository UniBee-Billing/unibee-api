package user

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/payment/service"

	"unibee/api/user/payment"
)

func (c *ControllerPayment) ItemList(ctx context.Context, req *payment.ItemListReq) (res *payment.ItemListRes, err error) {
	result, err := service.OneTimePaymentItemList(ctx, &service.PaymentItemListInternalReq{
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
	return &payment.ItemListRes{PaymentItems: result.PaymentItems, Total: result.Total}, nil
}
