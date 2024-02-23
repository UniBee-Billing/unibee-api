package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	_interface "unibee-api/internal/interface"
	metric2 "unibee-api/internal/logic/metric"
	"unibee-api/internal/query"

	"unibee-api/api/merchant/metric"
)

func (c *ControllerMetric) NewMerchantMetricPlanLimit(ctx context.Context, req *metric.NewMerchantMetricPlanLimitReq) (res *metric.NewMerchantMetricPlanLimitRes, err error) {
	one := query.GetMerchantInfoById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	me, err := metric2.NewMerchantMetricPlanLimit(ctx, &metric2.MerchantMetricPlanLimitInternalReq{
		MerchantId:  _interface.GetMerchantId(ctx),
		PlanId:      req.PlanId,
		MetricId:    req.MetricId,
		MetricLimit: req.MetricLimit,
	})
	if err != nil {
		return nil, err
	}
	return &metric.NewMerchantMetricPlanLimitRes{MerchantMetricPlanLimit: me}, nil
}
