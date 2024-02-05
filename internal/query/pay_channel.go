package query

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetGatewayById(ctx context.Context, id int64) (one *entity.MerchantGateway) {
	if id <= 0 {
		return nil
	}
	m := dao.MerchantGateway.Ctx(ctx)
	err := m.Where(entity.MerchantGateway{Id: uint64(id)}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetGatewayByGatewayName(ctx context.Context, gatewayName string) (one *entity.MerchantGateway) {
	if len(gatewayName) == 0 {
		return nil
	}
	err := dao.MerchantGateway.Ctx(ctx).Where(entity.MerchantGateway{GatewayName: gatewayName}).OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetGatewaysGroupByEnumKey(ctx context.Context) []*entity.MerchantGateway {
	var data []*entity.MerchantGateway
	err := dao.MerchantGateway.Ctx(ctx).Group(dao.MerchantGateway.Columns().EnumKey).
		OmitEmpty().Scan(&data)
	if err != nil {
		g.Log().Errorf(ctx, "GetGatewaysGroupByEnumKey error:%s", err)
		return nil
	}
	return data
}

func GetPaymentTypeGatewayById(ctx context.Context, id int64) (one *entity.MerchantGateway) {
	if id <= 0 {
		return nil
	}
	m := dao.MerchantGateway.Ctx(ctx)
	err := m.Where(entity.MerchantGateway{Id: uint64(id)}).
		Where(m.Builder().
			Where(entity.MerchantGateway{GatewayType: consts.GatewayTypeOneTimePayment}).WhereOr("gateway_type is null")).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetSubscriptionTypeGatewayById(ctx context.Context, id int64) (one *entity.MerchantGateway) {
	if id <= 0 {
		return nil
	}
	m := dao.MerchantGateway.Ctx(ctx)
	err := m.Where(entity.MerchantGateway{Id: uint64(id), GatewayType: consts.GatewayTypeSubscription}).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetListSubscriptionTypeGateways(ctx context.Context) (list []*entity.MerchantGateway) {
	var data []*entity.MerchantGateway
	err := dao.MerchantGateway.Ctx(ctx).Where(entity.MerchantGateway{GatewayType: consts.GatewayTypeSubscription}).
		OmitEmpty().Scan(&data)
	if err != nil {
		g.Log().Errorf(ctx, "GetListSubscriptionTypeGateways error:%s", err)
		return nil
	}
	return data
}

func SaveGatewayUniqueProductId(ctx context.Context, id int64, productId string) error {
	if len(productId) == 0 || id < 0 {
		return nil
	}
	_, err := dao.MerchantGateway.Ctx(ctx).Data(g.Map{
		dao.MerchantGateway.Columns().UniqueProductId: productId,
		dao.MerchantGateway.Columns().GmtModify:       gtime.Now(),
	}).Where(dao.MerchantGateway.Columns().Id, id).Update()
	if err != nil {
		return err
	}
	return nil
}

func UpdateGatewayWebhookSecret(ctx context.Context, id int64, secret string) error {
	if id <= 0 {
		return gerror.New("invalid id")
	}
	_, err := dao.MerchantGateway.Ctx(ctx).Data(g.Map{
		dao.MerchantGateway.Columns().WebhookSecret: secret,
		dao.MerchantGateway.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.MerchantGateway.Columns().Id, id).Update()
	if err != nil {
		return err
	}
	//rowAffected, err := update.RowsAffected()
	//if rowAffected != 1 {
	//	return gerror.Newf("UpdateGatewayWebhookSecret update err:%s", update)
	//}
	return nil
}
