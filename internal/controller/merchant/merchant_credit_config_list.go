package merchant

import (
	"context"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	entity "unibee/internal/model/entity/default"

	"unibee/api/merchant/credit"
)

func (c *ControllerCredit) ConfigList(ctx context.Context, req *credit.ConfigListReq) (res *credit.ConfigListRes, err error) {
	var list []*entity.CreditConfig
	query := dao.CreditConfig.Ctx(ctx).
		Where(dao.CreditConfig.Columns().MerchantId, _interface.GetMerchantId(ctx))
	if req.Types != nil && len(req.Types) > 0 {
		query = query.WhereIn(dao.CreditConfig.Columns().Type, req.Types)
	}
	if len(req.Currency) > 0 {
		query = query.Where(dao.CreditConfig.Columns().Currency, req.Currency)
	}
	_ = query.Scan(&list)
	var resultList = make([]*bean.CreditConfig, 0)
	for _, v := range list {
		resultList = append(resultList, bean.SimplifyCreditConfig(ctx, v))
	}
	return &credit.ConfigListRes{CreditConfigs: resultList}, nil
}
