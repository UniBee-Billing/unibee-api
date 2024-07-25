package query

import (
	"context"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetMerchantTaskExportTemplateById(ctx context.Context, id int64) (one *entity.MerchantBatchExportTemplate) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantBatchExportTemplate.Ctx(ctx).
		Where(dao.MerchantBatchExportTemplate.Columns().Id, id).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}
