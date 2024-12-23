package user

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/payment/service"

	"unibee/api/user/payment"
)

func (c *ControllerPayment) TimeLineList(ctx context.Context, req *payment.TimeLineListReq) (res *payment.TimeLineListRes, err error) {
	result, err := service.PaymentTimeLineList(ctx, &service.PaymentTimelineListInternalReq{
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
	return &payment.TimeLineListRes{PaymentTimelines: result.PaymentTimelines, Total: result.Total}, nil
}
