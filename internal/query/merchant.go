package query

import (
	"context"
	"fmt"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
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

func GetMerchantByOwnerEmail(ctx context.Context, email string) (one *entity.Merchant) {
	if len(email) == 0 {
		return nil
	}
	var member *entity.MerchantMember
	err := dao.MerchantMember.Ctx(ctx).
		Where(dao.MerchantMember.Columns().Email, email).
		Where(dao.MerchantMember.Columns().Role, "Owner").
		Scan(&member)
	if err != nil || member == nil {
		return nil
	}
	return GetMerchantById(ctx, member.MerchantId)
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

func GetMerchantOwnerMember(ctx context.Context, merchantId uint64) (one *entity.MerchantMember) {
	if merchantId <= 0 {
		return nil
	}
	err := dao.MerchantMember.Ctx(ctx).
		Where(dao.MerchantMember.Columns().MerchantId, merchantId).
		Where(dao.MerchantMember.Columns().Role, "Owner").
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetMerchantMemberCount(ctx context.Context, merchantId uint64) int {
	if merchantId <= 0 {
		return 100
	}
	count, err := dao.MerchantMember.Ctx(ctx).
		Where(dao.MerchantMember.Columns().MerchantId, merchantId).
		Where(dao.MerchantMember.Columns().IsDeleted, 0).
		Count()
	if err != nil {
		return 100
	}
	return count
}

func GetMerchantMembersByAuthJsProvider(ctx context.Context, Provider string, ProviderId string) (list []*entity.MerchantMember) {
	if len(Provider) == 0 || len(ProviderId) == 0 {
		return nil
	}
	err := dao.MerchantMember.Ctx(ctx).
		WhereLike(dao.MerchantMember.Columns().AuthJs, "%"+fmt.Sprintf("%s##%s", Provider, ProviderId)+"%").
		Scan(&list)
	if err != nil {
		return nil
	}
	return list
}
