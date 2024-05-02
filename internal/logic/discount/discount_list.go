package discount

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

func MerchantDiscountCodeList(ctx context.Context, merchantId uint64) []*bean.MerchantDiscountCodeSimplify {
	var mainList = make([]*bean.MerchantDiscountCodeSimplify, 0)
	var list []*entity.MerchantDiscountCode
	err := dao.MerchantDiscountCode.Ctx(ctx).
		Where(dao.MerchantDiscountCode.Columns().MerchantId, merchantId).
		Where(dao.MerchantDiscountCode.Columns().Type, 0).
		Where(dao.MerchantDiscountCode.Columns().IsDeleted, 0).
		OrderDesc(dao.MerchantDiscountCode.Columns().GmtCreate).
		Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "MerchantDiscountCodeList err:%s", err.Error())
		return mainList
	}
	for _, one := range list {
		mainList = append(mainList, bean.SimplifyMerchantDiscountCode(one))
	}

	return mainList
}
