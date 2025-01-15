package api

import (
	"context"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/gateway/util"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func GetGatewayUser(ctx context.Context, userId uint64, gatewayId uint64) (one *entity.GatewayUser) {
	utility.Assert(userId > 0, "invalid userId")
	utility.Assert(gatewayId > 0, "invalid gatewayId")
	err := dao.GatewayUser.Ctx(ctx).
		Where(dao.GatewayUser.Columns().UserId, userId).
		Where(dao.GatewayUser.Columns().GatewayId, gatewayId).
		Where(dao.GatewayUser.Columns().IsDeleted, 0).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetUserAccountById(ctx context.Context, id uint64) (one *entity.UserAccount) {
	if id <= 0 {
		return nil
	}
	err := dao.UserAccount.Ctx(ctx).
		Where(dao.UserAccount.Columns().Id, id).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func QueryAndCreateGatewayUser(ctx context.Context, gateway *entity.MerchantGateway, userId uint64) *entity.GatewayUser {
	gatewayUser := GetGatewayUser(ctx, userId, gateway.Id)
	if gatewayUser == nil {
		user := GetUserAccountById(ctx, userId)
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
			detailQuery, err := GetGatewayServiceProvider(ctx, gateway.Id).GatewayUserDetailQuery(ctx, gateway, gatewayUser.UserId)
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
		user := GetUserAccountById(ctx, userId)
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
