package vat_gateway

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/merchant_config"
	"unibee/internal/logic/subscription/config"
	"unibee/internal/logic/vat_gateway/default_vat_gateway"
	vat "unibee/internal/logic/vat_gateway/github"
	"unibee/internal/logic/vat_gateway/vatsense"
	"unibee/internal/logic/vat_gateway/vatstack"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

const (
	KeyMerchantVatName          = "KEY_MERCHANT_DEFAULT_VAT_NAME"
	KeyMerchantVatGeneralConfig = "KEY_MERCHANT_DEFAULT_GENERAL_CONFIG"
)

const (
	VAT_IMPLEMENT_NAMES  = "default|vatsense|github|vatstack"
	DEFAULT_GATEWAY_NAME = "default"
)

type VATGeneralConfig struct {
	ValidForNonEU     bool     `json:"valid_non_eu" dc:" valid for non eu"`
	ValidCountryCodes []string `json:"valid_country_codes" dc:"valid for all if empty list and blank"`
}

func GetDefaultVatGateway(ctx context.Context, merchantId uint64) VATGateway {
	vatName, vatData := GetDefaultMerchantVatConfig(ctx, merchantId)
	if len(vatName) == 0 {
		one := &default_vat_gateway.DefaultVatGateway{Name: DEFAULT_GATEWAY_NAME}
		return one
	}
	if strings.Compare(vatName, "vatsense") == 0 {
		one := &vatsense.VatSense{Password: vatData, Name: vatName}
		return one
	} else if strings.Compare(vatName, "github") == 0 {
		one := &vat.Github{Password: vatData, Name: vatName}
		return one
	} else if strings.Compare(vatName, "vatstack") == 0 {
		one := &vatstack.VatStack{
			ApiData: vatData,
			Name:    vatName,
		}
		return one
	} else if strings.Compare(vatName, DEFAULT_GATEWAY_NAME) == 0 {
		one := &default_vat_gateway.DefaultVatGateway{Name: DEFAULT_GATEWAY_NAME}
		return one
	}
	return nil
}

func GetMerchantVATGeneralConfig(ctx context.Context, merchantId uint64) *VATGeneralConfig {
	generalConfig := merchant_config.GetMerchantConfig(ctx, merchantId, KeyMerchantVatGeneralConfig)
	if generalConfig != nil && len(generalConfig.ConfigValue) > 0 {
		var one *VATGeneralConfig
		_ = utility.UnmarshalFromJsonString(generalConfig.ConfigValue, &one)
		return one
	}
	return &VATGeneralConfig{
		ValidForNonEU:     false,
		ValidCountryCodes: make([]string, 0),
	}
}

func GetDefaultMerchantVatConfig(ctx context.Context, merchantId uint64) (vatName string, data string) {
	nameConfig := merchant_config.GetMerchantConfig(ctx, merchantId, KeyMerchantVatName)
	if nameConfig != nil {
		vatName = nameConfig.ConfigValue
	}
	if len(vatName) == 0 {
		// default vat build-in gateway
		return DEFAULT_GATEWAY_NAME, ""
	}
	valueConfig := merchant_config.GetMerchantConfig(ctx, merchantId, vatName)
	if valueConfig != nil {
		data = valueConfig.ConfigValue
	}
	return
}

func ValidateVatNumberByDefaultGateway(ctx context.Context, merchantId uint64, userId uint64, vatNumber string, requestVatNumber string) (*bean.ValidResult, error) {
	if len(vatNumber) == 0 {
		return &bean.ValidResult{
			Valid:           false,
			VatNumber:       "",
			CountryCode:     "",
			CompanyName:     "",
			CompanyAddress:  "",
			ValidateMessage: "",
		}, nil
	}
	one := query.GetVatNumberValidateHistory(ctx, merchantId, vatNumber)
	if one != nil {
		var valid = false
		if one.Valid == 1 {
			valid = true
		}
		return &bean.ValidResult{
			Valid:           valid,
			VatNumber:       one.VatNumber,
			CountryCode:     one.CountryCode,
			CompanyName:     one.CompanyName,
			CompanyAddress:  one.CompanyAddress,
			ValidateMessage: one.ValidateMessage,
		}, nil
	}
	gateway := GetDefaultVatGateway(ctx, merchantId)
	if gateway == nil {
		return nil, gerror.New("Default Vat Gateway Need Setup")
	}
	result, validateError := gateway.ValidateVatNumber(vatNumber, requestVatNumber)
	if validateError != nil {
		return nil, validateError
	}
	var valid = 0
	if result.Valid {
		valid = 1
	}
	one = &entity.MerchantVatNumberVerifyHistory{
		MerchantId:      merchantId,
		VatNumber:       vatNumber,
		Valid:           int64(valid),
		ValidateGateway: gateway.GetGatewayName(),
		CountryCode:     result.CountryCode,
		CompanyName:     result.CompanyName,
		CompanyAddress:  result.CompanyAddress,
		ValidateMessage: result.ValidateMessage,
		CreateTime:      gtime.Now().Timestamp(),
	}
	_, err := dao.MerchantVatNumberVerifyHistory.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`ValidateVatNumberByDefaultGateway record insert failure %s`, err)
	}
	return result, nil
}

