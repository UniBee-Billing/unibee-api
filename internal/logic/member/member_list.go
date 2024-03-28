package member

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	dao "unibee/internal/dao/oversea_pay"
)

func MerchantMemberList(ctx context.Context, merchantId uint64) []*bean.MerchantMemberSimplify {
	var mainList []*bean.MerchantMemberSimplify
	err := dao.MerchantMember.Ctx(ctx).
		Where(dao.MerchantMember.Columns().MerchantId, merchantId).
		Where(dao.MerchantMember.Columns().IsDeleted, 0).
		Scan(&mainList)
	if err != nil {
		g.Log().Errorf(ctx, "MerchantMemberList err:%s", err.Error())
		return mainList
	}
	return mainList
}
