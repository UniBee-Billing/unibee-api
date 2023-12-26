package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetUserAccountById(ctx context.Context, id uint64) (one *entity.UserAccount) {
	err := dao.UserAccount.Ctx(ctx).Where(entity.UserAccount{Id: uint64(id)}).OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}
