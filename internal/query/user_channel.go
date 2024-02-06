package query

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/utility"
)

func GetGatewayUser(ctx context.Context, userId int64, gatewayId int64) (one *entity.GatewayUser) {
	utility.Assert(userId > 0, "invalid userId")
	utility.Assert(gatewayId > 0, "invalid gatewayId")
	err := dao.GatewayUser.Ctx(ctx).Where(entity.GatewayUser{UserId: userId, GatewayId: gatewayId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetGatewayUserByGatewayUserId(ctx context.Context, gatewayUserId string, gatewayId int64) (one *entity.GatewayUser) {
	utility.Assert(len(gatewayUserId) > 0, "invalid gatewayUserId")
	utility.Assert(gatewayId > 0, "invalid gatewayId")
	err := dao.GatewayUser.Ctx(ctx).Where(entity.GatewayUser{GatewayUserId: gatewayUserId, GatewayId: gatewayId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func CreateOrUpdateGatewayUser(ctx context.Context, userId int64, gatewayId int64, gatewayUserId string, gatewayDefaultPaymentMethod string) (*entity.GatewayUser, error) {
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
			CreateAt:                    gtime.Now().Timestamp(),
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
	} else {
		one.GatewayDefaultPaymentMethod = gatewayDefaultPaymentMethod
		_, err := dao.GatewayUser.Ctx(ctx).Data(g.Map{
			dao.GatewayUser.Columns().GatewayDefaultPaymentMethod: gatewayDefaultPaymentMethod,
		}).Where(dao.GatewayUser.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateGatewayUser record insert failure %s`, err)
			return nil, err
		}
	}
	return one, nil
}
