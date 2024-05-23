package user

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) UserPendingCryptoSubscriptionDetail(ctx context.Context, req *subscription.UserPendingCryptoSubscriptionDetailReq) (res *subscription.UserPendingCryptoSubscriptionDetailRes, err error) {
	user := query.GetUserAccountById(ctx, _interface.Context().Get(ctx).User.Id)
	utility.Assert(user != nil, "user not found")
	one := query.GetLatestCreateOrProcessingSubscriptionByUserId(ctx, user.Id, _interface.GetMerchantId(ctx))
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
