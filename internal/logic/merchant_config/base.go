package merchant_config

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

func SetMerchantConfig(ctx context.Context, merchantId uint64, configKey string, configValue string) error {
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
	var one *entity.MerchantConfig
	err := dao.MerchantConfig.Ctx(ctx).
		Where(entity.MerchantConfig{MerchantId: merchantId}).
		Where(entity.MerchantConfig{ConfigKey: configKey}).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return one
}
