package api

import (
	"context"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func QueryAndCreateChannelUser(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) *entity.GatewayUser {
	gatewayUser := query.GetGatewayUser(ctx, userId, gateway.Id)
	if gatewayUser == nil {
		user := query.GetUserAccountById(ctx, userId)
		utility.Assert(user != nil, "user not found")
		utility.Assert(len(user.Email) > 0, "invalid user email")
		create, err := GetGatewayServiceProvider(ctx, gateway.Id).GatewayUserCreate(ctx, gateway, user)
		utility.AssertError(err, "GatewayUserCreate")
		gatewayUser, err = query.CreateOrUpdateGatewayUser(ctx, userId, gateway.Id, create.GatewayUserId, "")
		utility.AssertError(err, "CreateOrUpdateGatewayUser")
		return gatewayUser
	} else {
		if len(gatewayUser.GatewayDefaultPaymentMethod) == 0 {
			//no default payment method, query it
			detailQuery, err := GetGatewayServiceProvider(ctx, gateway.Id).GatewayUserDetailQuery(ctx, gateway, gatewayUser.UserId)
			utility.AssertError(err, "GatewayUserDetailQuery")
			if len(detailQuery.DefaultPaymentMethod) > 0 {
				gatewayUser, err = query.CreateOrUpdateGatewayUser(ctx, userId, gateway.Id, gatewayUser.GatewayUserId, detailQuery.DefaultPaymentMethod)
				gatewayUser.GatewayDefaultPaymentMethod = detailQuery.DefaultPaymentMethod
			}
		}
		return gatewayUser
	}
}

func QueryAndCreateChannelUserWithOutPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) *entity.GatewayUser {
	gatewayUser := query.GetGatewayUser(ctx, userId, gateway.Id)
	if gatewayUser == nil {
		user := query.GetUserAccountById(ctx, userId)
		utility.Assert(user != nil, "user not found")
		utility.Assert(len(user.Email) > 0, "invalid user email")
		create, err := GetGatewayServiceProvider(ctx, gateway.Id).GatewayUserCreate(ctx, gateway, user)
		utility.AssertError(err, "GatewayUserCreate")
		gatewayUser, err = query.CreateOrUpdateGatewayUser(ctx, userId, gateway.Id, create.GatewayUserId, "")
		utility.AssertError(err, "CreateOrUpdateGatewayUser")
		return gatewayUser
	} else {
		return gatewayUser
	}
}
