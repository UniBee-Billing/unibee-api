package bean

import entity "unibee/internal/model/entity/default"

type MerchantCountryConfig struct {
	MerchantId  uint64 `json:"merchantId" description:"merchant id"` // merchant id
	CountryCode string `json:"countryCode" description:""`           //
	Name        string `json:"name"        description:""`           //
}

func SimplifyMerchantCountryConfig(one *entity.MerchantCountryConfig) *MerchantCountryConfig {
	if one == nil {
		return nil
	}
	return &MerchantCountryConfig{
		MerchantId:  one.MerchantId,
		CountryCode: one.CountryCode,
		Name:        one.Name,
	}
}

type MerchantVatRule struct {
	GatewayNames      string `json:"gatewayNames" dc:""`
	ValidCountryCodes string `json:"validCountryCodes" dc:""`
	TaxPercentage     *int64 `json:"taxPercentage" dc:""`
	IgnoreVatNumber   bool   `json:"ignoreVatNumber" dc:""`
}
