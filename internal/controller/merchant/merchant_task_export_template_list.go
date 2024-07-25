package merchant

import (
	"context"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"

	"unibee/api/merchant/task"
)

func (c *ControllerTask) ExportTemplateList(ctx context.Context, req *task.ExportTemplateListReq) (res *task.ExportTemplateListRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "No Permission")
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}
	var entities []*entity.MerchantBatchExportTemplate
	var total = 0
	var sortKey = "id desc"
	q := dao.MerchantBatchExportTemplate.Ctx(ctx).
		Where(dao.MerchantBatchExportTemplate.Columns().MerchantId, _interface.GetMerchantId(ctx)).
		Where(dao.MerchantBatchExportTemplate.Columns().MemberId, _interface.Context().Get(ctx).MerchantMember.Id).
		Where(dao.MerchantBatchExportTemplate.Columns().IsDeleted, 0).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count)
	if len(req.Task) > 0 {
		q = q.Where(dao.MerchantBatchExportTemplate.Columns().Task, req.Task)
	}
	err = q.ScanAndCount(&entities, &total, true)
	var list []*bean.MerchantBatchExportTemplateSimplify
	if err == nil && len(entities) > 0 {
		for _, one := range entities {
			list = append(list, bean.SimplifyMerchantBatchExportTemplate(one))
		}
	}
	return &task.ExportTemplateListRes{
		Templates: list,
		Total:     total,
	}, nil
}
