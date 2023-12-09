package openapi

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/service"
)

type sOpenApi struct{}

func (s sOpenApi) GetOpenApiConfig(ctx context.Context, key string) (res *entity.OpenApiConfig) {
	err := dao.OpenApiConfig.Ctx(ctx).Where(entity.OpenApiConfig{ApiKey: key}).OmitEmpty().Scan(&res)
	if err != nil {
		return nil
	}
	return res

}

func init() {
	service.RegisterOpenApi(New())
}

func New() *sOpenApi {
	return &sOpenApi{}
}
