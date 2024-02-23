package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	_interface "unibee-api/internal/interface"
	metric2 "unibee-api/internal/logic/metric"
	"unibee-api/internal/query"

	"unibee-api/api/merchant/metric"
)

func (c *ControllerMetric) MerchantMetricList(ctx context.Context, req *metric.MerchantMetricListReq) (res *metric.MerchantMetricListRes, err error) {
	one := query.GetMerchantInfoById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	return &metric.MerchantMetricListRes{MerchantMetrics: metric2.MerchantMetricList(ctx, _interface.GetMerchantId(ctx))}, nil
}
