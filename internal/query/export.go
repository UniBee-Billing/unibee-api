package query

import (
	"context"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

func GetMerchantBatchTask(ctx context.Context, id uint64) (one *entity.MerchantBatchTask) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantBatchTask.Ctx(ctx).
		Where(dao.MerchantBatchTask.Columns().Id, id).
		Scan(&one)
	if err != nil {
		one = nil
	}
	return
}