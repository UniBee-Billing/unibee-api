package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	_interface "unibee/internal/interface"
	metric2 "unibee/internal/logic/metric"
	"unibee/internal/query"

	"unibee/api/merchant/metric"
)

func (c *ControllerMetric) EditMerchantMetric(ctx context.Context, req *metric.EditMerchantMetricReq) (res *metric.EditMerchantMetricRes, err error) {
	one := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	me, err := metric2.EditMerchantMetric(ctx, _interface.GetMerchantId(ctx), req.MetricId, req.MetricName, req.MetricDescription)
	if err != nil {
		return nil, err
	}
	return &metric.EditMerchantMetricRes{MerchantMetric: me}, nil
}
