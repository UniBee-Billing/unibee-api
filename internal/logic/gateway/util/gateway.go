package util

import (
	"context"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	//"unibee/internal/query"
)

//	func GetGatewayById(ctx context.Context, id uint64) (one *entity.MerchantGateway) {
//		if id <= 0 {
//			return nil
//		}
//		err := dao.MerchantGateway.Ctx(ctx).
//			Where(dao.MerchantGateway.Columns().Id, id).
//			Scan(&one)
//		if err != nil {
//			one = nil
//		}
//		return
//	}

func GetGatewayById(ctx context.Context, id uint64) (gateway *entity.MerchantGateway) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantGateway.Ctx(ctx).
		Where(dao.MerchantGateway.Columns().Id, id).
		Scan(&gateway)
	if err != nil {
		gateway = nil
	}
	return
}
