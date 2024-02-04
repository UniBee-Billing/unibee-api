package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetUserAccountById(ctx context.Context, id uint64) (one *entity.UserAccount) {
	if id <= 0 {
		return nil
	}
	err := dao.UserAccount.Ctx(ctx).Where(entity.UserAccount{Id: id}).OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetUserAccountByEmail(ctx context.Context, email string) (one *entity.UserAccount) {
	if len(email) == 0 {
		return nil
	}
	err := dao.UserAccount.Ctx(ctx).Where(entity.UserAccount{Email: email}).OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}
