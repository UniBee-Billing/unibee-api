package merchant

import (
	"context"
	"unibee/api/merchant/discount"
	_interface "unibee/internal/interface/context"
	discount2 "unibee/internal/logic/discount"
)

func (c *ControllerDiscount) List(ctx context.Context, req *discount.ListReq) (res *discount.ListRes, err error) {
	list, total := discount2.MerchantDiscountCodeList(ctx, &discount2.ListInternalReq{
		MerchantId:   _interface.GetMerchantId(ctx),
		DiscountType: req.DiscountType,
		BillingType:  req.BillingType,
		Status:       req.Status,
		Code:         req.Code,
		SearchKey:    req.SearchKey,
		Currency:     req.Currency,
		SortField:    req.SortField,
		SortType:     req.SortType,
		Page:         req.Page,
		Count:        req.Count,
		CreateTimeStart: req.CreateTimeStart,
		CreateTimeEnd:   req.CreateTimeEnd,
	})
	return &discount.ListRes{Discounts: list, Total: total}, nil
}
