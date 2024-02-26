package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/payment/service"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) TimeLineList(ctx context.Context, req *payment.TimeLineListReq) (res *payment.TimeLineListRes, err error) {
	result, err := service.PaymentTimeLineList(ctx, &service.PaymentTimelineListInternalReq{
		MerchantId: _interface.GetMerchantId(ctx),
		UserId:     req.UserId,
		SortField:  req.SortField,
		SortType:   req.SortType,
		Page:       req.Page,
		Count:      req.Count,
	})
	if err != nil {
		return nil, err
	}
	return &payment.TimeLineListRes{PaymentTimeLines: result.PaymentTimelines}, nil
}
