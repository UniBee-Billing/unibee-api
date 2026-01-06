package merchant

import (
	"context"
	"fmt"
	"unibee/api/bean"
	"unibee/api/merchant/metric"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/metric_event"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerMetric) NewEvent(ctx context.Context, req *metric.NewEventReq) (res *metric.NewEventRes, err error) {
	var one *entity.UserAccount
	if req.UserId > 0 {
		one = query.GetUserAccountById(ctx, req.UserId)
	} else if len(req.Email) > 0 {
		one = query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email)
	} else if len(req.ExternalUserId) > 0 {
		one = query.GetUserAccountByExternalUserId(ctx, _interface.GetMerchantId(ctx), req.ExternalUserId)
	}
	utility.Assert(one != nil, "user not found, should provides one of three options, UserId, ExternalUserId, or Email")
	//if one == nil {
	//	return nil, gerror.New("user not found, should provides one of three options, UserId, ExternalUserId, or Email")
	//}
	event, err := metric_event.NewMerchantMetricEvent(ctx, &metric_event.MerchantMetricEventInternalReq{
		MerchantId:          _interface.GetMerchantId(ctx),
		MetricCode:          req.MetricCode,
		UserId:              one.Id,
		ExternalEventId:     req.ExternalEventId,
		MetricProperties:    req.MetricProperties,
		ProductId:           req.ProductId,
		AggregationValue:    req.AggregationValue,
		AggregationUniqueId: req.AggregationUniqueId,
	})
	if err != nil {
		utility.AssertError(err, fmt.Sprintf("New metric event error:%s", err.Error()))
	}
	if _interface.Context() != nil && _interface.Context().Get(ctx) != nil && _interface.Context().Get(ctx).MerchantMember != nil {
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     one.MerchantId,
			Target:         fmt.Sprintf("NewMetricEvent(%s-%s-%d)", one.Email, req.MetricCode, event.Id),
			Content:        "MetricTestTool",
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)
	}
	//if err != nil {
	//	return nil, err
	//}
	return &metric.NewEventRes{MerchantMetricEvent: bean.SimplifyMerchantMetricEvent(event)}, nil
}
