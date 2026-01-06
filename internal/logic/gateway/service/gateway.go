package service

import (
	"context"
	"fmt"
	"strings"
	"unibee/api/bean/detail"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/api"
	gatewayWebhook "unibee/internal/logic/gateway/webhook"
	"unibee/internal/logic/merchant_config"
	"unibee/internal/logic/multi_currencies/currency_exchange"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func SetupGateway(ctx context.Context, merchantId uint64, gatewayName string, gatewayKey string, gatewaySecret string, subGateway string, paymentTypes []string, displayName *string, gatewayIcon *[]string, sort *int64, currencyExchange []*detail.GatewayCurrencyExchange, metadata map[string]interface{}) *entity.MerchantGateway {
	utility.Assert(len(gatewayName) > 0, "gatewayName invalid")
	gatewayInfo := api.GetGatewayWebhookServiceProviderByGatewayName(ctx, gatewayName).GatewayInfo(ctx)
	utility.Assert(gatewayInfo != nil, "gateway not ready")
	utility.Assert(len(gatewayKey) > 0, "publicKey invalid")
	if len(gatewayKey) > 0 || len(gatewaySecret) > 0 {
		var gatewayPaymentTypes = make([]*_interface.GatewayPaymentType, 0)
		for _, paymentTypeStr := range paymentTypes {
			for _, infoPaymentType := range gatewayInfo.GatewayPaymentTypes {
				if paymentTypeStr == infoPaymentType.PaymentType {
					gatewayPaymentTypes = append(gatewayPaymentTypes, infoPaymentType)
				}
			}
		}
		_, _, err := api.GetGatewayWebhookServiceProviderByGatewayName(ctx, gatewayName).GatewayTest(ctx, &_interface.GatewayTestReq{
			Key:                 gatewayKey,
			Secret:              gatewaySecret,
			SubGateway:          subGateway,
			GatewayPaymentTypes: gatewayPaymentTypes,
		})
		utility.AssertError(err, "gateway test error, key or secret invalid")
	}
	utility.Assert(gatewayName != "wire_transfer", "gateway should not wire transfer type")
	if gatewayInfo.GatewayType == consts.GatewayTypeCrypto {
		exchangeApiKeyConfig := merchant_config.GetMerchantConfig(ctx, merchantId, currency_exchange.FiatExchangeApiKey)
		if config.GetConfigInstance().Mode != "cloud" {
			utility.Assert(exchangeApiKeyConfig != nil && len(exchangeApiKeyConfig.ConfigValue) > 0, "ExchangeApi Need Setup First For Crypto Gateway")
		}
	}

	//var one *entity.MerchantGateway
	//err := dao.MerchantGateway.Ctx(ctx).
	//	Where(dao.MerchantGateway.Columns().MerchantId, merchantId).
	//	Where(dao.MerchantGateway.Columns().GatewayName, gatewayName).
	//	Where(dao.MerchantGateway.Columns().IsDeleted, 0).
	//	OmitNil().
	//	Scan(&one)
	//utility.AssertError(err, "system error")
	//utility.Assert(one == nil, "same gateway exist")

	minData := GetMinData(ctx, merchantId, gatewayName)

	var name = ""
	if displayName != nil {
		name = *displayName
	}
	var logo = ""
	if gatewayIcon != nil {
		logo = unibee.StringValue(utility.ArrayPointJoinToStringPoint(gatewayIcon))
	}
	var gatewaySort int64 = 0
	if sort != nil {
		gatewaySort = *sort
	} else {
		gatewaySort = gatewayInfo.Sort
	}
	one := &entity.MerchantGateway{
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
		BrandData:     unibee.StringValue(utility.JoinToStringPoint(paymentTypes)),
		IsDeleted:     int(minData),
		MetaData:      utility.MarshalToJsonString(metadata),
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

func GetMinData(ctx context.Context, merchantId uint64, gatewayName string) float64 {
	minData, err := dao.MerchantGateway.Ctx(ctx).
		Where(dao.MerchantGateway.Columns().MerchantId, merchantId).
		Where(dao.MerchantGateway.Columns().GatewayName, gatewayName).
		Min(dao.MerchantGateway.Columns().IsDeleted)
	utility.AssertError(err, "system error")
	if minData > 0 {
		minData = 0
	} else if minData == 0 {
		if query.GetDefaultGatewayByGatewayName(ctx, merchantId, gatewayName) == nil {
			minData = 0
		} else {
			minData = -1
		}
	} else {
		minData = minData - 1
	}
	return minData
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

func EditGateway(ctx context.Context, merchantId uint64, gatewayId uint64, targetGatewayKey *string, targetGatewaySecret *string, targetSubGateway *string, paymentTypes []string, displayName *string, gatewayIcon *[]string, sort *int64, currencyExchange []*detail.GatewayCurrencyExchange, metadata *map[string]interface{}) *entity.MerchantGateway {
	utility.Assert(gatewayId > 0, "gatewayId invalid")
	one := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(one != nil, "gateway not found")
	utility.Assert(one.MerchantId == merchantId, "merchant not match")
	gatewayInfo := api.GetGatewayWebhookServiceProviderByGatewayName(ctx, one.GatewayName).GatewayInfo(ctx)
	utility.Assert(gatewayInfo != nil, "gateway not ready")

	if targetGatewayKey != nil || targetGatewaySecret != nil {
		utility.Assert(one.GatewayType != consts.GatewayTypeWireTransfer, "gateway should not wire transfer type")
		utility.Assert(one.IsDeleted <= 0, "gateway already archived")
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
		var gatewayPaymentTypes = make([]*_interface.GatewayPaymentType, 0)
		for _, paymentTypeStr := range paymentTypes {
			for _, infoPaymentType := range gatewayInfo.GatewayPaymentTypes {
				if paymentTypeStr == infoPaymentType.PaymentType {
					gatewayPaymentTypes = append(gatewayPaymentTypes, infoPaymentType)
				}
			}
		}
		_, _, err := api.GetGatewayServiceProvider(ctx, gatewayId).GatewayTest(ctx, &_interface.GatewayTestReq{
			OldKey:              one.GatewayKey,
			OldSecret:           one.GatewaySecret,
			Key:                 gatewayKey,
			Secret:              gatewaySecret,
			SubGateway:          subGateway,
			GatewayPaymentTypes: gatewayPaymentTypes,
		})
		utility.AssertError(err, "gateway test error, key or secret invalid")
		_, err = dao.MerchantGateway.Ctx(ctx).Data(g.Map{
			dao.MerchantGateway.Columns().GatewaySecret: gatewaySecret,
			dao.MerchantGateway.Columns().GatewayKey:    gatewayKey,
		}).Where(dao.MerchantGateway.Columns().Id, one.Id).OmitNil().Update()
		utility.AssertError(err, "system error")
		gatewayWebhook.CheckAndSetupGatewayWebhooks(ctx, one.Id)
	}

	_, err := dao.MerchantGateway.Ctx(ctx).Data(g.Map{
		dao.MerchantGateway.Columns().Logo:       utility.ArrayPointJoinToStringPoint(gatewayIcon),
		dao.MerchantGateway.Columns().Name:       displayName,
		dao.MerchantGateway.Columns().BrandData:  utility.JoinToStringPoint(paymentTypes),
		dao.MerchantGateway.Columns().EnumKey:    sort,
		dao.MerchantGateway.Columns().SubGateway: targetSubGateway,
		dao.MerchantGateway.Columns().Custom:     utility.MarshalMetadataToJsonString(currencyExchange),
		dao.MerchantGateway.Columns().MetaData:   utility.MarshalMetadataToJsonString(utility.MergeMetadata(one.MetaData, metadata)),
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
	utility.Assert(one.GatewayType != consts.GatewayTypeWireTransfer, "Invalid gateway: wire transfer not supported")
	utility.Assert(one != nil, "gateway not found")
	utility.Assert(one.MerchantId == merchantId, "merchant not match")
	list := make([]*entity.MerchantGateway, 0)
	err := dao.MerchantGateway.Ctx(ctx).
		Where(dao.MerchantGateway.Columns().MerchantId, merchantId).
		Where(dao.MerchantGateway.Columns().GatewayName, one.GatewayName).
		WhereLTE(dao.MerchantGateway.Columns().IsDeleted, 0).
		Scan(&list)
	utility.AssertError(err, "system error")
	utility.Assert(len(list) > 0, "no valid gateway found")
	utility.Assert(len(list) > 1, "Archiving failed: only one valid gateway available")
	var nextDefaultOne *entity.MerchantGateway
	if one.IsDeleted == 0 {
		for _, o := range list {
			if o.Id != one.Id && o.IsDeleted < 0 {
				nextDefaultOne = o
				break
			}
		}
		utility.Assert(nextDefaultOne != nil, "Archiving failed: only one valid gateway available")
	}

	_, err = dao.MerchantGateway.Ctx(ctx).Data(g.Map{
		dao.MerchantGateway.Columns().IsDeleted: gtime.Now().Timestamp(),
	}).Where(dao.MerchantGateway.Columns().Id, one.Id).OmitNil().Update()
	utility.AssertError(err, "system error")
	one = query.GetGatewayById(ctx, gatewayId)

	if nextDefaultOne != nil {
		_, err = dao.MerchantGateway.Ctx(ctx).Data(g.Map{
			dao.MerchantGateway.Columns().IsDeleted: 0,
		}).Where(dao.MerchantGateway.Columns().Id, nextDefaultOne.Id).OmitNil().Update()
		utility.AssertError(err, "system error")
	}
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

func RestoreGateway(ctx context.Context, merchantId uint64, gatewayId uint64) *entity.MerchantGateway {
	utility.Assert(gatewayId > 0, "gatewayId invalid")
	one := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(one.GatewayType != consts.GatewayTypeWireTransfer, "invalid gateway, wire transfer not supported")
	utility.Assert(one != nil, "gateway not found")
	utility.Assert(one.MerchantId == merchantId, "merchant not match")
	if one.IsDeleted > 0 {
		minData := GetMinData(ctx, merchantId, one.GatewayName)
		_, err := dao.MerchantGateway.Ctx(ctx).Data(g.Map{
			dao.MerchantGateway.Columns().IsDeleted: minData,
		}).Where(dao.MerchantGateway.Columns().Id, one.Id).OmitNil().Update()
		utility.AssertError(err, "system error")
		one = query.GetGatewayById(ctx, gatewayId)
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     one.MerchantId,
			Target:         fmt.Sprintf("Gateway(%v-%s)", one.Id, one.GatewayName),
			Content:        "Restore",
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)
	}

	return one
}

func SetDefaultGateway(ctx context.Context, merchantId uint64, gatewayId uint64) *entity.MerchantGateway {
	utility.Assert(gatewayId > 0, "gatewayId invalid")
	one := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(one.GatewayType != consts.GatewayTypeWireTransfer, "invalid gateway, wire transfer not supported")
	utility.Assert(one != nil, "gateway not found")
	utility.Assert(one.MerchantId == merchantId, "merchant not match")
	if one.IsDeleted != 0 {
		oldDefaultOne := query.GetDefaultGatewayByGatewayName(ctx, merchantId, one.GatewayName)
		var err error
		if oldDefaultOne != nil {
			minData := GetMinData(ctx, merchantId, oldDefaultOne.GatewayName)
			err = dao.MerchantGateway.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
				result, err := transaction.Update(dao.MerchantGateway.Table(), g.Map{dao.MerchantGateway.Columns().IsDeleted: minData},
					g.Map{dao.MerchantGateway.Columns().Id: oldDefaultOne.Id})
				if err != nil || result == nil {
					return err
				}
				result, err = transaction.Update(dao.MerchantGateway.Table(), g.Map{dao.MerchantGateway.Columns().IsDeleted: 0},
					g.Map{dao.MerchantGateway.Columns().Id: one.Id})
				if err != nil || result == nil {
					return err
				}
				return nil
			})
			utility.AssertError(err, "system error")
		} else {
			_, err = dao.MerchantGateway.Ctx(ctx).Data(g.Map{
				dao.MerchantGateway.Columns().IsDeleted: 0,
			}).Where(dao.MerchantGateway.Columns().Id, one.Id).OmitNil().Update()
			utility.AssertError(err, "system error")
		}
		one = query.GetGatewayById(ctx, gatewayId)
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     one.MerchantId,
			Target:         fmt.Sprintf("Gateway(%v-%s)", one.Id, one.GatewayName),
			Content:        "SetDefault",
			UserId:         0,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)
	}

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
	query.GetDefaultGatewayByGatewayName(ctx, req.MerchantId, gatewayName)
	utility.Assert(one == nil, "same gateway exist")
	var name = ""
	if req.DisplayName != nil {
		name = *req.DisplayName
	}
	var logo = ""
	if req.GatewayIcon != nil {
		logo = unibee.StringValue(utility.ArrayPointJoinToStringPoint(req.GatewayIcon))
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
		dao.MerchantGateway.Columns().Logo:          utility.ArrayPointJoinToStringPoint(req.GatewayIcon),
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
