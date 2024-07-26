package query

import (
	"context"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetProductById(ctx context.Context, id uint64) (one *entity.Product) {
	if id <= 0 {
		return nil
	}
	err := dao.Product.Ctx(ctx).Where(dao.Product.Columns().Id, id).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetProductsByIds(ctx context.Context, ids []int64) (list []*entity.Product) {
	err := dao.Product.Ctx(ctx).WhereIn(dao.Product.Columns().Id, ids).OmitEmpty().Scan(&list)
	if err != nil {
		return nil
	}
	return list
}
