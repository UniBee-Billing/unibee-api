package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/vat_gateway"

	"unibee/api/merchant/vat"
)

func (c *ControllerVat) NumberValidateHistory(ctx context.Context, req *vat.NumberValidateHistoryReq) (res *vat.NumberValidateHistoryRes, err error) {
	list, total := vat_gateway.MerchantNumberValidateHistoryList(ctx, &vat_gateway.NumberValidateHistoryInternalReq{
		MerchantId:      _interface.Context().Get(ctx).MerchantId,
		SearchKey:       req.SearchKey,
		VatNumber:       req.VatNumber,
		CountryCode:     req.CountryCode,
		ValidateGateway: req.ValidateGateway,
		Status:          req.Status,
		SortField:       req.SortField,
		SortType:        req.SortType,
		Page:            req.Page,
		Count:           req.Count,
		CreateTimeStart: req.CreateTimeStart,
		CreateTimeEnd:   req.CreateTimeEnd,
	})
	return &vat.NumberValidateHistoryRes{
		NumberValidateHistoryList: list,
		Total:                     total,
	}, nil
}
