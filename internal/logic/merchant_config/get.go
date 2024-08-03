package merchant_config

import (
	"context"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func GetMerchantConfig(ctx context.Context, merchantId uint64, configKey string) *entity.MerchantConfig {
	utility.Assert(merchantId > 0, "invalid merchantId")
	if len(configKey) == 0 {
		return nil
	}
	var one *entity.MerchantConfig
	err := dao.MerchantConfig.Ctx(ctx).
		Where(dao.MerchantConfig.Columns().MerchantId, merchantId).
		Where(dao.MerchantConfig.Columns().ConfigKey, configKey).
		Scan(&one)
	if err != nil {
		one = nil
	}
	return one
}
