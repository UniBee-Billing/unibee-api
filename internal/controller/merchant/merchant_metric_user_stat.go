package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/metric_event"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/metric"
)

func (c *ControllerMetric) UserStat(ctx context.Context, req *metric.UserStatReq) (res *metric.UserStatRes, err error) {
	utility.Assert(req.UserId > 0 || len(req.ExternalUserId) > 0, "UserId or ExternalUserId Needed")
	var user *entity.UserAccount
	if req.UserId > 0 {
		user = query.GetUserAccountById(ctx, uint64(req.UserId))
	} else if len(req.ExternalUserId) > 0 {
		user = query.GetUserAccountByExternalUserId(ctx, _interface.GetMerchantId(ctx), req.ExternalUserId)
	}
	utility.Assert(user != nil, "user not found")
	list := metric_event.GetUserMetricLimitStat(ctx, _interface.GetMerchantId(ctx), user)
	return &metric.UserStatRes{UserMerchantMetricStats: list}, nil
}
