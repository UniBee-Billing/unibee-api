package query

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/merchant_config"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

const MerchantDefaultProductKey = "KEYMERCHANTDEFAULTPRODUCT"

func GetDefaultProduct(ctx context.Context, merchantId uint64) *entity.Product {
	config := merchant_config.GetMerchantConfig(ctx, merchantId, MerchantDefaultProductKey)
	if config != nil && len(config.ConfigValue) > 0 {
		var one *entity.Product
		_ = utility.UnmarshalFromJsonString(config.ConfigValue, &one)
		if one != nil {
			one.Id = 0
			return one
		}
	}
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

func GetProductById(ctx context.Context, id uint64, merchantId uint64) (one *entity.Product) {
	if id <= 0 {
		return GetDefaultProduct(ctx, merchantId)
	}
	err := dao.Product.Ctx(ctx).Where(dao.Product.Columns().Id, id).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetProductByProductName(ctx context.Context, productName string, merchantId uint64) (one *entity.Product) {
	if len(productName) <= 0 || productName == "Default" {
		return GetDefaultProduct(ctx, merchantId)
	}
	err := dao.Product.Ctx(ctx).Where(dao.Product.Columns().ProductName, productName).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
