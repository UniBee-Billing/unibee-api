package vat_gateway

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/merchant_config"
	"go-oversea-pay/internal/logic/vat_gateway/base"
	"go-oversea-pay/internal/logic/vat_gateway/vatsense"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"strings"
)

const (
	KeyMerchantVatName = "KEY_MERCHANT_DEFAULT_VAT_NAME"
)

const (
	VAT_IMPLEMENT_NAMES = "vatsense"
)

type VatCountryRate struct {
	CountryCode           string `json:"countryCode"           `                                      // country_code
	CountryName           string `json:"countryName"           `                                      // country_name
	VatSupport            bool   `json:"vatSupport"          dc:"vat support,true or false"         ` // vat support true or false
	StandardTaxPercentage int64  `json:"standardTaxPercentage"  dc:"Standard Tax百分比，10 表示 10%"`       // Standard Tax百分比，10 表示 10%
}

func GetDefaultVatGateway(ctx context.Context, merchantId int64) base.Gateway {
	vatName, vatData := getDefaultMerchantVatConfig(ctx, merchantId)
	if len(vatName) == 0 {
		return nil
	}
	if strings.Compare(vatName, "vatsense") == 0 {
		one := &vatsense.VatSense{Password: vatData, Name: vatName}
		return one
	}
	return nil
}

func getDefaultMerchantVatConfig(ctx context.Context, merchantId int64) (vatName string, data string) {
	nameConfig := merchant_config.GetMerchantConfig(ctx, merchantId, KeyMerchantVatName)
	if nameConfig != nil {
		vatName = nameConfig.ConfigValue
	}
	valueConfig := merchant_config.GetMerchantConfig(ctx, merchantId, vatName)
	if valueConfig != nil {
		data = valueConfig.ConfigValue
	}
	return
}

func SetupMerchantVatConfig(ctx context.Context, merchantId int64, vatName string, data string, isDefault bool) error {
	if !strings.Contains(VAT_IMPLEMENT_NAMES, vatName) {
		return gerror.New("Vat gateway not support")
	}
	err := merchant_config.SetMerchantConfig(ctx, merchantId, vatName, data)
	if err != nil {
		return err
	}
	if isDefault {
		err = merchant_config.SetMerchantConfig(ctx, merchantId, KeyMerchantVatName, vatName)
	}
	return err
}

func InitMerchantDefaultVatGateway(ctx context.Context, merchantId int64) error {
	gateway := GetDefaultVatGateway(ctx, merchantId)
	if gateway == nil {
		g.Log().Infof(ctx, "InitMerchantDefaultVatGateway merchant gateway data not setup merchantId:%d gateway:%s", merchantId, gateway.GetVatName())
	}
	countries, err := gateway.ListAllCountries()
	if err != nil {
		g.Log().Infof(ctx, "InitMerchantDefaultVatGateway ListAllCountries err merchantId:%d gateway:%s err:%v", merchantId, gateway.GetVatName(), err)
		return err
	}
	_, err = dao.CountryRate.Ctx(ctx).Data(countries).OmitEmpty().Save(countries)
	if err != nil {
		g.Log().Infof(ctx, "InitMerchantDefaultVatGateway Save Countries err merchantId:%d gateway:%s err:%v", merchantId, gateway.GetVatName(), err)
		return err
	}
	countries, err = gateway.ListAllRates()
	if err != nil {
		g.Log().Infof(ctx, "InitMerchantDefaultVatGateway ListAllRates err merchantId:%d gateway:%s err:%v", merchantId, gateway.GetVatName(), err)
		return err
	}
	_, err = dao.CountryRate.Ctx(ctx).Data(countries).OnDuplicate(
		dao.CountryRate.Columns().StandardTypes,
		dao.CountryRate.Columns().StandardDescription,
		dao.CountryRate.Columns().StandardTaxPercentage,
		dao.CountryRate.Columns().Other).
		OmitEmpty().Save()
	if err != nil {
		g.Log().Infof(ctx, "InitMerchantDefaultVatGateway Save All Rates err merchantId:%d gateway:%s err:%v", merchantId, gateway.GetVatName(), err)
		return err
	}
	return nil
}

func ValidateVatNumberByDefaultGateway(ctx context.Context, merchantId int64, vatNumber string, requestVatNumber string) (*base.ValidResult, error) {
	gateway := GetDefaultVatGateway(ctx, merchantId)
	if gateway == nil {
		g.Log().Infof(ctx, "InitMerchantDefaultVatGateway merchant gateway data not setup merchantId:%d gateway:%s", merchantId, gateway.GetVatName())
		return nil, gerror.New("default vat gateway not setup")
	}
	return gateway.ValidateVatNumber(vatNumber, requestVatNumber)
}

func MerchantCountryRateList(ctx context.Context, merchantId int64) ([]*VatCountryRate, error) {
	gateway := GetDefaultVatGateway(ctx, merchantId)
	if gateway == nil {
		g.Log().Infof(ctx, "MerchantCountryRateList merchant gateway data not setup merchantId:%d gateway:%s", merchantId, gateway.GetVatName())
		return nil, gerror.New("default vat gateway not setup")
	}
	var countryRateList []*entity.CountryRate
	err := dao.CountryRate.Ctx(ctx).
		Where(dao.CountryRate.Columns().IsDeleted, 0).
		Where(dao.CountryRate.Columns().VatName, gateway.GetVatName()).
		Order("country_name").
		OmitEmpty().Scan(&countryRateList)
	if err != nil {
		return nil, err
	}
	var list []*VatCountryRate
	for _, countryRate := range countryRateList {
		var vatSupport = false
		if countryRate.Vat == 1 {
			vatSupport = true
		} else {
			vatSupport = false
		}
		list = append(list, &VatCountryRate{
			CountryCode:           countryRate.CountryCode,
			CountryName:           countryRate.CountryName,
			VatSupport:            vatSupport,
			StandardTaxPercentage: countryRate.StandardTaxPercentage,
		})
	}
	return list, nil
}
