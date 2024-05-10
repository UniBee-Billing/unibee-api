package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	_interface "unibee/internal/interface"
	metric2 "unibee/internal/logic/metric"
	"unibee/internal/query"

	"unibee/api/merchant/metric"
)

func (c *ControllerMetric) List(ctx context.Context, req *metric.ListReq) (res *metric.ListRes, err error) {
	one := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	list, total := metric2.MerchantMetricList(ctx, _interface.GetMerchantId(ctx))
	return &metric.ListRes{MerchantMetrics: list, Total: total}, nil
}
