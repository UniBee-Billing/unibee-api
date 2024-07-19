package merchant

import (
	"context"
	"unibee/internal/logic/batch"

	"unibee/api/merchant/task"
)

func (c *ControllerTask) ExportColumnList(ctx context.Context, req *task.ExportColumnListReq) (res *task.ExportColumnListRes, err error) {
	columns, commentMap := batch.ExportColumnList(ctx, req.Task)
	return &task.ExportColumnListRes{Columns: columns, ColumnComments: commentMap}, nil
}
