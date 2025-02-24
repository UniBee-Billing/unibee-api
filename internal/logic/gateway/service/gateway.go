package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/api/bean/detail"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/gateway/api"
	gatewayWebhook "unibee/internal/logic/gateway/webhook"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"
)

func JoinGatewayIcon(gatewayIcon *[]string) *string {
	if gatewayIcon == nil {
		return nil
	} else {
		return unibee.String(strings.Join(*gatewayIcon, "|"))
	}
}

func SetupGateway(ctx context.Context, merchantId uint64, gatewayName string, gatewayKey string, gatewaySecret string, subGateway string, displayName *string, gatewayIcon *[]string, sort *int64, currencyExchange []*detail.GatewayCurrencyExchange) *entity.MerchantGateway {
	utility.Assert(len(gatewayName) > 0, "gatewayName invalid")
	gatewayInfo := api.GetGatewayWebhookServiceProviderByGatewayName(ctx, gatewayName).GatewayInfo(ctx)
	utility.Assert(gatewayInfo != nil, "gateway not ready")
	if len(gatewayKey) > 0 || len(gatewaySecret) > 0 {
		_, _, err := api.GetGatewayWebhookServiceProviderByGatewayName(ctx, gatewayName).GatewayTest(ctx, gatewayKey, gatewaySecret, subGateway)
		utility.AssertError(err, "gateway test error, key or secret invalid")
	}
	utility.Assert(gatewayName != "wire_transfer", "gateway should not wire transfer type")
	var one *entity.MerchantGateway
	err := dao.MerchantGateway.Ctx(ctx).
		Where(dao.MerchantGateway.Columns().MerchantId, merchantId).
		Where(dao.MerchantGateway.Columns().GatewayName, gatewayName).
		Where(dao.MerchantGateway.Columns().IsDeleted, 0).
		OmitNil().
		Scan(&one)
	utility.AssertError(err, "system error")
	utility.Assert(one == nil, "same gateway exist")

	var name = ""
	if displayName != nil {
		name = *displayName
	}
	var logo = ""
	if gatewayIcon != nil {
		logo = unibee.StringValue(JoinGatewayIcon(gatewayIcon))
	}
	var gatewaySort int64 = 0
	if sort != nil {
		gatewaySort = *sort
	} else {
		gatewaySort = gatewayInfo.Sort
	}
	one = &entity.MerchantGateway{
		MerchantId:    merchantId,
		GatewayName:   gatewayName,
		GatewayKey:    gatewayKey,
		GatewaySecret: gatewaySecret,
		SubGateway:    subGateway,
		EnumKey:       gatewaySort,
		GatewayType:   gatewayInfo.GatewayType,
		Name:          name,
		Logo:          logo,
		Host:          gatewayInfo.Host,
		Custom:        utility.MarshalToJsonString(currencyExchange),
	}
	result, err := dao.MerchantGateway.Ctx(ctx).Data(one).OmitNil().Insert(one)
	utility.AssertError(err, "system error")
	id, _ := result.LastInsertId()
	one.Id = uint64(id)

	if len(gatewayKey) > 0 || len(gatewaySecret) > 0 {
		gatewayWebhook.CheckAndSetupGatewayWebhooks(ctx, one.Id)
	}

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

func UpdateGatewaySort(ctx context.Context, merchantId uint64, gatewayId uint64, sort int64) {
	utility.Assert(gatewayId > 0, "gatewayId invalid")
	one := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(one != nil, "gateway not found")
	utility.Assert(one.MerchantId == merchantId, "merchant not match")
	_, err := dao.MerchantGateway.Ctx(ctx).Data(g.Map{
		dao.MerchantGateway.Columns().EnumKey:   sort,
		dao.MerchantGateway.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantGateway.Columns().Id, one.Id).OmitNil().Update()
	utility.AssertError(err, "system error")
}

func EditGateway(ctx context.Context, merchantId uint64, gatewayId uint64, targetGatewayKey *string, targetGatewaySecret *string, targetSubGateway *string, displayName *string, gatewayIcon *[]string, sort *int64, currencyExchange []*detail.GatewayCurrencyExchange) *entity.MerchantGateway {
	utility.Assert(gatewayId > 0, "gatewayId invalid")
	one := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(one != nil, "gateway not found")
	utility.Assert(one.MerchantId == merchantId, "merchant not match")

	if targetGatewayKey != nil || targetGatewaySecret != nil {
		utility.Assert(one.GatewayType != consts.GatewayTypeWireTransfer, "gateway should not wire transfer type")
		gatewayKey := one.GatewayKey
		gatewaySecret := one.GatewaySecret
		subGateway := one.SubGateway
		if targetGatewayKey != nil {
			gatewayKey = *targetGatewayKey
		}
		if targetGatewaySecret != nil {
			gatewaySecret = *targetGatewaySecret
		}
		if targetSubGateway != nil {
			subGateway = *targetSubGateway
		}
		_, _, err := api.GetGatewayServiceProvider(ctx, gatewayId).GatewayTest(ctx, gatewayKey, gatewaySecret, subGateway)
		utility.AssertError(err, "gateway test error, key or secret invalid")
		_, err = dao.MerchantGateway.Ctx(ctx).Data(g.Map{
			dao.MerchantGateway.Columns().GatewaySecret: gatewaySecret,
			dao.MerchantGateway.Columns().GatewayKey:    gatewayKey,
		}).Where(dao.MerchantGateway.Columns().Id, one.Id).OmitNil().Update()
		utility.AssertError(err, "system error")
		gatewayWebhook.CheckAndSetupGatewayWebhooks(ctx, one.Id)
	}

	_, err := dao.MerchantGateway.Ctx(ctx).Data(g.Map{
		dao.MerchantGateway.Columns().Logo:       JoinGatewayIcon(gatewayIcon),
		dao.MerchantGateway.Columns().Name:       displayName,
		dao.MerchantGateway.Columns().EnumKey:    sort,
		dao.MerchantGateway.Columns().SubGateway: targetSubGateway,
		dao.MerchantGateway.Columns().Custom:     utility.MarshalMetadataToJsonString(currencyExchange),
		dao.MerchantGateway.Columns().GmtModify:  gtime.Now(),
	}).Where(dao.MerchantGateway.Columns().Id, one.Id).OmitNil().Update()
	utility.AssertError(err, "system error")

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

func ArchiveGateway(ctx context.Context, merchantId uint64, gatewayId uint64) *entity.MerchantGateway {
	utility.Assert(gatewayId > 0, "gatewayId invalid")
	one := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(one.GatewayType != consts.GatewayTypeWireTransfer, "invalid gateway, wire transfer not supported")
	utility.Assert(one != nil, "gateway not found")
	utility.Assert(one.MerchantId == merchantId, "merchant not match")
	_, err := dao.MerchantGateway.Ctx(ctx).Data(g.Map{
		dao.MerchantGateway.Columns().IsDeleted: gtime.Now().Timestamp(),
	}).Where(dao.MerchantGateway.Columns().Id, one.Id).OmitNil().Update()
	utility.AssertError(err, "system error")
	one = query.GetGatewayById(ctx, gatewayId)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Gateway(%v-%s)", one.Id, one.GatewayName),
		Content:        "Archive",
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
	gatewaySimplify := detail.ConvertGatewayDetail(ctx, gateway)
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

func GetMerchantAvailableGatewaysByCountryCode(ctx context.Context, merchantId uint64, countryCode string) []*detail.Gateway {
	var availableGateways []*detail.Gateway
	gateways := query.GetMerchantGatewayList(ctx, merchantId, unibee.Bool(false))
	for _, one := range gateways {
		if IsGatewaySupportCountryCode(ctx, one, countryCode) {
			availableGateways = append(availableGateways, detail.ConvertGatewayDetail(ctx, one))
		}
	}
	return availableGateways
}

type WireTransferSetupReq struct {
	GatewayId     uint64              `json:"gatewayId"  dc:"The id of payment gateway" v:"required"`
	MerchantId    uint64              `json:"merchantId"   dc:"The merchantId of wire transfer" v:"required" `
	Currency      string              `json:"currency"   dc:"The currency of wire transfer " v:"required" `
	MinimumAmount int64               `json:"minimumAmount"   dc:"The minimum amount of wire transfer" v:"required" `
	Bank          *detail.GatewayBank `json:"bank"   dc:"The receiving bank of wire transfer " v:"required" `
	DisplayName   *string
	GatewayIcon   *[]string
	Sort          *int64
}
type WireTransferSetupRes struct {
}

func SetupWireTransferGateway(ctx context.Context, req *WireTransferSetupReq) *entity.MerchantGateway {
	gatewayName := "wire_transfer"
	var one *entity.MerchantGateway
	gatewayInfo := api.GetGatewayWebhookServiceProviderByGatewayName(ctx, gatewayName).GatewayInfo(ctx)
	utility.Assert(gatewayInfo != nil, "gateway not ready")
	err := dao.MerchantGateway.Ctx(ctx).
		Where(dao.MerchantGateway.Columns().MerchantId, req.MerchantId).
		Where(dao.MerchantGateway.Columns().GatewayName, gatewayName).
		OmitEmpty().
		Scan(&one)
	utility.AssertError(err, "system error")
	utility.Assert(one == nil, "same gateway exist")
	var name = ""
	if req.DisplayName != nil {
		name = *req.DisplayName
	}
	var logo = ""
	if req.GatewayIcon != nil {
		logo = unibee.StringValue(JoinGatewayIcon(req.GatewayIcon))
	}
	var gatewaySort int64 = 0
	if req.Sort != nil {
		gatewaySort = *req.Sort
	} else {
		gatewaySort = gatewayInfo.Sort
	}
	one = &entity.MerchantGateway{
		MerchantId:    req.MerchantId,
		GatewayName:   gatewayName,
		Currency:      strings.ToUpper(req.Currency),
		MinimumAmount: req.MinimumAmount,
		GatewayType:   consts.GatewayTypeWireTransfer,
		BankData:      utility.MarshalToJsonString(req.Bank),
		Name:          name,
		Logo:          logo,
		EnumKey:       gatewaySort,
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
	utility.Assert(one.GatewayType == consts.GatewayTypeWireTransfer, "gateway should be wire transfer type")
	utility.Assert(one.MerchantId == req.MerchantId, "merchant not match")

	_, err := dao.MerchantGateway.Ctx(ctx).Data(g.Map{
		dao.MerchantGateway.Columns().Logo:          JoinGatewayIcon(req.GatewayIcon),
		dao.MerchantGateway.Columns().Name:          req.DisplayName,
		dao.MerchantGateway.Columns().EnumKey:       req.Sort,
		dao.MerchantGateway.Columns().BankData:      utility.MarshalToJsonString(req.Bank),
		dao.MerchantGateway.Columns().Currency:      strings.ToUpper(req.Currency),
		dao.MerchantGateway.Columns().MinimumAmount: req.MinimumAmount,
		dao.MerchantGateway.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.MerchantGateway.Columns().Id, one.Id).OmitNil().Update()
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
