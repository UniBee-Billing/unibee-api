package api

import (
	"context"
	"unibee/internal/logic/gateway/util"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func QueryAndCreateGatewayUser(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) *entity.GatewayUser {
	gatewayUser := util.GetGatewayUser(ctx, userId, gateway.Id)
	if gatewayUser == nil {
		user := util.GetUserAccountById(ctx, userId)
		utility.Assert(user != nil, "user not found")
		utility.Assert(len(user.Email) > 0, "invalid user email")
		create, err := GetGatewayServiceProvider(ctx, gateway.Id).GatewayUserCreate(ctx, gateway, user)
		utility.AssertError(err, "GatewayUserCreate")
		gatewayUser, err = util.CreateOrUpdateGatewayUser(ctx, userId, gateway.Id, create.GatewayUserId, "")
		utility.AssertError(err, "CreateOrUpdateGatewayUser")
		return gatewayUser
	} else {
		if len(gatewayUser.GatewayDefaultPaymentMethod) == 0 {
			//no default payment method, query
			detailQuery, err := GetGatewayServiceProvider(ctx, gateway.Id).GatewayUserDetailQuery(ctx, gateway, gatewayUser.GatewayUserId)
			utility.AssertError(err, "GatewayUserDetailQuery")
			if len(detailQuery.DefaultPaymentMethod) > 0 {
				gatewayUser, err = util.CreateOrUpdateGatewayUser(ctx, userId, gateway.Id, gatewayUser.GatewayUserId, detailQuery.DefaultPaymentMethod)
				gatewayUser.GatewayDefaultPaymentMethod = detailQuery.DefaultPaymentMethod
			}
		}
		return gatewayUser
	}
}

func QueryAndCreateGatewayUserWithOutPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) *entity.GatewayUser {
	gatewayUser := util.GetGatewayUser(ctx, userId, gateway.Id)
	if gatewayUser == nil {
		user := util.GetUserAccountById(ctx, userId)
		utility.Assert(user != nil, "user not found")
		utility.Assert(len(user.Email) > 0, "invalid user email")
		create, err := GetGatewayServiceProvider(ctx, gateway.Id).GatewayUserCreate(ctx, gateway, user)
		utility.AssertError(err, "GatewayUserCreate")
		gatewayUser, err = util.CreateOrUpdateGatewayUser(ctx, userId, gateway.Id, create.GatewayUserId, "")
		utility.AssertError(err, "CreateOrUpdateGatewayUser")
		return gatewayUser
	} else {
		return gatewayUser
	}
}
