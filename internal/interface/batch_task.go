package _interface

import (
	"context"
	entity "unibee/internal/model/entity/oversea_pay"
)

type BatchExportTask interface {
	TaskName() string
	Header() interface{}
	PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error)
}

type BatchImportTask interface {
	TaskName() string
	TemplateHeader() interface{}
	ImportRow(ctx context.Context, task *entity.MerchantBatchTask, data map[string]string) (interface{}, error)
}
