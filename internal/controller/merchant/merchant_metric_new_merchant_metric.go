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

func (c *ControllerMetric) NewMerchantMetric(ctx context.Context, req *metric.NewMerchantMetricReq) (res *metric.NewMerchantMetricRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	one := query.GetMerchantInfoById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	me, err := metric2.NewMerchantMetric(ctx, &metric2.NewMerchantMetricInternalReq{
		MerchantId:          _interface.GetMerchantId(ctx),
		Code:                req.Code,
		Name:                req.MetricName,
		Description:         req.MetricDescription,
		AggregationType:     req.AggregationType,
		AggregationProperty: req.AggregationProperty,
	})
	if err != nil {
		return nil, err
	}
	return &metric.NewMerchantMetricRes{MerchantMetric: me}, nil
}
