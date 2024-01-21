package user

import (
	"context"
	"go-oversea-pay/internal/logic/payment/service"

	"go-oversea-pay/api/user/payment"
)

func (c *ControllerPayment) TimeLineList(ctx context.Context, req *payment.TimeLineListReq) (res *payment.TimeLineListRes, err error) {
	result, err := service.PaymentTimeLineList(ctx, &service.PaymentTimelineListInternalReq{
		MerchantId: req.MerchantId,
		UserId:     req.UserId,
		SortField:  req.SortField,
		SortType:   req.SortType,
		Page:       req.Page,
		Count:      req.Count,
	})
	if err != nil {
		return nil, err
	}
	return &payment.TimeLineListRes{PaymentTimelines: result.PaymentTimelines}, nil
}
