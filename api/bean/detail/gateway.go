package detail

import (
	"context"
	"fmt"
	"strings"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	gateway2 "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
	"unicode"

	"github.com/gogf/gf/v2/encoding/gjson"
)

type GatewaySort struct {
	GatewayName string `json:"gatewayName" description:"Required, The gateway name, stripe|paypal|changelly|unitpay|payssion|cryptadium"`
	Id          uint64 `json:"gatewayId" description:"The gateway id"`
	Sort        int64  `json:"sort" description:"Required, The sort value of payment gateway, should greater than 0, The higher the value, the lower the ranking"`
}

type GatewayCurrencyExchange struct {
	FromCurrency string  `json:"from_currency" description:"the currency of gateway exchange from"`
	ToCurrency   string  `json:"to_currency" description:"the currency of gateway exchange to"`
	ExchangeRate float64 `json:"exchange_rate"  description:"the exchange rate of gateway, set to 0 if using https://app.exchangerate-api.com/ instead of fixed exchange rate"`
}

type Gateway struct {
	Id                            uint64                           `json:"gatewayId"`
	Name                          string                           `json:"name" description:"The name of gateway"`
	Description                   string                           `json:"description" description:"The description of gateway"`
	GatewayName                   string                           `json:"gatewayName" description:"The gateway name, stripe|paypal|changelly|unitpay|payssion|cryptadium"`
	DisplayName                   string                           `json:"displayName" description:"The gateway display name, used at user portal"`
	GatewayIcons                  []string                         `json:"gatewayIcons"  description:"The gateway display name, used at user portal"`
	GatewayWebsiteLink            string                           `json:"gatewayWebsiteLink" description:"The gateway website link"`
	GatewayWebhookIntegrationLink string                           `json:"gatewayWebhookIntegrationLink" description:"The gateway webhook integration guide link, gateway webhook need setup if not blank"`
	GatewayLogo                   string                           `json:"gatewayLogo"`
	GatewayKey                    string                           `json:"gatewayKey"            description:""`
	GatewaySecret                 string                           `json:"gatewaySecret"            description:""`
	SubGateway                    string                           `json:"subGateway"            description:""`
	GatewayType                   int64                            `json:"gatewayType"           description:"gateway type，1-Bank Card ｜ 2-Crypto | 3 - Wire Transfer"`
	CountryConfig                 map[string]bool                  `json:"countryConfig"`
	CreateTime                    int64                            `json:"createTime"            description:"create utc time"` // create utc time
	MinimumAmount                 int64                            `json:"minimumAmount"   description:"The minimum amount of wire transfer" `
	Currency                      string                           `json:"currency"   description:"The currency of wire transfer " `
	Bank                          *GatewayBank                     `json:"bank"   dc:"The receiving bank of wire transfer" `
	WebhookEndpointUrl            string                           `json:"webhookEndpointUrl"   description:"The endpoint url of gateway webhook " `
	WebhookSecret                 string                           `json:"webhookSecret"  dc:"The secret of gateway webhook"`
	Sort                          int64                            `json:"sort"               description:"The sort value of payment gateway, The higher the value, the lower the ranking"`
	IsSetupFinished               bool                             `json:"IsSetupFinished"  dc:"Whether the gateway finished setup process" `
	CurrencyExchange              []*GatewayCurrencyExchange       `json:"currencyExchange" dc:"The currency exchange for gateway payment, effect at start of payment creation when currency matched"`
	CurrencyExchangeEnabled       bool                             `json:"currencyExchangeEnabled"            description:"whether to enable currency exchange"`
	GatewayPaymentTypes           []*_interface.GatewayPaymentType `json:"gatewayPaymentTypes" dc:"gatewayPaymentTypes"`
	SetupGatewayPaymentTypes      []*_interface.GatewayPaymentType `json:"setupGatewayPaymentTypes"  dc:"The total list of gateway payment types, used for setup"`
	Archive                       bool                             `json:"archive"  dc:""`
	PublicKeyName                 string                           `json:"publicKeyName"  dc:""`
	PrivateSecretName             string                           `json:"privateSecretName"  dc:""`
	SubGatewayName                string                           `json:"subGatewayName"  dc:""`
	AutoChargeEnabled             bool                             `json:"autoChargeEnabled"  dc:""`
	IsDefault                     bool                             `json:"isDefault"  dc:""`
	CompanyIssuer                 *GatewayCompanyIssuer            `json:"companyIssuer" dc:""`
	Metadata                      map[string]interface{}           `json:"metadata"                  description:""`
}

type GatewayBank struct {
	AccountHolder    string `json:"accountHolder"   dc:"The AccountHolder of wire transfer " v:"required" `
	AccountNumber    string `json:"accountNumber,omitempty"  dc:"The Account Number"`
	SwiftCode        string `json:"swiftCode,omitempty"  dc:"The Swift Code"`
	BankName         string `json:"bankName,omitempty"  dc:"The Bank Name"`
	BSBCode          string `json:"bsbCode,omitempty"  dc:"The BSB Code"`
	BIC              string `json:"bic"   dc:"The BIC of wire transfer " `
	IBAN             string `json:"iban"   dc:"The IBAN of wire transfer " `
	ABARoutingNumber string `json:"ABARoutingNumber"   dc:"The ABARoutingNumber of wire transfer " `
	CNAPS            string `json:"CNAPS"   dc:"The CNAPS of wire transfer " `
	Address          string `json:"address"   dc:"The address of wire transfer " v:"required" `
	Remarks          string `json:"Remarks"   dc:"The Remarks additional content " `
}

