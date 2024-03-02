package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	_interface "unibee/internal/interface"
	metric2 "unibee/internal/logic/metric"
	"unibee/internal/query"

	"unibee/api/merchant/metric"
)

func (c *ControllerMetric) Delete(ctx context.Context, req *metric.DeleteReq) (res *metric.DeleteRes, err error) {
	one := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	err = metric2.DeleteMerchantMetric(ctx, _interface.GetMerchantId(ctx), req.MetricId)
	if err != nil {
		return nil, err
	}
	return &metric.DeleteRes{}, nil
}
