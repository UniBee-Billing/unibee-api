package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/operation_log"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/task"
)

func (c *ControllerTask) DeleteTemplate(ctx context.Context, req *task.DeleteTemplateReq) (res *task.DeleteTemplateRes, err error) {
	utility.Assert(req.TemplateId > 0, "invalid templateId")
	one := query.GetMerchantTaskExportTemplateById(ctx, req.TemplateId)
	utility.Assert(one != nil, "invalid templateId")
	utility.Assert(one.MerchantId == _interface.GetMerchantId(ctx), "No Permission")
	if _interface.Context().Get(ctx) != nil && _interface.Context().Get(ctx).MerchantMember != nil {
		utility.Assert(one.MemberId == _interface.Context().Get(ctx).MerchantMember.Id, "No Permission")
	}
	if one.IsDeleted != 0 {
		// already deleted
		return &task.DeleteTemplateRes{}, nil
	}
	_, err = dao.MerchantBatchExportTemplate.Ctx(ctx).Data(g.Map{
		dao.MerchantBatchExportTemplate.Columns().IsDeleted: 1,
		dao.MerchantBatchExportTemplate.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantBatchExportTemplate.Columns().Id, one.Id).Where(dao.MerchantWebhook.Columns().IsDeleted, 0).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("MemberExportTemplate(%v)", one.Id),
		Content:        "Delete",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)

	return &task.DeleteTemplateRes{}, nil
}
