package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/gateway/api"
	gatewayWebhook "unibee/internal/logic/gateway/webhook"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func SetupGateway(ctx context.Context, merchantId uint64, gatewayName string, gatewayKey string, gatewaySecret string) {
	utility.Assert(len(gatewayName) > 0, "gatewayName invalid")
	if strings.Compare(strings.ToLower(gatewayName), "stripe") == 0 {
		utility.Assert(len(gatewaySecret) > 0, "invalid gatewaySecret")
		utility.Assert(strings.HasPrefix(gatewaySecret, "sk_"), "invalid gatewaySecret, should start with 'sk_'")
		err := api.GetGatewayWebhookServiceProviderByGatewayName(ctx, gatewayName).GatewayTest(ctx, gatewayKey, gatewaySecret)
		utility.AssertError(err, "gateway test error, key or secret invalid")

	} else if strings.Compare(strings.ToLower(gatewayName), "paypal") == 0 {
		utility.Assert(false, "not support")
	} else {
		utility.Assert(false, "invalid gatewayName")
	}
	one := query.GetGatewayByGatewayName(ctx, gatewayName)
	utility.Assert(one == nil, "exist same gateway")
	one = &entity.MerchantGateway{
		MerchantId:    merchantId,
		GatewayName:   gatewayName,
		Name:          gatewayName,
		GatewayKey:    gatewayKey,
		GatewaySecret: gatewaySecret,
	}
	result, err := dao.MerchantGateway.Ctx(ctx).Data(one).OmitNil().Insert(one)
	utility.AssertError(err, "system error")
	id, _ := result.LastInsertId()
	one.Id = uint64(id)

	gatewayWebhook.CheckAndSetupGatewayWebhooks(ctx, one.Id)
}

func EditGateway(ctx context.Context, merchantId uint64, gatewayId uint64, gatewayKey string, gatewaySecret string) {
	utility.Assert(gatewayId > 0, "gatewayId invalid")
	one := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(one != nil, "gateway not found")
	err := api.GetGatewayServiceProvider(ctx, gatewayId).GatewayTest(ctx, gatewayKey, gatewaySecret)
	utility.AssertError(err, "gateway test error, key or secret invalid")

	_, err = dao.MerchantGateway.Ctx(ctx).Data(g.Map{
		dao.MerchantGateway.Columns().GatewaySecret: gatewaySecret,
		dao.MerchantGateway.Columns().GatewayKey:    gatewayKey,
		dao.MerchantGateway.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.MerchantGateway.Columns().Id, one.Id).Update()
	utility.AssertError(err, "system error")

	gatewayWebhook.CheckAndSetupGatewayWebhooks(ctx, one.Id)
}
