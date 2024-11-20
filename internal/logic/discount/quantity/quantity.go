package quantity

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

const CacheKeyPrefixDiscountUsageCount = "CacheKeyPrefixDiscountUsageCountV2"

func GetDiscountQuantityUsedCountCacheKey(code string) string {
	key := fmt.Sprintf("%s_%s", CacheKeyPrefixDiscountUsageCount, code)
	return key
}

func GetDiscountQuantityUsedCount(ctx context.Context, id uint64) (count int) {
	one := query.GetDiscountById(ctx, id)
	if one == nil {
		return 0
	}
	key := GetDiscountQuantityUsedCountCacheKey(one.Code)
	get, _ := g.Redis().Get(ctx, key)
	if get != nil && !get.IsNil() {
		return get.Int()
	}
	count = getDiscountQuantityUsedCountFromDatabase(ctx, one)
	_, _ = g.Redis().Set(ctx, key, count)
	return count
}

func getDiscountQuantityUsedCountFromDatabase(ctx context.Context, one *entity.MerchantDiscountCode) (count int) {
	if one == nil {
		return count
	}
	count, err := dao.MerchantUserDiscountCode.Ctx(ctx).
		Where(dao.MerchantUserDiscountCode.Columns().Code, one.Code).
		Where(dao.MerchantUserDiscountCode.Columns().Status, 1).
		Where(dao.MerchantUserDiscountCode.Columns().Recurring, 0).
		Where(dao.MerchantUserDiscountCode.Columns().IsDeleted, 0).
		Count()
	if err != nil {
		return 0
	}
	return count
}
