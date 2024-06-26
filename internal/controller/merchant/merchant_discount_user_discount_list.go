package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	discount2 "unibee/internal/logic/discount"

	"unibee/api/merchant/discount"
)

func (c *ControllerDiscount) UserDiscountList(ctx context.Context, req *discount.UserDiscountListReq) (res *discount.UserDiscountListRes, err error) {
	list, total := discount2.MerchantUserDiscountCodeList(ctx, &discount2.UserDiscountListInternalReq{
		MerchantId:      _interface.GetMerchantId(ctx),
		Id:              req.Id,
		SortField:       req.SortField,
		SortType:        req.SortType,
		Page:            req.Page,
		Count:           req.Count,
		CreateTimeStart: req.CreateTimeStart,
		CreateTimeEnd:   req.CreateTimeEnd,
	})
	return &discount.UserDiscountListRes{
		UserDiscounts: list,
		Total:         total,
	}, nil
}
