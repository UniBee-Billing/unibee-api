package merchant

import (
	"context"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/metric_event"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"

	"unibee-api/api/merchant/metric"
)

func (c *ControllerMetric) UserMerchantMetricStat(ctx context.Context, req *metric.UserMerchantMetricStatReq) (res *metric.UserMerchantMetricStatRes, err error) {
	utility.Assert(req.UserId > 0 || len(req.ExternalUserId) > 0, "UserId or ExternalUserId Needed")
	var user *entity.UserAccount
	if req.UserId > 0 {
		user = query.GetUserAccountById(ctx, uint64(req.UserId))
	} else if len(req.ExternalUserId) > 0 {
		user = query.GetUserAccountByExternalUserId(ctx, _interface.GetMerchantId(ctx), req.ExternalUserId)
	}
	utility.Assert(user != nil, "user not found")
	list := metric_event.GetUserMetricLimitStat(ctx, _interface.GetMerchantId(ctx), user)
	return &metric.UserMerchantMetricStatRes{UserMerchantMetricStats: list}, nil
}
