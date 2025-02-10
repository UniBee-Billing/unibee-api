package query

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consts"
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

func GetMerchantGatewayList(ctx context.Context, merchantId uint64, archive *bool) (list []*entity.MerchantGateway) {
	var data []*entity.MerchantGateway
	q := dao.MerchantGateway.Ctx(ctx).
		Where(dao.MerchantGateway.Columns().MerchantId, merchantId)
	if archive != nil && *archive {
		q = q.WhereGT(dao.MerchantGateway.Columns().IsDeleted, 0)
	} else if archive != nil && !*archive {
		q = q.Where(dao.MerchantGateway.Columns().IsDeleted, 0)
	}
	err := q.Order("is_deleted asc, enum_key desc").
		Scan(&data)
	if err != nil {
		g.Log().Errorf(ctx, "GetMerchantGatewayList error:%s", err)
		return nil
	}
	var validGateways []*entity.MerchantGateway
	for _, v := range data {
		if v.GatewayType == consts.GatewayTypeWireTransfer {
			validGateways = append(validGateways, v)
		} else if len(v.GatewayKey) > 0 {
			validGateways = append(validGateways, v)
		}
	}
	return validGateways
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
