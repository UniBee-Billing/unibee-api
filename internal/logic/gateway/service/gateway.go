package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/api/bean"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/gateway/api"
	gatewayWebhook "unibee/internal/logic/gateway/webhook"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func SetupGateway(ctx context.Context, merchantId uint64, gatewayName string, gatewayKey string, gatewaySecret string) *entity.MerchantGateway {
	utility.Assert(len(gatewayName) > 0, "gatewayName invalid")
	icon, gatewayType, err := api.GetGatewayWebhookServiceProviderByGatewayName(ctx, gatewayName).GatewayTest(ctx, gatewayKey, gatewaySecret)
	utility.AssertError(err, "gateway test error, key or secret invalid")
	one := query.GetGatewayByGatewayName(ctx, merchantId, gatewayName)
	utility.Assert(one == nil, "exist same gateway")
	if config.GetConfigInstance().IsProd() {
		err = dao.MerchantGateway.Ctx(ctx).
			Where(dao.MerchantGateway.Columns().MerchantId, merchantId).
			Where(dao.MerchantGateway.Columns().GatewayName, gatewayName).
			Where(dao.MerchantGateway.Columns().GatewayKey, gatewayKey).
			Where(dao.MerchantGateway.Columns().GatewaySecret, gatewaySecret).
			OmitEmpty().
			Scan(&one)
		utility.AssertError(err, "system error")
		utility.Assert(one == nil, "same gateway exist")
	}
	one = &entity.MerchantGateway{
		MerchantId:    merchantId,
		GatewayName:   gatewayName,
		Name:          gatewayName,
		GatewayKey:    gatewayKey,
		GatewaySecret: gatewaySecret,
		GatewayType:   gatewayType,
		Logo:          icon,
	}
	result, err := dao.MerchantGateway.Ctx(ctx).Data(one).OmitNil().Insert(one)
	utility.AssertError(err, "system error")
	id, _ := result.LastInsertId()
	one.Id = uint64(id)

	gatewayWebhook.CheckAndSetupGatewayWebhooks(ctx, one.Id)

	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Gateway(%v-%s)", one.Id, one.GatewayName),
		Content:        "Setup",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return one
}

