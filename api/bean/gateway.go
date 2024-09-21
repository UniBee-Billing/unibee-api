package bean

import (
	"unibee/internal/consts"
	gateway2 "unibee/internal/logic/gateway"
	entity "unibee/internal/model/entity/default"
	"unicode"
)

type Gateway struct {
	Id                 uint64          `json:"gatewayId"`
	GatewayName        string          `json:"gatewayName"`
	DisplayName        string          `json:"displayName"`
	GatewayLogo        string          `json:"gatewayLogo"`
	GatewayKey         string          `json:"gatewayKey"            description:""`                                                        //
	GatewayType        int64           `json:"gatewayType"           description:"gateway type，1-Bank Card ｜ 2-Crypto | 3 - Wire Transfer"` // gateway type，1-Default｜ 2-Crypto
	CountryConfig      map[string]bool `json:"countryConfig"`
	CreateTime         int64           `json:"createTime"            description:"create utc time"` // create utc time
	MinimumAmount      int64           `json:"minimumAmount"   description:"The minimum amount of wire transfer" `
	Currency           string          `json:"currency"   description:"The currency of wire transfer " `
	Bank               *GatewayBank    `json:"bank"   dc:"The receiving bank of wire transfer" `
	WebhookEndpointUrl string          `json:"webhookEndpointUrl"   description:"The endpoint url of gateway webhook " `
	WebhookSecret      string          `json:"webhookSecret"  dc:"The secret of gateway webhook"`
}

type GatewayBank struct {
	AccountHolder string `json:"accountHolder"   dc:"The AccountHolder of wire transfer " v:"required" `
	BIC           string `json:"bic"   dc:"The BIC of wire transfer " v:"required" `
	IBAN          string `json:"iban"   dc:"The IBAN of wire transfer " v:"required" `
	Address       string `json:"address"   dc:"The address of wire transfer " v:"required" `
}

func SimplifyGateway(one *entity.MerchantGateway) *Gateway {
	if one == nil {
		return nil
	}
	var countryConfig map[string]bool
	_ = UnmarshalFromJsonString(one.CountryConfig, &countryConfig)
	var bank *GatewayBank
	_ = UnmarshalFromJsonString(one.BankData, &bank)
	var webhookEndpointUrl = ""
	if one.GatewayType != consts.GatewayTypeWireTransfer {
		webhookEndpointUrl = gateway2.GetPaymentWebhookEntranceUrl(one.Id)
	}

	return &Gateway{
		Id:                 one.Id,
		GatewayLogo:        one.Logo,
		GatewayName:        one.GatewayName,
		DisplayName:        toUpperFirst(one.GatewayName, one.Name),
		GatewayType:        one.GatewayType,
		CountryConfig:      countryConfig,
		CreateTime:         one.CreateTime,
		Currency:           one.Currency,
		MinimumAmount:      one.MinimumAmount,
		Bank:               bank,
		WebhookEndpointUrl: webhookEndpointUrl,
		GatewayKey:         one.GatewayKey,
		WebhookSecret:      one.WebhookSecret,
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

func SimplifyGatewayList(ones []*entity.MerchantGateway) (list []*Gateway) {
	if len(ones) == 0 {
		return make([]*Gateway, 0)
	}
	for _, one := range ones {
		list = append(list, SimplifyGateway(one))
	}
	return list
}
