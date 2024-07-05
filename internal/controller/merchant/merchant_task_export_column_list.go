package merchant

import (
	"context"
	"unibee/internal/logic/batch"

	"unibee/api/merchant/task"
)

func (c *ControllerTask) ExportColumnList(ctx context.Context, req *task.ExportColumnListReq) (res *task.ExportColumnListRes, err error) {
	return &task.ExportColumnListRes{Columns: batch.ExportColumnList(ctx, req.Task)}, nil
}
