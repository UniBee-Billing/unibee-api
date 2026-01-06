package merchant

import (
	"context"
	"fmt"
	"sort"
	"unibee/api/bean/detail"
	"unibee/api/merchant/gateway"
	"unibee/internal/cmd/config"
	_interface2 "unibee/internal/interface"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/merchant_config"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerGateway) SetupList(ctx context.Context, req *gateway.SetupListReq) (res *gateway.SetupListRes, err error) {
	var list = make([]*detail.Gateway, 0)
	sortConfig := merchant_config.GetMerchantConfig(ctx, _interface.GetMerchantId(ctx), _interface2.KEY_MERCHANT_GATEWAY_SORT)
	var sortMap = make(map[string]int64)
	if sortConfig != nil {
		_ = utility.UnmarshalFromJsonString(sortConfig.ConfigValue, &sortMap)
	}
	gatewayList := query.GetMerchantGatewaySetupList(ctx, _interface.GetMerchantId(ctx))
	gatewayMap := make(map[string][]*entity.MerchantGateway)
	for _, one := range gatewayList {
		if slice, ok := gatewayMap[one.GatewayName]; ok {
			gatewayMap[one.GatewayName] = append(slice, one)
		} else {
			gatewayMap[one.GatewayName] = []*entity.MerchantGateway{one}
		}
	}
	for _, gatewayName := range api.ExportGatewaySetupListKeys() {
		if info, exists := api.ExportGatewaySetupList[gatewayName]; exists {
			if config.GetConfigInstance().IsProd() {
				if info.IsStaging {
					continue
				}
			}
			var needSetupValidOne = true
			targetGatewayList := gatewayMap[gatewayName]
			if targetGatewayList != nil && len(targetGatewayList) > 0 {
				//needSetupValidOne = true
				//} else {
				for _, one := range targetGatewayList {
					gatewayDetail := detail.ConvertGatewayDetail(ctx, one)
					gatewayDetail.SetupGatewayPaymentTypes = info.GatewayPaymentTypes
					gatewayDetail.Name = fmt.Sprintf("[ID: %d] %s", gatewayDetail.Id, gatewayDetail.Name)
					if gatewayDetail.Archive {
						gatewayDetail.Name = fmt.Sprintf("%s ARV", gatewayDetail.Name)
						gatewayDetail.Sort = 99999999
						// archived gateway sort to end
					}
					if one.IsDeleted == 0 {
						needSetupValidOne = false
					}
					list = append(list, gatewayDetail)
				}
			}

			//one := query.GetDefaultGatewayByGatewayName(ctx, _interface.GetMerchantId(ctx), gatewayName)
			//if one != nil && one.IsDeleted == 0 {
			//	gatewayDetail := detail.ConvertGatewayDetail(ctx, one)
			//	gatewayDetail.SetupGatewayPaymentTypes = info.GatewayPaymentTypes
			//	list = append(list, gatewayDetail)
			//} else {
			if needSetupValidOne {
				gatewaySort := info.Sort
				if _, ok := sortMap[gatewayName]; ok {
					gatewaySort = sortMap[gatewayName]
				}
				var publicKeyName = "Public Key"
				var privateSecretName = "Private Key"
				var subGatewayName = ""

				if len(info.PublicKeyName) > 0 {
					publicKeyName = info.PublicKeyName
				}
				if len(info.PrivateSecretName) > 0 {
					privateSecretName = info.PrivateSecretName
				}
				if len(info.SubGatewayName) > 0 {
					subGatewayName = info.SubGatewayName
				}
				name := info.Name
				if info.IsStaging {
					name = fmt.Sprintf("%s Beta", info.Name)
				}
				list = append(list, &detail.Gateway{
					Id:                            0,
					Name:                          name,
					Description:                   info.Description,
					GatewayName:                   gatewayName,
					DisplayName:                   info.DisplayName,
					GatewayIcons:                  info.GatewayIcons,
					GatewayWebsiteLink:            info.GatewayWebsiteLink,
					GatewayWebhookIntegrationLink: info.GatewayWebhookIntegrationLink,
					GatewayLogo:                   info.GatewayLogo,
					GatewayKey:                    "",
					GatewayType:                   info.GatewayType,
					CountryConfig:                 nil,
					CreateTime:                    0,
					MinimumAmount:                 0,
					Currency:                      "",
					Bank:                          nil,
					WebhookEndpointUrl:            "",
					WebhookSecret:                 "",
					Sort:                          gatewaySort,
					IsSetupFinished:               false,
					Archive:                       false,
					CurrencyExchangeEnabled:       info.CurrencyExchangeEnabled,
					SetupGatewayPaymentTypes:      info.GatewayPaymentTypes,
					GatewayPaymentTypes:           make([]*_interface2.GatewayPaymentType, 0),
					PublicKeyName:                 publicKeyName,
					PrivateSecretName:             privateSecretName,
					SubGatewayName:                subGatewayName,
					AutoChargeEnabled:             info.AutoChargeEnabled,
				})
			}
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Sort < list[j].Sort
	})
	return &gateway.SetupListRes{Gateways: list}, nil
}
