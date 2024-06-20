package member

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

func MerchantMemberList(ctx context.Context, merchantId uint64, page int, count int) ([]*detail.MerchantMemberDetail, int) {
	if count <= 0 {
		count = 20
	}
	if page < 0 {
		page = 0
	}
	var resultList = make([]*detail.MerchantMemberDetail, 0)
	var mainList = make([]*entity.MerchantMember, 0)
	err := dao.MerchantMember.Ctx(ctx).
		Where(dao.MerchantMember.Columns().MerchantId, merchantId).
		Where(dao.MerchantMember.Columns().IsDeleted, 0).
		Limit(page*count, count).
		Scan(&mainList)
	if err != nil {
		g.Log().Errorf(ctx, "MerchantMemberList err:%s", err.Error())
		return resultList, len(resultList)
	}
	for _, one := range mainList {
		resultList = append(resultList, detail.ConvertMemberToDetail(ctx, one))
	}
	return resultList, len(resultList)
}

func MerchantMemberTotalList(ctx context.Context, merchantId uint64) ([]*detail.MerchantMemberDetail, int) {
	var resultList = make([]*detail.MerchantMemberDetail, 0)
	var mainList = make([]*entity.MerchantMember, 0)
	err := dao.MerchantMember.Ctx(ctx).
		Where(dao.MerchantMember.Columns().MerchantId, merchantId).
		Where(dao.MerchantMember.Columns().IsDeleted, 0).
		Scan(&mainList)
	if err != nil {
		g.Log().Errorf(ctx, "MerchantMemberList err:%s", err.Error())
		return resultList, len(resultList)
	}
	for _, one := range mainList {
		resultList = append(resultList, detail.ConvertMemberToDetail(ctx, one))
	}
	return resultList, len(resultList)
}
