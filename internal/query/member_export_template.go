package query

import (
	"context"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
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
