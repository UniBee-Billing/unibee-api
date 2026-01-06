package user

import (
	"context"
	"unibee/api/bean"
	"unibee/api/user/subscription"
	"unibee/internal/consts"
	"unibee/internal/logic/subscription/service/detail"
	"unibee/internal/query"
)

func (c *ControllerSubscription) Detail(ctx context.Context, req *subscription.DetailReq) (res *subscription.DetailRes, err error) {
	subscriptionDetail, err := detail.SubscriptionDetail(ctx, req.SubscriptionId)
	if err != nil {
		return nil, err
	}
	return &subscription.DetailRes{
		User:                                subscriptionDetail.User,
		PromoCreditAccounts:                 bean.SimplifyCreditAccountList(ctx, query.GetCreditAccountListByUserId(ctx, subscriptionDetail.User.Id, consts.CreditAccountTypePromo)),
		CreditAccounts:                      bean.SimplifyCreditAccountList(ctx, query.GetCreditAccountListByUserId(ctx, subscriptionDetail.User.Id, consts.CreditAccountTypeMain)),
		Subscription:                        subscriptionDetail.Subscription,
		Plan:                                subscriptionDetail.Plan,
		Gateway:                             subscriptionDetail.Gateway,
		AddonParams:                         subscriptionDetail.AddonParams,
		Addons:                              subscriptionDetail.Addons,
		LatestInvoice:                       subscriptionDetail.LatestInvoice,
		UnfinishedSubscriptionPendingUpdate: subscriptionDetail.UnfinishedSubscriptionPendingUpdate,
	}, nil
}
