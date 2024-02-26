package openapi

import (
	"context"
	"unibee/internal/interface"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
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
