package member

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

func MerchantMemberList(ctx context.Context, merchantId uint64) ([]*bean.MerchantMemberSimplify, int) {
	var mainList = make([]*entity.MerchantMember, 0)
	var resultList = make([]*bean.MerchantMemberSimplify, 0)
	err := dao.MerchantMember.Ctx(ctx).
		Where(dao.MerchantMember.Columns().MerchantId, merchantId).
		Where(dao.MerchantMember.Columns().IsDeleted, 0).
		Scan(&mainList)
	for _, one := range mainList {
		resultList = append(resultList, bean.SimplifyMerchantMember(one))
	}
	if err != nil {
		g.Log().Errorf(ctx, "MerchantMemberList err:%s", err.Error())
		return resultList, len(mainList)
	}
	return resultList, len(mainList)
}
