package merchant

import (
	"context"
	"unibee/api/merchant/download"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/batch"
	"unibee/utility"
)

func (c *ControllerDownload) List(ctx context.Context, req *download.ListReq) (res *download.ListRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "no permission")
	list, total := batch.MerchantBatchTaskList(ctx, _interface.GetMerchantId(ctx), _interface.Context().Get(ctx).MerchantMember.Id, req.Page, req.Count)
	return &download.ListRes{Downloads: list, Total: total}, nil
}