func MerchantCountryRateList(ctx context.Context, merchantId uint64) ([]*bean.VatCountryRate, error) {
	gateway := GetDefaultVatGateway(ctx, merchantId)
	if gateway == nil {
		return make([]*bean.VatCountryRate, 0), gerror.New("Default Vat Gateway Need Setup")
	}
	var countryRateList []*entity.CountryRate
	err := dao.CountryRate.Ctx(ctx).
		Where(dao.CountryRate.Columns().MerchantId, merchantId).
		Where(dao.CountryRate.Columns().IsDeleted, 0).
		Where(dao.CountryRate.Columns().Gateway, gateway.GetGatewayName()).
		Order("country_name").
		Scan(&countryRateList)
	if err != nil {
		return nil, err
	}
	generalConfig := GetMerchantVATGeneralConfig(ctx, merchantId)
	var list []*bean.VatCountryRate
	for _, countryRate := range countryRateList {
		var vatSupport = false
		if countryRate.Vat == 1 {
			vatSupport = true
		} else {
			vatSupport = false
		}
		var standardTaxPercentage = countryRate.StandardTaxPercentage
		if len(generalConfig.ValidCountryCodes) > 0 {
			if !utility.IsStringInArray(generalConfig.ValidCountryCodes, countryRate.CountryCode) {
				standardTaxPercentage = 0
			}
		}
		list = append(list, &bean.VatCountryRate{
			CountryCode:           countryRate.CountryCode,
			CountryName:           countryRate.CountryName,
			VatSupport:            vatSupport,
			IsEU:                  countryRate.Eu == 1,
			StandardTaxPercentage: standardTaxPercentage,
		})
	}
	return list, nil
}

func QueryVatCountryRateByMerchant(ctx context.Context, merchantId uint64, countryCode string) (*bean.VatCountryRate, error) {
	gateway := GetDefaultVatGateway(ctx, merchantId)
	if gateway == nil {
		return nil, gerror.New("Vat Gateway Need Setup")
	}
	var one *entity.CountryRate
	err := dao.CountryRate.Ctx(ctx).
		Where(dao.CountryRate.Columns().MerchantId, merchantId).
		Where(dao.CountryRate.Columns().IsDeleted, 0).
		Where(dao.CountryRate.Columns().Gateway, gateway.GetGatewayName()).
		Where(dao.CountryRate.Columns().CountryCode, countryCode).
		Scan(&one)
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, gerror.New("vat data not found")
	}
	var vatSupport = false
	if one.Vat == 1 {
		vatSupport = true
	} else {
		vatSupport = false
	}
	var standardTaxPercentage = one.StandardTaxPercentage
	generalConfig := GetMerchantVATGeneralConfig(ctx, merchantId)
	if len(generalConfig.ValidCountryCodes) > 0 {
		if !utility.IsStringInArray(generalConfig.ValidCountryCodes, one.CountryCode) {
			standardTaxPercentage = 0
		}
	}
	return &bean.VatCountryRate{
		Id:                    one.Id,
		Gateway:               one.Gateway,
		CountryCode:           one.CountryCode,
		CountryName:           one.CountryName,
		VatSupport:            vatSupport,
		IsEU:                  one.Eu == 1,
		StandardTaxPercentage: standardTaxPercentage,
	}, nil
}

func ComputeMerchantVatPercentage(ctx context.Context, merchantId uint64, countryCode string, gatewayId uint64, validVatNumber string) (taxPercentage int64, countryName string) {
	if GetDefaultVatGateway(ctx, merchantId).VatRatesEnabled() {
		vatCountryRate, err := QueryVatCountryRateByMerchant(ctx, merchantId, countryCode)
		if err == nil && vatCountryRate != nil {
			countryName = vatCountryRate.CountryName
			var ignoreVatNumber = false
			if len(config.GetMerchantSubscriptionConfig(ctx, merchantId).GatewayVATRule) > 0 {
				var gatewayName string
				gateway := query.GetGatewayById(ctx, gatewayId)
				if gateway != nil {
					gatewayName = gateway.GatewayName
				}
				var gatewayVATRules = make([]*bean.MerchantVatRule, 0)
				_ = utility.UnmarshalFromJsonString(config.GetMerchantSubscriptionConfig(ctx, merchantId).GatewayVATRule, &gatewayVATRules)
				if len(gatewayVATRules) > 0 {
					for _, gatewayVatRule := range gatewayVATRules {
						if ruleContain(gatewayVatRule.GatewayNames, gatewayName) && ruleContain(gatewayVatRule.ValidCountryCodes, countryCode) {
							if gatewayVatRule.TaxPercentage != nil && *gatewayVatRule.TaxPercentage > 0 {
								taxPercentage = *gatewayVatRule.TaxPercentage
							} else {
								taxPercentage = vatCountryRate.StandardTaxPercentage
							}
							ignoreVatNumber = gatewayVatRule.IgnoreVatNumber
							break
						}
					}
				} else {
					taxPercentage = vatCountryRate.StandardTaxPercentage
				}
			} else {
				taxPercentage = vatCountryRate.StandardTaxPercentage
			}
			if len(validVatNumber) > 0 && !ignoreVatNumber {
				taxPercentage = 0
			}
		}
	} else {
		g.Log().Infof(ctx, "Vat Gateway Need Setup:"+strconv.FormatUint(merchantId, 10))
	}
	return taxPercentage, countryName
}

func ruleContain(rules string, target string) bool {
	if rules == "*" || (len(target) > 0 && strings.Contains(rules, target)) {
		return true
	} else {
		return false
	}
}
