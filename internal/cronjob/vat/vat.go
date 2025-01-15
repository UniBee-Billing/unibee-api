package vat

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/vat_gateway/setup"
	entity "unibee/internal/model/entity/default"
)

func TaskForSyncVatData(ctx context.Context) {
	var list []*entity.Merchant
	_ = dao.Merchant.Ctx(ctx).
		Where(dao.Merchant.Columns().IsDeleted, 0).
		Scan(&list)
	for _, v := range list {
		err := setup.InitMerchantDefaultVatGateway(ctx, v.Id)
		if err != nil {
			g.Log().Errorf(ctx, "TaskForSyncVatData merchantId:%d err:%s", v.Id, err.Error())
		} else {
			g.Log().Infof(ctx, "TaskForSyncVatData success merchantId:%d", v.Id)
		}
	}

}
