package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/subscription/service"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) UserPendingCryptoSubscriptionDetail(ctx context.Context, req *subscription.UserPendingCryptoSubscriptionDetailReq) (res *subscription.UserPendingCryptoSubscriptionDetailRes, err error) {
	var user *entity.UserAccount
	if _interface.Context().Get(ctx).IsOpenApiCall {
		if req.UserId == 0 {
			utility.Assert(len(req.ExternalUserId) > 0, "ExternalUserId|UserId is nil, one of it is required")
			user = query.GetUserAccountByExternalUserId(ctx, _interface.GetMerchantId(ctx), req.ExternalUserId)
			utility.AssertError(err, "Server Error")
		} else {
			user = query.GetUserAccountById(ctx, req.UserId)
		}
	} else {
		user = query.GetUserAccountById(ctx, req.UserId)
	}
	utility.Assert(user != nil, "user not found")
	one := query.GetLatestCreateOrProcessingSubscriptionByUserId(ctx, user.Id, _interface.GetMerchantId(ctx), req.ProductId)
	if one != nil {
		gateway := query.GetGatewayById(ctx, one.GatewayId)
		if gateway.GatewayType == consts.GatewayTypeCrypto {
			detail, err := service.SubscriptionDetail(ctx, one.SubscriptionId)
			if err == nil {
				return &subscription.UserPendingCryptoSubscriptionDetailRes{
					Subscription: detail,
				}, nil
			}
		}
	}
	return &subscription.UserPendingCryptoSubscriptionDetailRes{}, nil
}
