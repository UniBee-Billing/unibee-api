package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	_interface "unibee/internal/interface"
	metric2 "unibee/internal/logic/metric"
	"unibee/internal/query"

	"unibee/api/merchant/metric"
)

func (c *ControllerMetric) DelMerchantMetric(ctx context.Context, req *metric.DelMerchantMetricReq) (res *metric.DelMerchantMetricRes, err error) {
	one := query.GetMerchantInfoById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	err = metric2.DeleteMerchantMetric(ctx, _interface.GetMerchantId(ctx), req.MetricId)
	if err != nil {
		return nil, err
	}
	return &metric.DelMerchantMetricRes{}, nil
}
