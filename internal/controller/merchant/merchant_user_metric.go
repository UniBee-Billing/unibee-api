package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/metric_event"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/metric"
)

func (c *ControllerMetric) UserMetric(ctx context.Context, req *metric.UserMetricReq) (res *metric.UserMetricRes, err error) {
	utility.Assert(req.UserId > 0 || len(req.ExternalUserId) > 0, "UserId or ExternalUserId Needed")
	var user *entity.UserAccount
	if req.UserId > 0 {
		user = query.GetUserAccountById(ctx, uint64(req.UserId))
	} else if len(req.ExternalUserId) > 0 {
		user = query.GetUserAccountByExternalUserId(ctx, _interface.GetMerchantId(ctx), req.ExternalUserId)
	}
	utility.Assert(user != nil, "user not found")
	return &metric.UserMetricRes{UserMetric: metric_event.GetUserMetricStat(ctx, _interface.GetMerchantId(ctx), user)}, nil
}
