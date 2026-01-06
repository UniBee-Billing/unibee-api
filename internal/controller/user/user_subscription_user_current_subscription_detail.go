package user

import (
	"context"
	"unibee/api/bean"
	"unibee/internal/consts"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/subscription/service/detail"
	"unibee/internal/query"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) UserCurrentSubscriptionDetail(ctx context.Context, req *subscription.UserCurrentSubscriptionDetailReq) (res *subscription.UserCurrentSubscriptionDetailRes, err error) {
	user := query.GetUserAccountById(ctx, _interface.Context().Get(ctx).User.Id)
	one := query.GetLatestActiveOrIncompleteOrCreateSubscriptionByUserId(ctx, user.Id, _interface.GetMerchantId(ctx), req.ProductId)
	if one != nil {
		subscriptionDetail, err := detail.SubscriptionDetail(ctx, one.SubscriptionId)
		if err != nil {
			return nil, err
		}
		if subscriptionDetail != nil {
			return &subscription.UserCurrentSubscriptionDetailRes{
				User:                                subscriptionDetail.User,
				PromoCreditAccounts:                 bean.SimplifyCreditAccountList(ctx, query.GetCreditAccountListByUserId(ctx, user.Id, consts.CreditAccountTypePromo)),
				CreditAccounts:                      bean.SimplifyCreditAccountList(ctx, query.GetCreditAccountListByUserId(ctx, user.Id, consts.CreditAccountTypeMain)),
				Subscription:                        subscriptionDetail.Subscription,
				Plan:                                subscriptionDetail.Plan,
				Gateway:                             subscriptionDetail.Gateway,
				AddonParams:                         subscriptionDetail.AddonParams,
				Addons:                              subscriptionDetail.Addons,
				LatestInvoice:                       subscriptionDetail.LatestInvoice,
				UnfinishedSubscriptionPendingUpdate: subscriptionDetail.UnfinishedSubscriptionPendingUpdate,
			}, nil
		} else {
			return nil, nil
		}
	}
	return nil, nil
}
