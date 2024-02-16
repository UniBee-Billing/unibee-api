package openapi

import (
	"context"
	"unibee-api/internal/interface"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
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
