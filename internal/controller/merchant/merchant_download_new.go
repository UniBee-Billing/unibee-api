package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/batch"
	"unibee/utility"

	"unibee/api/merchant/download"
)

func (c *ControllerDownload) New(ctx context.Context, req *download.NewReq) (res *download.NewRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "no permission")
	err = batch.NewBatchDownloadTask(ctx, &batch.MerchantBatchTaskInternalRequest{
		MerchantId: _interface.GetMerchantId(ctx),
		MemberId:   _interface.Context().Get(ctx).MerchantMember.Id,
		Task:       req.Task,
		Payload:    req.Payload,
	})
	if err != nil {
		return nil, err
	}
	return &download.NewRes{}, nil
}
