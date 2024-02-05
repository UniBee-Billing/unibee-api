package api

import (
	"context"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func queryAndCreateChannelUser(ctx context.Context, gateway *entity.MerchantGateway, userId int64) *entity.GatewayUser {
	channelUser := query.GetGatewayUser(ctx, userId, int64(gateway.Id))
	if channelUser == nil {
		user := query.GetUserAccountById(ctx, uint64(userId))
		utility.Assert(user != nil, "user not found")
		utility.Assert(len(user.Email) > 0, "invalid user email")
		create, err := GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewayUserCreate(ctx, gateway, user)
		utility.AssertError(err, "GatewayUserCreate")
		channelUser, err = query.CreateOrUpdateGatewayUser(ctx, userId, int64(gateway.Id), create.GatewayUserId, "")
		utility.AssertError(err, "CreateOrUpdateGatewayUser")
		return channelUser
	} else {
		if len(channelUser.GatewayDefaultPaymentMethod) == 0 {
			//no default payment method, query it
			detailQuery, err := GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewayUserDetailQuery(ctx, gateway, channelUser.UserId)
			utility.AssertError(err, "GatewayUserDetailQuery")
			if len(detailQuery.DefaultPaymentMethod) > 0 {
				channelUser, err = query.CreateOrUpdateGatewayUser(ctx, userId, int64(gateway.Id), channelUser.GatewayUserId, detailQuery.DefaultPaymentMethod)
				channelUser.GatewayDefaultPaymentMethod = detailQuery.DefaultPaymentMethod
			}
		}
		return channelUser
	}
}

func queryAndCreateChannelUserWithOutPaymentMethod(ctx context.Context, gateway *entity.MerchantGateway, userId int64) *entity.GatewayUser {
	channelUser := query.GetGatewayUser(ctx, userId, int64(gateway.Id))
	if channelUser == nil {
		user := query.GetUserAccountById(ctx, uint64(userId))
		utility.Assert(user != nil, "user not found")
		utility.Assert(len(user.Email) > 0, "invalid user email")
		create, err := GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewayUserCreate(ctx, gateway, user)
		utility.AssertError(err, "GatewayUserCreate")
		channelUser, err = query.CreateOrUpdateGatewayUser(ctx, userId, int64(gateway.Id), create.GatewayUserId, "")
		utility.AssertError(err, "CreateOrUpdateGatewayUser")
		return channelUser
	} else {
		return channelUser
	}
}
