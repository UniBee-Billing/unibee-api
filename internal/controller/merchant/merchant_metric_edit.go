package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	metric2 "unibee/internal/logic/metric"
	"unibee/internal/query"

	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/metric"
)

func (c *ControllerMetric) Edit(ctx context.Context, req *metric.EditReq) (res *metric.EditRes, err error) {
	one := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	me, err := metric2.EditMerchantMetric(ctx, _interface.GetMerchantId(ctx), req.MetricId, req.Type, req.MetricName, req.MetricDescription, req.Unit, req.MetaData)
	if err != nil {
		return nil, err
	}
	return &metric.EditRes{MerchantMetric: me}, nil
}
