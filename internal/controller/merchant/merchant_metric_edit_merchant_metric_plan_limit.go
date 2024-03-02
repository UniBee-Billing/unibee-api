package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	_interface "unibee/internal/interface"
	metric2 "unibee/internal/logic/metric"
	"unibee/internal/query"

	"unibee/api/merchant/metric"
)

func (c *ControllerMetric) EditMerchantMetricPlanLimit(ctx context.Context, req *metric.EditMerchantMetricPlanLimitReq) (res *metric.EditMerchantMetricPlanLimitRes, err error) {
	one := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	me, err := metric2.EditMerchantMetricPlanLimit(ctx, &metric2.MerchantMetricPlanLimitInternalReq{
		MerchantId:        _interface.GetMerchantId(ctx),
		MetricPlanLimitId: uint64(req.MetricPlanLimitId),
		MetricLimit:       req.MetricLimit,
	})
	if err != nil {
		return nil, err
	}
	return &metric.EditMerchantMetricPlanLimitRes{MerchantMetricPlanLimit: me}, nil
}
