package role

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

func MerchantRoleList(ctx context.Context, merchantId uint64) []*bean.MerchantRoleSimplify {
	var mainList = make([]*bean.MerchantRoleSimplify, 0)
	var list []*entity.MerchantRole
	err := dao.MerchantRole.Ctx(ctx).
		Where(dao.MerchantRole.Columns().MerchantId, merchantId).
		Where(dao.MerchantRole.Columns().IsDeleted, 0).
		Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "MerchantRoleList err:%s", err.Error())
		return mainList
	}
	for _, one := range list {
		mainList = append(mainList, bean.SimplifyMerchantRole(one))
	}

	return mainList
}
