package merchant

import (
	"context"
	"sort"
	"unibee/api/bean/detail"
	"unibee/api/merchant/gateway"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/query"
)

func (c *ControllerGateway) SetupList(ctx context.Context, req *gateway.SetupListReq) (res *gateway.SetupListRes, err error) {
	var list = make([]*detail.Gateway, 0)
	for _, gatewayName := range api.ExportGatewaySetupListKeys {
		if info, exists := api.ExportGatewaySetupList[gatewayName]; exists {
			one := query.GetGatewayByGatewayName(ctx, _interface.GetMerchantId(ctx), gatewayName)
			if one != nil && one.IsDeleted == 0 {
				gatewayDetail := detail.ConvertGatewayDetail(ctx, one)
				gatewayDetail.SubGatewayConfigs = info.SubGatewayConfigs
				list = append(list, gatewayDetail)
			} else {
				list = append(list, &detail.Gateway{
					Id:                            0,
					Name:                          info.Name,
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
					Sort:                          0,
					IsSetupFinished:               false,
					CurrencyExchangeEnabled:       info.CurrencyExchangeEnabled,
					SubGatewayConfigs:             info.SubGatewayConfigs,
				})
			}
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Sort < list[j].Sort
	})
	return &gateway.SetupListRes{Gateways: list}, nil
}
