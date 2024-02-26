package _interface

import (
	"context"
	entity "unibee/internal/model/entity/oversea_pay"
)

type IOpenApi interface {
	GetOpenApiConfig(ctx context.Context, key string) (res *entity.OpenApiConfig)
}

var localOpenApi IOpenApi

func OpenApi() IOpenApi {
	if localOpenApi == nil {
		panic("implement not found for interface ISession, forgot register?")
	}
	return localOpenApi
}

func RegisterOpenApi(i IOpenApi) {
	localOpenApi = i
}
