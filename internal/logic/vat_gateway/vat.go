package vat_gateway

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	vat2 "unibee/internal/interface"
	"unibee/internal/logic/merchant_config"
	"unibee/internal/logic/merchant_config/update"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/subscription/config"
	vat "unibee/internal/logic/vat_gateway/github"
	"unibee/internal/logic/vat_gateway/vatsense"
	"unibee/internal/logic/vat_gateway/vatstack"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

const (
	KeyMerchantVatName = "KEY_MERCHANT_DEFAULT_VAT_NAME"
)

const (
	VAT_IMPLEMENT_NAMES = "vatsense|github|vatstack"
)

func GetDefaultVatGateway(ctx context.Context, merchantId uint64) vat2.VATGateway {
	vatName, vatData := GetDefaultMerchantVatConfig(ctx, merchantId)
	if len(vatName) == 0 {
		return nil
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
	}
	return nil
}

func GetDefaultMerchantVatConfig(ctx context.Context, merchantId uint64) (vatName string, data string) {
	nameConfig := merchant_config.GetMerchantConfig(ctx, merchantId, KeyMerchantVatName)
	if nameConfig != nil {
		vatName = nameConfig.ConfigValue
	}
	if len(vatName) == 0 {
		return
	}
	valueConfig := merchant_config.GetMerchantConfig(ctx, merchantId, vatName)
	if valueConfig != nil {
		data = valueConfig.ConfigValue
	}
	return
}

func SetupMerchantVatConfig(ctx context.Context, merchantId uint64, vatName string, data string, isDefault bool) error {
	utility.Assert(strings.Contains(VAT_IMPLEMENT_NAMES, vatName), "gateway not support, should be "+VAT_IMPLEMENT_NAMES)
	err := update.SetMerchantConfig(ctx, merchantId, vatName, data)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     merchantId,
		Target:         fmt.Sprintf("Vat(%s)", vatName),
		Content:        "SetupVatGateway",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	if err != nil {
		return err
	}
	if isDefault {
		err = update.SetMerchantConfig(ctx, merchantId, KeyMerchantVatName, vatName)
	}
	return err
}

func CleanMerchantDefaultVatConfig(ctx context.Context, merchantId uint64) error {
	return update.SetMerchantConfig(ctx, merchantId, KeyMerchantVatName, "")
}

