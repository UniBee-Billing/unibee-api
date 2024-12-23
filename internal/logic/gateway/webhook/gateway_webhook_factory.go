package webhook

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	dao "unibee/internal/dao/default"
	"unibee/internal/interface"
	context2 "unibee/internal/interface/context"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func GetGatewayWebhookServiceProvider(ctx context.Context, gatewayId uint64) (one _interface.GatewayWebhookInterface) {
	proxy := &GatewayWebhookProxy{}
	proxy.Gateway = query.GetGatewayById(ctx, gatewayId)
	proxy.GatewayName = proxy.Gateway.GatewayName
	utility.Assert(proxy.Gateway != nil, fmt.Sprintf("gateway not found %d", gatewayId))
	return proxy
}

func GetGatewayWebhookServiceProviderByGatewayName(ctx context.Context, gatewayName string) (one _interface.GatewayWebhookInterface) {
	proxy := &GatewayWebhookProxy{}
	proxy.GatewayName = gatewayName
	return proxy
}

func CheckAndSetupGatewayWebhooks(ctx context.Context, gatewayId uint64) {
	gateway := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(gateway != nil, fmt.Sprintf("gateway not found %d", gatewayId))
	err := GetGatewayWebhookServiceProvider(ctx, gateway.Id).GatewayCheckAndSetupWebhook(ctx, gateway)
	if err != nil {
		g.Log().Errorf(ctx, "CheckAndSetupGatewayWebhooks GatewayName:%s Error:%s", gateway.GatewayName, err)
	} else {
		g.Log().Infof(ctx, "CheckAndSetupGatewayWebhooks GatewayName:%s Success", gateway.GatewayName)
	}
	utility.AssertError(err, "CheckAndSetupGatewayWebhooks Error")
	if context2.Context().Get(ctx) != nil && (context2.Context().Get(ctx).MerchantMember != nil || context2.Context().Get(ctx).IsOpenApiCall) {
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     gateway.MerchantId,
			Target:         fmt.Sprintf("Gateway(%v-%s)", gateway.Id, gateway.GatewayName),
			Content:        "SetupWebhook",
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)
	}
}

func SetupAllWebhooksBackground() {
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
				printChannelPanic(ctx, err)
				return
			}
		}()

		var list []*entity.MerchantGateway
		err = dao.MerchantGateway.Ctx(ctx).
			WhereIn(dao.MerchantGateway.Columns().GatewayName, []string{"stripe", "paypal"}).
			Where(dao.MerchantGateway.Columns().IsDeleted, 0).
			Scan(&list)
		if err != nil {
			g.Log().Errorf(ctx, "SetupAllWebhooksBackground error:%s", err)
		}
		for _, gateway := range list {
			CheckAndSetupGatewayWebhooks(ctx, gateway.Id)
		}
	}()
}
