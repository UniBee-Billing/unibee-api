package multi_currency

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/consumer/webhook/log"
	"unibee/internal/logic/merchant_config"
	"unibee/internal/logic/multi_currencies"
	"unibee/internal/logic/multi_currencies/currency_exchange"
	"unibee/internal/query"
)

func TaskForSyncMerchantsMultiCurrencyConfigs(ctx context.Context) {
	g.Log().Infof(ctx, "TaskForSyncMerchantsMultiCurrencyConfigs start")
	merchants, _ := query.GetMerchantList(ctx)
	for _, merchant := range merchants {
		data := merchant_config.GetMerchantConfig(ctx, merchant.Id, currency_exchange.MerchantMultiCurrenciesConfig)
		if data != nil {
			SyncMerchantMultiCurrencyConfigs(merchant.Id)
		}
	}
	g.Log().Infof(ctx, "TaskForSyncMerchantsMultiCurrencyConfigs end")
}

func SyncMerchantMultiCurrencyConfigs(merchantId uint64) {
	go func() {
		ctx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				log.PrintPanic(ctx, err)
				return
			}
		}()
		g.Log().Infof(ctx, "SyncMerchantMultiCurrencyConfigs:%s", merchantId)
		multi_currencies.UpdateMerchantMultiCurrenciesConfigExchangeRate(ctx, merchantId)
	}()
}
