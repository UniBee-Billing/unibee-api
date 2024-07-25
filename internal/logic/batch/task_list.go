package batch

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func MerchantBatchTaskList(ctx context.Context, merchantId uint64, memberId uint64, page int, count int) ([]*bean.MerchantBatchTaskSimplify, int) {
	if count <= 0 {
		count = 20
	}
	if page < 0 {
		page = 0
	}
	var total = 0
	var resultList = make([]*bean.MerchantBatchTaskSimplify, 0)
	var mainList = make([]*entity.MerchantBatchTask, 0)
	err := dao.MerchantBatchTask.Ctx(ctx).
		Where(dao.MerchantBatchTask.Columns().MerchantId, merchantId).
		Where(dao.MerchantBatchTask.Columns().MemberId, memberId).
		OrderDesc(dao.MerchantBatchTask.Columns().CreateTime).
		Limit(page*count, count).
		ScanAndCount(&mainList, &total, true)
	if err != nil {
		g.Log().Errorf(ctx, "MerchantMemberList err:%s", err.Error())
		return resultList, len(resultList)
	}
	for _, one := range mainList {
		resultList = append(resultList, bean.SimplifyMerchantBatchTask(one))
	}
	return resultList, total
}
