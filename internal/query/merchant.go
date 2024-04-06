package query

import (
	"context"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

func GetMerchantByApiKey(ctx context.Context, apiKey string) (one *entity.Merchant) {
	if len(apiKey) <= 0 {
		return nil
	}
	err := dao.Merchant.Ctx(ctx).
		Where(dao.Merchant.Columns().ApiKey, apiKey).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetMerchantById(ctx context.Context, id uint64) (one *entity.Merchant) {
	if id <= 0 {
		return nil
	}
	err := dao.Merchant.Ctx(ctx).
		Where(dao.Merchant.Columns().Id, id).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetMerchantByHost(ctx context.Context, host string) (one *entity.Merchant) {
	if len(host) <= 0 {
		return nil
	}
	err := dao.Merchant.Ctx(ctx).
		Where(dao.Merchant.Columns().Host, host).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetMerchantList(ctx context.Context) (list []*entity.Merchant, err error) {
	err = dao.Merchant.Ctx(ctx).
		Scan(&list)
	if err != nil {
		return make([]*entity.Merchant, 0), err
	}
	return list, nil
}

func GetActiveMerchantList(ctx context.Context) (list []*entity.Merchant) {
	err := dao.Merchant.Ctx(ctx).
		Where(dao.Merchant.Columns().IsDeleted, 0).
		Scan(&list)
	if err != nil {
		return make([]*entity.Merchant, 0)
	}
	return
}

func GetMerchantMemberById(ctx context.Context, id uint64) (one *entity.MerchantMember) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantMember.Ctx(ctx).
		Where(dao.MerchantMember.Columns().Id, id).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetMerchantMemberByEmail(ctx context.Context, email string) (one *entity.MerchantMember) {
	if len(email) == 0 {
		return nil
	}
	err := dao.MerchantMember.Ctx(ctx).
		Where(dao.MerchantMember.Columns().Email, email).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}