func EditGateway(ctx context.Context, merchantId uint64, gatewayId uint64, gatewayKey string, gatewaySecret string) *entity.MerchantGateway {
	utility.Assert(gatewayId > 0, "gatewayId invalid")
	one := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(one != nil, "gateway not found")
	utility.Assert(one.MerchantId == merchantId, "merchant not match")
	icon, _, err := api.GetGatewayServiceProvider(ctx, gatewayId).GatewayTest(ctx, gatewayKey, gatewaySecret)
	utility.AssertError(err, "gateway test error, key or secret invalid")

	_, err = dao.MerchantGateway.Ctx(ctx).Data(g.Map{
		dao.MerchantGateway.Columns().Logo:          icon,
		dao.MerchantGateway.Columns().GatewaySecret: gatewaySecret,
		dao.MerchantGateway.Columns().GatewayKey:    gatewayKey,
		dao.MerchantGateway.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.MerchantGateway.Columns().Id, one.Id).Update()
	utility.AssertError(err, "system error")

	gatewayWebhook.CheckAndSetupGatewayWebhooks(ctx, one.Id)
	one = query.GetGatewayById(ctx, gatewayId)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Gateway(%v-%s)", one.Id, one.GatewayName),
		Content:        "Edit",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return one
}

func EditGatewayCountryConfig(ctx context.Context, merchantId uint64, gatewayId uint64, countryConfig map[string]bool) (err error) {
	utility.Assert(gatewayId > 0, "gatewayId invalid")
	one := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(one != nil, "gateway not found")
	utility.Assert(one.MerchantId == merchantId, "merchant not match")
	_, err = dao.MerchantGateway.Ctx(ctx).Data(g.Map{
		dao.MerchantGateway.Columns().CountryConfig: utility.MarshalToJsonString(countryConfig),
		dao.MerchantGateway.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.MerchantGateway.Columns().Id, one.Id).Update()
	utility.AssertError(err, "system error")
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Gateway(%v-%s)", one.Id, one.GatewayName),
		Content:        "EditCountryConfig",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return err
}

func IsGatewaySupportCountryCode(ctx context.Context, gateway *entity.MerchantGateway, countryCode string) bool {
	gatewaySimplify := bean.SimplifyGateway(gateway)
	var support = true
	if gatewaySimplify.CountryConfig != nil {
		if _, ok := gatewaySimplify.CountryConfig[countryCode]; ok {
			if !gatewaySimplify.CountryConfig[countryCode] {
				support = false
			}
		}
	}
	return support
}

func GetMerchantAvailableGatewaysByCountryCode(ctx context.Context, merchantId uint64, countryCode string) []*bean.Gateway {
	var availableGateways []*bean.Gateway
	gateways := query.GetMerchantGatewayList(ctx, merchantId)
	for _, one := range gateways {
		if IsGatewaySupportCountryCode(ctx, one, countryCode) {
			availableGateways = append(availableGateways, bean.SimplifyGateway(one))
		}
	}
	return availableGateways
}

type WireTransferSetupReq struct {
	GatewayId     uint64            `json:"gatewayId"  dc:"The id of payment gateway" v:"required"`
	MerchantId    uint64            `json:"merchantId"   dc:"The merchantId of wire transfer" v:"required" `
	Currency      string            `json:"currency"   dc:"The currency of wire transfer " v:"required" `
	MinimumAmount int64             `json:"minimumAmount"   dc:"The minimum amount of wire transfer" v:"required" `
	Bank          *bean.GatewayBank `json:"bank"   dc:"The receiving bank of wire transfer " v:"required" `
}
type WireTransferSetupRes struct {
}

func SetupWireTransferGateway(ctx context.Context, req *WireTransferSetupReq) *entity.MerchantGateway {
	gatewayName := "wire_transfer"
	one := query.GetGatewayByGatewayName(ctx, req.MerchantId, gatewayName)
	utility.Assert(one == nil, "exist same gateway")
	if config.GetConfigInstance().IsProd() {
		err := dao.MerchantGateway.Ctx(ctx).
			Where(dao.MerchantGateway.Columns().MerchantId, req.MerchantId).
			Where(dao.MerchantGateway.Columns().GatewayName, gatewayName).
			OmitEmpty().
			Scan(&one)
		utility.AssertError(err, "system error")
		utility.Assert(one == nil, "same gateway exist")
	}
	one = &entity.MerchantGateway{
		MerchantId:    req.MerchantId,
		GatewayName:   gatewayName,
		Name:          "Wire Transfer",
		Currency:      strings.ToUpper(req.Currency),
		MinimumAmount: req.MinimumAmount,
		GatewayType:   consts.GatewayTypeWireTransfer,
		BankData:      utility.MarshalToJsonString(req.Bank),
	}
	result, err := dao.MerchantGateway.Ctx(ctx).Data(one).OmitNil().Insert(one)
	utility.AssertError(err, "system error")
	id, _ := result.LastInsertId()
	one.Id = uint64(id)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Gateway(%v-%s)", one.Id, one.GatewayName),
		Content:        "Setup-WireTransfer",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return one
}

func EditWireTransferGateway(ctx context.Context, req *WireTransferSetupReq) *entity.MerchantGateway {
	utility.Assert(req.GatewayId > 0, "gatewayId invalid")
	one := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(one != nil, "gateway not found")
	utility.Assert(one.MerchantId == req.MerchantId, "merchant not match")

	_, err := dao.MerchantGateway.Ctx(ctx).Data(g.Map{
		dao.MerchantGateway.Columns().BankData:      utility.MarshalToJsonString(req.Bank),
		dao.MerchantGateway.Columns().Currency:      strings.ToUpper(req.Currency),
		dao.MerchantGateway.Columns().MinimumAmount: req.MinimumAmount,
		dao.MerchantGateway.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.MerchantGateway.Columns().Id, one.Id).Update()
	utility.AssertError(err, "system error")
	one = query.GetGatewayById(ctx, req.GatewayId)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Gateway(%v-%s)", one.Id, one.GatewayName),
		Content:        "Edit-WireTransfer",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return one
}
