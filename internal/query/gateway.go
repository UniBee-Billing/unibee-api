package query

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetGatewayByGatewayName(ctx context.Context, merchantId uint64, gatewayName string) (one *entity.MerchantGateway) {
	if len(gatewayName) == 0 {
		return nil
	}
	err := dao.MerchantGateway.Ctx(ctx).
		Where(dao.MerchantGateway.Columns().MerchantId, merchantId).
		Where(dao.MerchantGateway.Columns().GatewayName, gatewayName).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetGatewayById(ctx context.Context, id uint64) (one *entity.MerchantGateway) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantGateway.Ctx(ctx).
		Where(dao.MerchantGateway.Columns().Id, id).
		Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetMerchantGatewayList(ctx context.Context, merchantId uint64) (list []*entity.MerchantGateway) {
	var data []*entity.MerchantGateway
	err := dao.MerchantGateway.Ctx(ctx).
		Where(dao.MerchantGateway.Columns().MerchantId, merchantId).
		Where(dao.MerchantGateway.Columns().IsDeleted, 0).
		OrderAsc(dao.MerchantGateway.Columns().GatewayType).
		Scan(&data)
	if err != nil {
		g.Log().Errorf(ctx, "GetMerchantGatewayList error:%s", err)
		return nil
	}
	return data
}

func UpdateGatewayWebhookSecret(ctx context.Context, id uint64, secret string) error {
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
	return nil
}