type GatewayCompanyIssuer struct {
	IssueVatNumber   string `json:"issueVatNumber"  dc:""`
	IssueRegNumber   string `json:"issueRegNumber"  dc:""`
	IssueCompanyName string `json:"issueCompanyName"  dc:""`
	IssueAddress     string `json:"issueAddress"  dc:""`
	IssueLogo        string `json:"issueLogo"  dc:""`
}

func ConvertGatewayDetail(ctx context.Context, one *entity.MerchantGateway) *Gateway {
	if one == nil {
		return nil
	}
	var metadata = make(map[string]interface{})
	if len(one.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("SimplifyPlan Unmarshal Metadata error:%s", err.Error())
		}
	}
	var countryConfig map[string]bool
	_ = utility.UnmarshalFromJsonString(one.CountryConfig, &countryConfig)
	var bank *GatewayBank
	_ = utility.UnmarshalFromJsonString(one.BankData, &bank)
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
	var gatewayPaymentTypes = make([]*_interface.GatewayPaymentType, 0)
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
		if len(one.BrandData) > 0 {
			for _, paymentTypeStr := range utility.SplitToArray(one.BrandData) {
				for _, infoPaymentType := range gatewayInfo.GatewayPaymentTypes {
					if paymentTypeStr == infoPaymentType.PaymentType {
						gatewayPaymentTypes = append(gatewayPaymentTypes, infoPaymentType)
					}
				}
			}
		}
	}
	var currencyExchangeList = make([]*GatewayCurrencyExchange, 0)
	_ = utility.UnmarshalFromJsonString(one.Custom, &currencyExchangeList)

	if one.EnumKey <= 0 && gatewayInfo != nil {
		one.EnumKey = gatewayInfo.Sort
	}

	var companyIssuer = &GatewayCompanyIssuer{}
	if v, ok := metadata["IssueVatNumber"]; ok {
		companyIssuer.IssueVatNumber = fmt.Sprintf("%s", v)
	}
	if v, ok := metadata["IssueRegNumber"]; ok {
		companyIssuer.IssueRegNumber = fmt.Sprintf("%s", v)
	}
	if v, ok := metadata["IssueCompanyName"]; ok {
		companyIssuer.IssueCompanyName = fmt.Sprintf("%s", v)
	}
	if v, ok := metadata["IssueAddress"]; ok {
		companyIssuer.IssueAddress = fmt.Sprintf("%s", v)
	}
	if v, ok := metadata["IssueLogo"]; ok {
		companyIssuer.IssueLogo = fmt.Sprintf("%s", v)
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
		Archive:                       one.IsDeleted > 0,
		PublicKeyName:                 publicKeyName,
		PrivateSecretName:             privateSecretName,
		SubGatewayName:                subGatewayName,
		AutoChargeEnabled:             autoChargeEnabled,
		GatewayPaymentTypes:           gatewayPaymentTypes,
		IsDefault:                     one.IsDeleted == 0 && one.Id > 0,
		Metadata:                      metadata,
		CompanyIssuer:                 companyIssuer,
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

func CopyGatewayCompanyIssuer(one *entity.MerchantGateway, targetMetaData map[string]interface{}) {
	if one != nil && targetMetaData != nil {
		var metadata = make(map[string]interface{})
		if len(one.MetaData) > 0 {
			err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
			if err != nil {
				fmt.Printf("CopyGatewayCompanyIssuer Unmarshal Metadata error:%s", err.Error())
			}
		}
		issueCompanyName := ""
		issueAddress := ""
		issueVatNumber := ""
		issueRegNumber := ""
		issueLogo := ""
		if v, ok := metadata["IssueCompanyName"]; ok {
			issueCompanyName = fmt.Sprintf("%s", v)
		}
		if v, ok := metadata["IssueAddress"]; ok {
			issueAddress = fmt.Sprintf("%s", v)
		}
		if v, ok := metadata["IssueVatNumber"]; ok {
			issueVatNumber = fmt.Sprintf("%s", v)
		}
		if v, ok := metadata["IssueRegNumber"]; ok {
			issueRegNumber = fmt.Sprintf("%s", v)
		}
		if v, ok := metadata["IssueLogo"]; ok {
			issueLogo = fmt.Sprintf("%s", v)
		}
		if len(issueCompanyName) > 0 && len(issueAddress) > 0 {
			targetMetaData["IssueCompanyName"] = issueCompanyName
			targetMetaData["IssueAddress"] = issueAddress
			targetMetaData["IssueVatNumber"] = issueVatNumber
			targetMetaData["IssueRegNumber"] = issueRegNumber
		}
		if len(issueLogo) > 0 {
			targetMetaData["IssueLogo"] = issueLogo
		}
	}
}
