package merchant

import (
	"context"
	"fmt"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface/context"
	session2 "unibee/internal/logic/session"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/session"
)

func (c *ControllerSession) NewSubUpdatePage(ctx context.Context, req *session.NewSubUpdatePageReq) (res *session.NewSubUpdatePageRes, err error) {
	if req.UserId == 0 {
		utility.Assert(len(req.ExternalUserId) > 0 || len(req.Email) > 0, "ExternalUserId|Email is nil, one of it is required when UserId not specified")
		if len(req.ExternalUserId) > 0 {
			user := query.GetUserAccountByExternalUserId(ctx, _interface.GetMerchantId(ctx), req.ExternalUserId)
			if user != nil {
				req.UserId = user.Id
			}
		} else if len(req.Email) > 0 {
			user := query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email)
			if user != nil {
				req.UserId = user.Id
			}
		}
	}
	utility.Assert(req.UserId > 0, "user not found")
	one := query.GetLatestActiveOrIncompleteOrCreateSubscriptionByUserId(ctx, req.UserId, _interface.GetMerchantId(ctx), req.ProductId)
	if one == nil {
		one = query.GetLatestSubscriptionByUserId(ctx, req.UserId, _interface.GetMerchantId(ctx), req.ProductId)
	}
	userSession, _, err := session2.NewUserSession(ctx, one.MerchantId, req.UserId, req.ReturnUrl, req.CancelUrl)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/hosted/sub_update_page?subscriptionId=%s&session=%s", config.GetConfigInstance().Server.GetServerPath(), one.SubscriptionId, userSession)
	if req.PlanId > 0 {
		url = fmt.Sprintf("%s/hosted/sub_update_page?subscriptionId=%s&session=%s&planId=%d", config.GetConfigInstance().Server.GetServerPath(), one.SubscriptionId, userSession, req.PlanId)
	}
	return &session.NewSubUpdatePageRes{Url: url}, nil
}
