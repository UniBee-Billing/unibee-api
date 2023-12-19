package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetMerchantInfoById(ctx context.Context, id int64) (one *entity.MerchantInfo) {
	err := dao.MerchantInfo.Ctx(ctx).Where(entity.MerchantInfo{Id: id}).OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}
