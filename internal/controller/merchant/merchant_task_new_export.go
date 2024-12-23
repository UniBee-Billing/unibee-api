package merchant

import (
	"context"
	"fmt"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/batch"
	"unibee/utility"

	"unibee/api/merchant/task"
)

func (c *ControllerTask) New(ctx context.Context, req *task.NewReq) (res *task.NewRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "no permission")
	UpperColumns := make([]string, 0)
	for _, value := range req.ExportColumns {
		if utility.IsStartLower(fmt.Sprintf("%s", value)) {
			UpperColumns = append(UpperColumns, utility.ToFirstCharUpperCase(value))
		} else {
			UpperColumns = append(UpperColumns, fmt.Sprintf("%s", value))
		}
	}
	err = batch.NewBatchExportTask(ctx, &batch.MerchantBatchExportTaskInternalRequest{
		MerchantId:    _interface.GetMerchantId(ctx),
		MemberId:      _interface.Context().Get(ctx).MerchantMember.Id,
		Task:          req.Task,
		Payload:       req.Payload,
		ExportColumns: UpperColumns,
		Format:        req.Format,
	})
	if err != nil {
		return nil, err
	}
	return &task.NewRes{}, nil
}
