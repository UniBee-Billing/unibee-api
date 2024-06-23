package invoice

import (
	"context"
	"fmt"
	entity "unibee/internal/model/entity/oversea_pay"
)

type TaskInvoice struct {
}

func (t TaskInvoice) TaskName() string {
	return fmt.Sprintf("InvoiceExport")
}

func (t TaskInvoice) Header() interface{} {
	//TODO implement me
	panic("implement me")
}

func (t TaskInvoice) PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error) {
	//TODO implement me
	panic("implement me")
}
