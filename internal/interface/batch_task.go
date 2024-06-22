package _interface

import (
	"context"
	entity "unibee/internal/model/entity/oversea_pay"
)

type BatchTask interface {
	TableName() string
	Header() []interface{}
	PageData(ctx context.Context, page int, count int, task *entity.MerchantBatchTask) ([][]interface{}, error)
}
