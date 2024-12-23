package merchant

import (
	"context"
	"unibee/api/merchant/task"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/batch"
	"unibee/utility"
)

func (c *ControllerTask) List(ctx context.Context, req *task.ListReq) (res *task.ListRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "no permission")
	list, total := batch.MerchantBatchTaskList(ctx, _interface.GetMerchantId(ctx), _interface.Context().Get(ctx).MerchantMember.Id, req.Page, req.Count)
	return &task.ListRes{Downloads: list, Total: total}, nil
}
