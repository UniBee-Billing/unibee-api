package _interface

import (
	"context"
	entity "unibee/internal/model/entity/default"
)

type BatchExportTask interface {
	TaskName() string
	Header() interface{}
	PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([]interface{}, error)
}

type BatchImportTask interface {
	TaskName() string
	TemplateVersion() string
	TemplateHeader() interface{}
	ImportRow(ctx context.Context, task *entity.MerchantBatchTask, row map[string]string) (interface{}, error)
}
