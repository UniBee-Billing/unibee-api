package merchant

import (
	"context"
	"fmt"
	"unibee/internal/logic/batch"
	"unibee/utility"

	"unibee/api/merchant/task"
)

func (c *ControllerTask) ExportColumnList(ctx context.Context, req *task.ExportColumnListReq) (res *task.ExportColumnListRes, err error) {
	columns, commentMap := batch.ExportColumnList(ctx, req.Task)
	lowerColumns := make([]interface{}, 0)
	for _, value := range columns {
		if utility.IsStartUpper(fmt.Sprintf("%s", value)) {
			lowerColumns = append(lowerColumns, utility.ToFirstCharLowerCase(fmt.Sprintf("%s", value)))
		} else {
			lowerColumns = append(lowerColumns, fmt.Sprintf("%s", value))
		}
	}
	var lowerColumnCommentMap = make(map[string]string)
	for key, value := range commentMap {
		if utility.IsStartUpper(key) {
			lowerColumnCommentMap[utility.ToFirstCharLowerCase(key)] = value
		} else {
			lowerColumnCommentMap[key] = value
		}
	}
	return &task.ExportColumnListRes{Columns: lowerColumns, ColumnComments: lowerColumnCommentMap}, nil
}
