package openapi

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/interface"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

type SOpenApi struct{}

func (s SOpenApi) GetOpenApiConfig(ctx context.Context, key string) (res *entity.OpenApiConfig) {
	err := dao.OpenApiConfig.Ctx(ctx).Where(entity.OpenApiConfig{ApiKey: key}).OmitEmpty().Scan(&res)
	if err != nil {
		return nil
	}
	return res

}

func init() {
	_interface.RegisterOpenApi(New())
}

func New() *SOpenApi {
	return &SOpenApi{}
}
