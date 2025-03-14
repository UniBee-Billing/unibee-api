package util

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

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

func GetGatewayUserByGatewayUserId(ctx context.Context, gatewayUserId string, gatewayId uint64) (one *entity.GatewayUser) {
	utility.Assert(len(gatewayUserId) > 0, "invalid gatewayUserId")
	utility.Assert(gatewayId > 0, "invalid gatewayId")
	err := dao.GatewayUser.Ctx(ctx).
		Where(dao.GatewayUser.Columns().GatewayUserId, gatewayUserId).
		Where(dao.GatewayUser.Columns().GatewayId, gatewayId).
		Where(dao.GatewayUser.Columns().IsDeleted, 0).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func CreateOrUpdateGatewayUser(ctx context.Context, userId uint64, gatewayId uint64, gatewayUserId string, gatewayDefaultPaymentMethod string) (*entity.GatewayUser, error) {
	utility.Assert(userId > 0, "invalid userId")
	utility.Assert(gatewayId > 0, "invalid gatewayId")
	utility.Assert(len(gatewayUserId) > 0, "invalid gatewayUserId")
	one := GetGatewayUser(ctx, userId, gatewayId)
	if one == nil {
		one = &entity.GatewayUser{
			UserId:                      userId,
			GatewayId:                   gatewayId,
			GatewayUserId:               gatewayUserId,
			GatewayDefaultPaymentMethod: gatewayDefaultPaymentMethod,
			CreateTime:                  gtime.Now().Timestamp(),
		}
		result, err := dao.GatewayUser.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateGatewayUser record insert failure %s`, err)
			return nil, err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		one.Id = uint64(uint(id))
	} else if len(gatewayDefaultPaymentMethod) > 0 {
		one.GatewayDefaultPaymentMethod = gatewayDefaultPaymentMethod
		_, err := dao.GatewayUser.Ctx(ctx).Data(g.Map{
			dao.GatewayUser.Columns().GatewayDefaultPaymentMethod: gatewayDefaultPaymentMethod,
		}).Where(dao.GatewayUser.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateGatewayUser update gatewayDefaultPaymentMethod failure %s`, err)
			return nil, err
		}
	}
	return one, nil
}

func CreateGatewayUser(ctx context.Context, userId uint64, gatewayId uint64, gatewayUserId string) (*entity.GatewayUser, error) {
	utility.Assert(userId > 0, "invalid userId")
	utility.Assert(gatewayId > 0, "invalid gatewayId")
	utility.Assert(len(gatewayUserId) > 0, "invalid gatewayUserId")
	one := GetGatewayUser(ctx, userId, gatewayId)
	if one == nil {
		one = &entity.GatewayUser{
			UserId:        userId,
			GatewayId:     gatewayId,
			GatewayUserId: gatewayUserId,
			CreateTime:    gtime.Now().Timestamp(),
		}
		result, err := dao.GatewayUser.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateGatewayUser record insert failure %s`, err)
			return nil, err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		one.Id = uint64(uint(id))
		return one, nil
	} else {
		return nil, gerror.New("same gatewayUserId exist")
	}
}
