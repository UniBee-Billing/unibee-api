package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	_interface "unibee/internal/interface"
	metric2 "unibee/internal/logic/metric"
	"unibee/internal/query"

	"unibee/api/merchant/metric"
)

func (c *ControllerMetric) NewPlanLimit(ctx context.Context, req *metric.NewPlanLimitReq) (res *metric.NewPlanLimitRes, err error) {
	one := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
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
	return &metric.NewPlanLimitRes{MerchantMetricPlanLimit: me}, nil
}
