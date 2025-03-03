package merchant

import (
	"context"
	"unibee/api/merchant/metric"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/metric_event"
)

func (c *ControllerMetric) EventList(ctx context.Context, req *metric.EventListReq) (res *metric.EventListRes, err error) {
	result, err := metric_event.EventList(ctx, &metric_event.EventListInternalReq{
		MerchantId:      _interface.GetMerchantId(ctx),
		UserIds:         req.UserIds,
		MetricIds:       req.MetricIds,
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
	return &metric.EventListRes{Events: result.Events, Total: result.Total}, nil
}
