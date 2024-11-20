package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/log"
	"unibee/internal/logic/email"
	"unibee/internal/logic/gateway/service"
	"unibee/internal/logic/vat_gateway"
	"unibee/internal/query"
	"unibee/utility"
)

func SetupForCloudMode(ctx context.Context, merchantId uint64) error {
	if config.GetConfigInstance().Mode == "cloud" {
		//if cloud version setup default sendgrid, vat, stripe gateway
		{
			name, data := email.GetDefaultMerchantEmailConfig(ctx, consts.CloudModeManagerMerchantId)
			utility.Assert(len(name) > 0 && len(data) > 0, "Server Error")
			err := email.SetupMerchantEmailConfig(ctx, merchantId, name, data, true)
			if err != nil {
				return err
			}
		}
		{
			name, data := vat_gateway.GetDefaultMerchantVatConfig(ctx, consts.CloudModeManagerMerchantId)
			utility.Assert(len(name) > 0 && len(data) > 0, "Server Error")
			err := vat_gateway.SetupMerchantVatConfig(ctx, merchantId, name, data, true)
			if err != nil {
				return err
			}
			err = vat_gateway.InitMerchantDefaultVatGateway(ctx, merchantId)
			if err != nil {
				return err
			}
		}
		{
			stripeGateway := query.GetGatewayByGatewayName(ctx, consts.CloudModeManagerMerchantId, "stripe")
			if stripeGateway != nil {
				service.SetupGateway(ctx, merchantId, stripeGateway.GatewayName, stripeGateway.GatewayKey, stripeGateway.GatewaySecret)
			}
		}
	}
	return nil
}

func ReloadAllMerchantsCacheForSDKAuthBackground() {
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
		list := query.GetActiveMerchantList(ctx)
		if len(list) > 0 {
			_, _ = g.Redis().Set(ctx, "UniBee#AllMerchants", utility.MarshalToJsonString(list))
			for _, one := range list {
				PutOpenApiKeyToCache(ctx, one.ApiKey, one.Id)
			}
		}
	}()
}