func InitMerchantDefaultVatGateway(ctx context.Context, merchantId uint64) error {
	gateway := GetDefaultVatGateway(ctx, merchantId)
	if gateway == nil {
		return gerror.New("Default Vat Gateway Need Setup")
	}
	countries, err := gateway.ListAllCountries()
	if err != nil {
		g.Log().Infof(ctx, "InitMerchantDefaultVatGateway ListAllCountries err merchantId:%d gatewayName:%s err:%v", merchantId, gateway.GetGatewayName(), err)
		return err
	}
	for _, country := range countries {
		country.MerchantId = merchantId
	}
	if countries != nil && len(countries) > 0 {
		for _, newOne := range countries {
			var one *entity.CountryRate
			err = dao.CountryRate.Ctx(ctx).
				Where(dao.CountryRate.Columns().MerchantId, newOne.MerchantId).
				Where(dao.CountryRate.Columns().Gateway, newOne.Gateway).
				Where(dao.CountryRate.Columns().CountryCode, newOne.CountryCode).
				Scan(&one)
			if err != nil {
				return err
			}
			if one != nil {
				_, err = dao.CountryRate.Ctx(ctx).Data(g.Map{
					dao.CountryRate.Columns().CountryName: newOne.CountryName,
					dao.CountryRate.Columns().Latitude:    newOne.Latitude,
					dao.CountryRate.Columns().Longitude:   newOne.Longitude,
					dao.CountryRate.Columns().Vat:         newOne.Vat,
					dao.CountryRate.Columns().GmtModify:   gtime.Now(),
				}).Where(dao.CountryRate.Columns().Id, one.Id).OmitNil().Update()
				if err != nil {
					g.Log().Errorf(ctx, "InitMerchantDefaultVatGateway Save Countries error:%s", err.Error())
					return err
				}
			} else {
				_, err = dao.CountryRate.Ctx(ctx).Data(newOne).OmitEmpty().Insert()
				if err != nil {
					g.Log().Infof(ctx, "InitMerchantDefaultVatGateway Save Countries err merchantId:%d gatewayName:%s err:%v", merchantId, gateway.GetGatewayName(), err)
					return err
				}
			}
		}
	}
	countryRates, err := gateway.ListAllRates()
	if err != nil {
		g.Log().Infof(ctx, "InitMerchantDefaultVatGateway ListAllRates err merchantId:%d gatewayName:%s err:%v", merchantId, gateway.GetGatewayName(), err)
		return err
	}
	for _, country := range countryRates {
		country.MerchantId = merchantId
	}
	if countryRates != nil && len(countryRates) > 0 {
		for _, newOne := range countryRates {
			var one *entity.CountryRate
			err = dao.CountryRate.Ctx(ctx).
				Where(dao.CountryRate.Columns().MerchantId, newOne.MerchantId).
				Where(dao.CountryRate.Columns().Gateway, newOne.Gateway).
				Where(dao.CountryRate.Columns().CountryCode, newOne.CountryCode).
				Scan(&one)
			if err != nil {
				g.Log().Infof(ctx, "InitMerchantDefaultVatGateway Save All Rates err merchantId:%d gatewayName:%s err:%v", merchantId, gateway.GetGatewayName(), err)
				return err
			}
			if one != nil {
				_, err = dao.CountryRate.Ctx(ctx).Data(g.Map{
					dao.CountryRate.Columns().CountryName:           newOne.CountryName,
					dao.CountryRate.Columns().Latitude:              newOne.Latitude,
					dao.CountryRate.Columns().Longitude:             newOne.Longitude,
					dao.CountryRate.Columns().Vat:                   newOne.Vat,
					dao.CountryRate.Columns().Other:                 newOne.Other,
					dao.CountryRate.Columns().Provinces:             newOne.Provinces,
					dao.CountryRate.Columns().Mamo:                  newOne.Mamo,
					dao.CountryRate.Columns().Eu:                    newOne.Eu,
					dao.CountryRate.Columns().StandardTaxPercentage: newOne.StandardTaxPercentage,
					dao.CountryRate.Columns().StandardTypes:         newOne.StandardTypes,
					dao.CountryRate.Columns().StandardDescription:   newOne.StandardDescription,
					dao.CountryRate.Columns().GmtModify:             gtime.Now(),
				}).Where(dao.CountryRate.Columns().Id, one.Id).OmitNil().Update()
				if err != nil {
					g.Log().Errorf(ctx, "InitMerchantDefaultVatGateway Save Countries error:%s", err.Error())
					return err
				}
			} else {
				_, err = dao.CountryRate.Ctx(ctx).Data(newOne).OmitNil().Insert()
				if err != nil {
					g.Log().Infof(ctx, "InitMerchantDefaultVatGateway Save All Rates err merchantId:%d gatewayName:%s err:%v", merchantId, gateway.GetGatewayName(), err)
					return err
				}
			}
		}
	}

	return nil
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
	var list []*bean.VatCountryRate
	for _, countryRate := range countryRateList {
		var vatSupport = false
		if countryRate.Vat == 1 {
			vatSupport = true
		} else {
			vatSupport = false
		}
		// disable tax for non-eu country
		var standardTaxPercentage = countryRate.StandardTaxPercentage
		if countryRate.Eu != 1 {
			standardTaxPercentage = 0
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
		return nil, gerror.New("Default Vat Gateway Need Setup")
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
	// disable tax for non-eu country
	var standardTaxPercentage = one.StandardTaxPercentage

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
	if GetDefaultVatGateway(ctx, merchantId) != nil {
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
		g.Log().Infof(ctx, "Default Vat Gateway Need Setup:"+strconv.FormatUint(merchantId, 10))
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
