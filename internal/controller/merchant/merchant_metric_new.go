package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	_interface "unibee/internal/interface/context"
	metric2 "unibee/internal/logic/metric"
	"unibee/internal/query"

	"unibee/api/merchant/metric"
)

func (c *ControllerMetric) New(ctx context.Context, req *metric.NewReq) (res *metric.NewRes, err error) {
	one := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	me, err := metric2.NewMerchantMetric(ctx, &metric2.NewMerchantMetricInternalReq{
		MerchantId:          _interface.GetMerchantId(ctx),
		Code:                req.Code,
		Type:                req.Type,
		Name:                req.MetricName,
		Description:         req.MetricDescription,
		AggregationType:     req.AggregationType,
		AggregationProperty: req.AggregationProperty,
	})
	if err != nil {
		return nil, err
	}
	return &metric.NewRes{MerchantMetric: me}, nil
}
