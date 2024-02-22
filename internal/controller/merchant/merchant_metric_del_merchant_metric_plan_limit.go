package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	metric2 "unibee-api/internal/logic/metric"
	"unibee-api/internal/query"
	"unibee-api/utility"

	"github.com/gogf/gf/v2/errors/gerror"

	"unibee-api/api/merchant/metric"
)

func (c *ControllerMetric) DelMerchantMetricPlanLimit(ctx context.Context, req *metric.DelMerchantMetricPlanLimitReq) (res *metric.DelMerchantMetricPlanLimitRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	one := query.GetMerchantInfoById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	err = metric2.DeleteMerchantMetricPlanLimit(ctx, _interface.GetMerchantId(ctx), req.MetricPlanLimitId)
	if err != nil {
		return nil, err
	}
	return &metric.DelMerchantMetricPlanLimitRes{}, nil
}
