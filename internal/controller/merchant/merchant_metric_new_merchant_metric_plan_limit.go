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

func (c *ControllerMetric) NewMerchantMetricPlanLimit(ctx context.Context, req *metric.NewMerchantMetricPlanLimitReq) (res *metric.NewMerchantMetricPlanLimitRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.MerchantId > 0, "merchantUserId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.MerchantId == uint64(req.MerchantId), "merchantId not match")
	}
	one := query.GetMerchantInfoById(ctx, req.MerchantId)
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	me, err := metric2.NewMerchantMetricPlanLimit(ctx, &metric2.MerchantMetricPlanLimitInternalReq{
		MerchantId:  req.MerchantId,
		PlanId:      req.PlanId,
		MetricId:    req.MetricId,
		MetricLimit: req.MetricLimit,
	})
	if err != nil {
		return nil, err
	}
	return &metric.NewMerchantMetricPlanLimitRes{MerchantMetric: me}, nil
}
