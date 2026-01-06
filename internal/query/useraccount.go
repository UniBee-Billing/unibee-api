package query

import (
	"context"
	"strings"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	entity "unibee/internal/model/entity/default"
)

func GetUserAccountsByIds(ctx context.Context, ids []uint64) (list []*entity.UserAccount) {
	if len(ids) <= 0 {
		return make([]*entity.UserAccount, 0)
	}
	err := dao.UserAccount.Ctx(ctx).
		WhereIn(dao.UserAccount.Columns().Id, ids).
		Scan(&list)
	if err != nil {
		return make([]*entity.UserAccount, 0)
	}
	return
}

func GetUserAccountById(ctx context.Context, id uint64) (one *entity.UserAccount) {
	if id <= 0 {
		return nil
	}
	one = _interface.GetUserFromPreloadContext(ctx, id)
	if one != nil {
		return one
	}
	err := dao.UserAccount.Ctx(ctx).
		Where(dao.UserAccount.Columns().Id, id).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetUserAccountByEmail(ctx context.Context, merchantId uint64, email string) (one *entity.UserAccount) {
	if len(email) == 0 {
		return nil
	}
	email = strings.TrimSpace(email)
	err := dao.UserAccount.Ctx(ctx).
		Where(dao.UserAccount.Columns().Email, email).
		Where(dao.UserAccount.Columns().MerchantId, merchantId).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetUserAccountByExternalUserId(ctx context.Context, merchantId uint64, externalUserId string) (one *entity.UserAccount) {
	if len(externalUserId) <= 0 {
		return nil
	}
	err := dao.UserAccount.Ctx(ctx).
		Where(dao.UserAccount.Columns().ExternalUserId, externalUserId).
		Where(dao.UserAccount.Columns().MerchantId, merchantId).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}
