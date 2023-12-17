package openapi

import (
	"context"
	"go-oversea-pay/internal/interface"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
)

type SOpenApi struct{}

func (s SOpenApi) GetOpenApiConfig(ctx context.Context, key string) (res *entity.OpenApiConfig) {
	return query.GetOpenApiConfigByKey(ctx, key)
}

func init() {
	_interface.RegisterOpenApi(New())
}

func New() *SOpenApi {
	return &SOpenApi{}
}
