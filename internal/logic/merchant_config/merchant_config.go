package merchant_config

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func SetMerchantConfig(ctx context.Context, merchantId uint64, configKey string, configValue string) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(len(configKey) > 0, "invalid key")
	one := GetMerchantConfig(ctx, merchantId, configKey)
	if one != nil {
		_, err := dao.MerchantConfig.Ctx(ctx).Data(g.Map{
			dao.MerchantConfig.Columns().ConfigValue: configValue,
			dao.MerchantConfig.Columns().GmtModify:   gtime.Now(),
		}).Where(dao.MerchantConfig.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			g.Log().Errorf(ctx, "SetMerchantConfig error:%s", err.Error())
			err = gerror.Newf(`SetMerchantConfig %s`, err.Error())
			return err
		}
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     merchantId,
			Target:         fmt.Sprintf("MerchantConfig(%v)", configKey),
			Content:        fmt.Sprintf("NewValue(%s)", configValue),
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)
		return nil
	} else {
		one = &entity.MerchantConfig{
			MerchantId:  merchantId,
			ConfigKey:   configKey,
			ConfigValue: configValue,
		}
		_, err := dao.MerchantConfig.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			g.Log().Errorf(ctx, "SetMerchantConfig error:%s", err.Error())
			err = gerror.Newf(`SetMerchantConfig %s`, err.Error())
			return err
		}
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     merchantId,
			Target:         fmt.Sprintf("MerchantConfig(%v)", configKey),
			Content:        fmt.Sprintf("NewValue(%s)", configValue),
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)
		return nil
	}
}

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
