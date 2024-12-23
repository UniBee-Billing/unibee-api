package merchant

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/task"
)

func (c *ControllerTask) NewTemplate(ctx context.Context, req *task.NewTemplateReq) (res *task.NewTemplateRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "No Permission")
	utility.Assert(len(req.Task) > 0, "Invalid Task")
	one := &entity.MerchantBatchExportTemplate{
		MerchantId:    _interface.GetMerchantId(ctx),
		MemberId:      _interface.Context().Get(ctx).MerchantMember.Id,
		Name:          req.Name,
		Task:          req.Task,
		Format:        req.Format,
		Payload:       utility.MarshalToJsonString(req.Payload),
		ExportColumns: utility.MarshalToJsonString(req.ExportColumns),
		IsDeleted:     0,
		CreateTime:    gtime.Now().Timestamp(),
	}
	result, err := dao.MerchantBatchExportTemplate.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		g.Log().Errorf(ctx, "New MerchantBatchExportTemplate Insert err:%s", err.Error())
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(id)

	return &task.NewTemplateRes{Template: bean.SimplifyMerchantBatchExportTemplate(one)}, nil
}
