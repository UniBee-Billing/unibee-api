package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/operation_log"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"

	"unibee/api/merchant/task"
)

func (c *ControllerTask) EditTemplate(ctx context.Context, req *task.EditTemplateReq) (res *task.EditTemplateRes, err error) {
	utility.Assert(req.TemplateId > 0, "invalid templateId")
	one := query.GetMerchantTaskExportTemplateById(ctx, req.TemplateId)
	utility.Assert(one != nil, "invalid templateId")
	utility.Assert(one.MerchantId == _interface.GetMerchantId(ctx), "No Permission")
	if _interface.Context().Get(ctx) != nil && _interface.Context().Get(ctx).MerchantMember != nil {
		utility.Assert(one.MemberId == _interface.Context().Get(ctx).MerchantMember.Id, "No Permission")
	}
	utility.Assert(one.IsDeleted == 0, "Template Already Deleted")
	var payload *string
	if req.Payload != nil {
		payload = unibee.String(utility.MarshalToJsonString(req.Payload))
	}
	var exportColumns *string
	if req.ExportColumns != nil {
		exportColumns = unibee.String(utility.MarshalToJsonString(req.ExportColumns))
	}
	_, err = dao.MerchantBatchExportTemplate.Ctx(ctx).Data(g.Map{
		dao.MerchantBatchExportTemplate.Columns().Name:          req.Name,
		dao.MerchantBatchExportTemplate.Columns().Task:          req.Task,
		dao.MerchantBatchExportTemplate.Columns().Format:        req.Format,
		dao.MerchantBatchExportTemplate.Columns().Payload:       payload,
		dao.MerchantBatchExportTemplate.Columns().ExportColumns: exportColumns,
		dao.MerchantBatchExportTemplate.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.MerchantBatchExportTemplate.Columns().Id, one.Id).Where(dao.MerchantWebhook.Columns().IsDeleted, 0).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("MemberExportTemplate(%v)", one.Id),
		Content:        "Edit",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	one = query.GetMerchantTaskExportTemplateById(ctx, req.TemplateId)
	return &task.EditTemplateRes{Template: bean.SimplifyMerchantBatchExportTemplate(one)}, nil
}
