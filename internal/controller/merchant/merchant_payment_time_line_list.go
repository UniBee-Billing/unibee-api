package merchant

import (
	"context"
	"unibee/api/merchant/payment"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/payment/service"
)

func (c *ControllerPayment) TimeLineList(ctx context.Context, req *payment.TimeLineListReq) (res *payment.TimeLineListRes, err error) {
	result, err := service.PaymentTimeLineList(ctx, &service.PaymentTimelineListInternalReq{
		MerchantId:      _interface.GetMerchantId(ctx),
		UserId:          req.UserId,
		AmountStart:     req.AmountStart,
		AmountEnd:       req.AmountEnd,
		Status:          req.Status,
		TimelineTypes:   req.TimelineTypes,
		GatewayIds:      req.GatewayIds,
		Currency:        req.Currency,
		SortField:       req.SortField,
		SortType:        req.SortType,
		Page:            req.Page,
		Count:           req.Count,
		CreateTimeStart: req.CreateTimeStart,
		CreateTimeEnd:   req.CreateTimeEnd,
	})
	if err != nil {
		return nil, err
	}
	return &payment.TimeLineListRes{PaymentTimeLines: result.PaymentTimelines, Total: result.Total}, nil
}
