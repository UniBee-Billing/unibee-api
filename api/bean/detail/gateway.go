package detail

import (
	"context"
	"strings"
	"unibee/api/bean"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	gateway2 "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
	"unicode"
)

type GatewaySort struct {
	GatewayName string `json:"gatewayName" description:"Required, The gateway name, stripe|paypal|changelly|unitpay|payssion|cryptadium"`
	Id          uint64 `json:"gatewayId" description:"The gateway id"`
	Sort        int64  `json:"sort" description:"Required, The sort value of payment gateway, should greater than 0, The bigger, the closer to the front"`
}

type GatewayCurrencyExchange struct {
	FromCurrency string  `json:"from_currency" description:"the currency of gateway exchange from"`
	ToCurrency   string  `json:"to_currency" description:"the currency of gateway exchange to"`
	ExchangeRate float64 `json:"exchange_rate"  description:"the exchange rate of gateway, set to 0 if using https://app.exchangerate-api.com/ instead of fixed exchange rate"`
}

type Gateway struct {
	Id                            uint64                         `json:"gatewayId"`
	Name                          string                         `json:"name" description:"The name of gateway"`
	Description                   string                         `json:"description" description:"The description of gateway"`
	GatewayName                   string                         `json:"gatewayName" description:"The gateway name, stripe|paypal|changelly|unitpay|payssion|cryptadium"`
	DisplayName                   string                         `json:"displayName" description:"The gateway display name, used at user portal"`
	GatewayIcons                  []string                       `json:"gatewayIcons"  description:"The gateway display name, used at user portal"`
	GatewayWebsiteLink            string                         `json:"gatewayWebsiteLink" description:"The gateway website link"`
	GatewayWebhookIntegrationLink string                         `json:"gatewayWebhookIntegrationLink" description:"The gateway webhook integration guide link, gateway webhook need setup if not blank"`
	GatewayLogo                   string                         `json:"gatewayLogo"`
	GatewayKey                    string                         `json:"gatewayKey"            description:""`
	GatewaySecret                 string                         `json:"gatewaySecret"            description:""`
	SubGateway                    string                         `json:"subGateway"            description:""`
	GatewayType                   int64                          `json:"gatewayType"           description:"gateway type，1-Bank Card ｜ 2-Crypto | 3 - Wire Transfer"`
	CountryConfig                 map[string]bool                `json:"countryConfig"`
	CreateTime                    int64                          `json:"createTime"            description:"create utc time"` // create utc time
	MinimumAmount                 int64                          `json:"minimumAmount"   description:"The minimum amount of wire transfer" `
	Currency                      string                         `json:"currency"   description:"The currency of wire transfer " `
	Bank                          *GatewayBank                   `json:"bank"   dc:"The receiving bank of wire transfer" `
	WebhookEndpointUrl            string                         `json:"webhookEndpointUrl"   description:"The endpoint url of gateway webhook " `
	WebhookSecret                 string                         `json:"webhookSecret"  dc:"The secret of gateway webhook"`
	Sort                          int64                          `json:"sort"               description:"The sort value of payment gateway, The bigger, the closer to the front"`
	IsSetupFinished               bool                           `json:"IsSetupFinished"  dc:"Whether the gateway finished setup process" `
	CurrencyExchange              []*GatewayCurrencyExchange     `json:"currencyExchange" dc:"The currency exchange for gateway payment, effect at start of payment creation when currency matched"`
	CurrencyExchangeEnabled       bool                           `json:"currencyExchangeEnabled"            description:"whether to enable currency exchange"`
	SubGatewayConfigs             []*_interface.SubGatewayConfig `json:"subGatewayConfigs"  dc:""`
	Archive                       bool                           `json:"archive"  dc:""`
	PublicKeyName                 string                         `json:"publicKeyName"  dc:""`
	PrivateSecretName             string                         `json:"privateSecretName"  dc:""`
	SubGatewayName                string                         `json:"subGatewayName"  dc:""`
	AutoChargeEnabled             bool                           `json:"autoChargeEnabled"  dc:""`
}

type GatewayBank struct {
	AccountHolder string `json:"accountHolder"   dc:"The AccountHolder of wire transfer " v:"required" `
	BIC           string `json:"bic"   dc:"The BIC of wire transfer " v:"required" `
	IBAN          string `json:"iban"   dc:"The IBAN of wire transfer " v:"required" `
	Address       string `json:"address"   dc:"The address of wire transfer " v:"required" `
}

