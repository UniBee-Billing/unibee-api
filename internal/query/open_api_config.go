package query

import (
	"context"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetOpenApiConfigByKey(ctx context.Context, key string) (res *entity.OpenApiConfig) {
	if len(key) == 0 {
		return nil
	}
	err := dao.OpenApiConfig.Ctx(ctx).Where(dao.OpenApiConfig.Columns().ApiKey, key).OmitEmpty().Scan(&res)
	if err != nil {
		return nil
	}
	return res
}
