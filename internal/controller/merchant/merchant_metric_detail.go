package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/api/merchant/metric"
	_interface "unibee/internal/interface"
	metric2 "unibee/internal/logic/metric"
	"unibee/internal/query"
)

func (c *ControllerMetric) Detail(ctx context.Context, req *metric.DetailReq) (res *metric.DetailRes, err error) {
	one := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	return &metric.DetailRes{MerchantMetric: metric2.MerchantMetricDetail(ctx, _interface.GetMerchantId(ctx), req.MetricId)}, nil
}