func ConvertGatewayDetail(ctx context.Context, one *entity.MerchantGateway) *Gateway {
	if one == nil {
		return nil
	}
	var countryConfig map[string]bool
	_ = bean.UnmarshalFromJsonString(one.CountryConfig, &countryConfig)
	var bank *GatewayBank
	_ = bean.UnmarshalFromJsonString(one.BankData, &bank)
	var webhookEndpointUrl = ""
	if one.GatewayType != consts.GatewayTypeWireTransfer {
		webhookEndpointUrl = gateway2.GetPaymentWebhookEntranceUrl(one.Id)
	}

	var gatewayIcons = make([]string, 0)
	gatewayInfo := api.GetGatewayServiceProvider(ctx, one.Id).GatewayInfo(ctx)
	if gatewayInfo != nil {
		gatewayIcons = gatewayInfo.GatewayIcons
	}
	if len(one.Logo) > 0 && one.Logo != "http://unibee.top/files/invoice/changelly.png" && one.Logo != "http://unibee.top/files/invoice/stripe.png" && one.Logo != "https://www.paypalobjects.com/webstatic/icon/favicon.ico" {
		gatewayIcons = strings.Split(one.Logo, "|")
	}

	var displayName = ""
	if gatewayInfo != nil {
		displayName = gatewayInfo.DisplayName
	}
	if len(one.Name) > 0 && one.Name != "stripe" && one.Name != "changelly" {
		displayName = one.Name
	}
	name := one.Name
	if gatewayInfo != nil {
		name = gatewayInfo.Name
	}
	description := one.Description
	if gatewayInfo != nil {
		description = gatewayInfo.Description
	}
	gatewayLogo := ""
	if gatewayInfo != nil {
		gatewayLogo = gatewayInfo.GatewayLogo
	}
	gatewayWebsiteLink := ""
	if gatewayInfo != nil {
		gatewayWebsiteLink = gatewayInfo.GatewayWebsiteLink
	}
	gatewayWebhookIntegrationLink := ""
	if gatewayInfo != nil {
		gatewayWebhookIntegrationLink = gatewayInfo.GatewayWebhookIntegrationLink
	}
	isSetupFinished := true
	if one.GatewayType != consts.GatewayTypeWireTransfer {
		if len(one.GatewayKey) == 0 {
			isSetupFinished = false
		}
		if gatewayInfo != nil {
			if len(gatewayInfo.GatewayWebhookIntegrationLink) > 0 {
				if len(one.WebhookSecret) == 0 {
					isSetupFinished = false
				}
			}
		}
	}
	currencyExchangeEnabled := false
	var publicKeyName = "Public Key"
	var privateSecretName = "Private Key"
	var subGatewayName = ""
	var autoChargeEnabled = false
	if gatewayInfo != nil {
		currencyExchangeEnabled = gatewayInfo.CurrencyExchangeEnabled
		if len(gatewayInfo.PublicKeyName) > 0 {
			publicKeyName = gatewayInfo.PublicKeyName
		}
		if len(gatewayInfo.PrivateSecretName) > 0 {
			privateSecretName = gatewayInfo.PrivateSecretName
		}
		if len(gatewayInfo.SubGatewayName) > 0 {
			subGatewayName = gatewayInfo.SubGatewayName
		}
		autoChargeEnabled = gatewayInfo.AutoChargeEnabled
	}
	var currencyExchangeList = make([]*GatewayCurrencyExchange, 0)
	_ = utility.UnmarshalFromJsonString(one.Custom, &currencyExchangeList)

	//subGatewayConfigs := make([]*_interface.SubGatewayConfig, 0)
	//if gatewayInfo != nil {
	//	subGatewayConfigs = gatewayInfo.SubGatewayConfigs
	//}
	if one.EnumKey <= 0 && gatewayInfo != nil {
		one.EnumKey = gatewayInfo.Sort
	}

	return &Gateway{
		Id:                            one.Id,
		Name:                          name,
		Description:                   description,
		GatewayLogo:                   gatewayLogo,
		GatewayWebsiteLink:            gatewayWebsiteLink,
		GatewayWebhookIntegrationLink: gatewayWebhookIntegrationLink,
		GatewayIcons:                  gatewayIcons,
		GatewayName:                   one.GatewayName,
		DisplayName:                   displayName,
		GatewayType:                   one.GatewayType,
		CountryConfig:                 countryConfig,
		CreateTime:                    one.CreateTime,
		Currency:                      one.Currency,
		MinimumAmount:                 one.MinimumAmount,
		Bank:                          bank,
		WebhookEndpointUrl:            webhookEndpointUrl,
		GatewayKey:                    utility.HideStar(one.GatewayKey),
		GatewaySecret:                 utility.HideStar(one.GatewaySecret),
		WebhookSecret:                 utility.HideStar(one.WebhookSecret),
		SubGateway:                    one.SubGateway,
		Sort:                          one.EnumKey,
		IsSetupFinished:               isSetupFinished,
		CurrencyExchange:              currencyExchangeList,
		CurrencyExchangeEnabled:       currencyExchangeEnabled,
		Archive:                       one.IsDeleted != 0,
		PublicKeyName:                 publicKeyName,
		PrivateSecretName:             privateSecretName,
		SubGatewayName:                subGatewayName,
		AutoChargeEnabled:             autoChargeEnabled,
		//SubGatewayConfigs:             subGatewayConfigs,
	}
}

func toUpperFirst(s string, target string) string {
	if len(target) > 0 {
		return target
	}
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func ConvertGatewayList(ctx context.Context, ones []*entity.MerchantGateway) (list []*Gateway) {
	if len(ones) == 0 {
		return make([]*Gateway, 0)
	}
	for _, one := range ones {
		list = append(list, ConvertGatewayDetail(ctx, one))
	}
	return list
}
