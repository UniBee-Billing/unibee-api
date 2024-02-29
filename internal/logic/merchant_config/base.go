package merchant_config

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

func SetMerchantConfig(ctx context.Context, merchantId uint64, configKey string, configValue string) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(len(configKey) > 0, "invalid key")
	utility.Assert(len(configValue) > 0, "invalid value")
	one := &entity.MerchantConfig{
		MerchantId:  merchantId,
		ConfigKey:   configKey,
		ConfigValue: configValue,
	}
	_, err := dao.MerchantConfig.Ctx(ctx).Data(one).OmitEmpty().Save(one)
	if err != nil {
		err = gerror.Newf(`SetMerchantConfig %s`, err)
		return err
	}
	return nil
}

func GetMerchantConfig(ctx context.Context, merchantId uint64, configKey string) *entity.MerchantConfig {
	if merchantId == 0 {
		return nil
	}
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
