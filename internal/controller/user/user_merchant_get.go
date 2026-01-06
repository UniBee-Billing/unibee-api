package user

import (
	"context"
	"unibee/api/bean"
	"unibee/api/user/merchant"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

func (c *ControllerMerchant) Get(ctx context.Context, req *merchant.GetReq) (res *merchant.GetRes, err error) {
	var list []*entity.CreditConfig
	q := dao.CreditConfig.Ctx(ctx).
		Where(dao.CreditConfig.Columns().MerchantId, _interface.GetMerchantId(ctx))
	_ = q.Scan(&list)
	var promoCreditList = make([]*bean.CreditConfig, 0)
	var creditList = make([]*bean.CreditConfig, 0)
	for _, v := range list {
		if v.Type == consts.CreditAccountTypeMain {
			creditList = append(creditList, bean.SimplifyCreditConfig(ctx, v))
		} else if v.Type == consts.CreditAccountTypePromo {
			promoCreditList = append(promoCreditList, bean.SimplifyCreditConfig(ctx, v))
		}
	}
	return &merchant.GetRes{
		Merchant:           bean.SimplifyMerchant(query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))),
		PromoCreditConfigs: promoCreditList,
		CreditConfigs:      creditList,
	}, nil
}
