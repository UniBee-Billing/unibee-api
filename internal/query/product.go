package query

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetDefaultProduct() *entity.Product {
	return &entity.Product{
		Id:          0,
		GmtCreate:   gtime.Now(),
		GmtModify:   gtime.Now(),
		CompanyId:   0,
		MerchantId:  0,
		ProductName: "Default",
		Description: "System Default Product",
		ImageUrl:    "",
		HomeUrl:     "",
		Status:      1,
		IsDeleted:   0,
		CreateTime:  gtime.Now().Timestamp(),
		MetaData:    "",
	}
}

func GetProductById(ctx context.Context, id uint64) (one *entity.Product) {
	if id <= 0 {
		return GetDefaultProduct()
	}
	err := dao.Product.Ctx(ctx).Where(dao.Product.Columns().Id, id).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
