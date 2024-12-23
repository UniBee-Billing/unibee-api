package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	_interface "unibee/internal/interface/context"
	metric2 "unibee/internal/logic/metric"
	"unibee/internal/query"

	"unibee/api/merchant/metric"
)

func (c *ControllerMetric) DeletePlanLimit(ctx context.Context, req *metric.DeletePlanLimitReq) (res *metric.DeletePlanLimitRes, err error) {
	one := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	err = metric2.DeleteMerchantMetricPlanLimit(ctx, _interface.GetMerchantId(ctx), req.MetricPlanLimitId)
	if err != nil {
		return nil, err
	}
	return &metric.DeletePlanLimitRes{}, nil
}
