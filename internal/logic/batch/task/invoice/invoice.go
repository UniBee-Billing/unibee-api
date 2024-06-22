package invoice

import (
	"context"
	"fmt"
	entity "unibee/internal/model/entity/oversea_pay"
)

type TaskInvoice struct {
}

func (t TaskInvoice) TableName(task *entity.MerchantBatchTask) string {
	return fmt.Sprintf("Invoice_%v_%v_%v", task.Id, task.MerchantId, task.MemberId)
}

func (t TaskInvoice) Header() []interface{} {
	//TODO implement me
	panic("implement me")
}

func (t TaskInvoice) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([][]interface{}, error) {
	//TODO implement me
	panic("implement me")
}
