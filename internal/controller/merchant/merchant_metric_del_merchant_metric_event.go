package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee-api/api/merchant/metric"
)

func (c *ControllerMetric) DelMerchantMetricEvent(ctx context.Context, req *metric.DelMerchantMetricEventReq) (res *metric.DelMerchantMetricEventRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
