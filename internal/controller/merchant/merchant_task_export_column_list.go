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
	var lowerColumnHeaders = make(map[string]string)
	for _, value := range columns {
		var column = fmt.Sprintf("%s", value)
		if utility.IsStartUpper(fmt.Sprintf("%s", column)) {
			column = utility.ToFirstCharLowerCase(column)
		}
		lowerColumns = append(lowerColumns, column)
		lowerColumnHeaders[column] = batch.AddSpaceBeforeUpperCaseExceptFirst(fmt.Sprintf("%s", value))
	}
	var lowerColumnCommentMap = make(map[string]string)
	for key, value := range commentMap {
		if utility.IsStartUpper(key) {
			lowerColumnCommentMap[utility.ToFirstCharLowerCase(key)] = value
		} else {
			lowerColumnCommentMap[key] = value
		}
	}
	return &task.ExportColumnListRes{Columns: lowerColumns, ColumnComments: lowerColumnCommentMap, ColumnHeaders: lowerColumnHeaders}, nil
}
