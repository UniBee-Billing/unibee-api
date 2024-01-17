package vat_gateway

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/merchant_config"
	"go-oversea-pay/internal/logic/vat_gateway/base"
	vat "go-oversea-pay/internal/logic/vat_gateway/github"
	"go-oversea-pay/internal/logic/vat_gateway/vatsense"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"strings"
)

const (
	KeyMerchantVatName = "KEY_MERCHANT_DEFAULT_VAT_NAME"
)

const (
	VAT_IMPLEMENT_NAMES = "vatsense|github"
)

type VatCountryRate struct {
	CountryCode           string `json:"countryCode"           `                                      // country_code
	CountryName           string `json:"countryName"           `                                      // country_name
	VatSupport            bool   `json:"vatSupport"          dc:"vat support,true or false"         ` // vat support true or false
	StandardTaxPercentage int64  `json:"standardTaxPercentage"  dc:"Tax税率，万分位，1000 表示 10%"`
}

func GetDefaultVatGateway(ctx context.Context, merchantId int64) base.Gateway {
	vatName, vatData := getDefaultMerchantVatConfig(ctx, merchantId)
	if len(vatName) == 0 {
		return nil
	}
	if strings.Compare(vatName, "vatsense") == 0 {
		one := &vatsense.VatSense{Password: vatData, Name: vatName}
		return one
	} else if strings.Compare(vatName, "github") == 0 {
		one := &vat.Github{Password: vatData, Name: vatName}
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
	if countries != nil && len(countries) > 0 {
		_, err = dao.CountryRate.Ctx(ctx).Data(countries).OmitEmpty().Save(countries)
		if err != nil {
			g.Log().Infof(ctx, "InitMerchantDefaultVatGateway Save Countries err merchantId:%d gateway:%s err:%v", merchantId, gateway.GetVatName(), err)
			return err
		}
	}
	countryRates, err := gateway.ListAllRates()
	if err != nil {
		g.Log().Infof(ctx, "InitMerchantDefaultVatGateway ListAllRates err merchantId:%d gateway:%s err:%v", merchantId, gateway.GetVatName(), err)
		return err
	}
	if countryRates != nil && len(countryRates) > 0 {
		if countries == nil || len(countries) == 0 {
			//Country 没数据，全覆盖
			_, err = dao.CountryRate.Ctx(ctx).Data(countryRates).OmitEmpty().Replace()
		} else {
			_, err = dao.CountryRate.Ctx(ctx).Data(countryRates).OnDuplicate(
				dao.CountryRate.Columns().StandardTypes,
				dao.CountryRate.Columns().StandardDescription,
				dao.CountryRate.Columns().StandardTaxPercentage,
				dao.CountryRate.Columns().Other).
				OmitEmpty().Save()
		}
		if err != nil {
			g.Log().Infof(ctx, "InitMerchantDefaultVatGateway Save All Rates err merchantId:%d gateway:%s err:%v", merchantId, gateway.GetVatName(), err)
			return err
		}
	}

	return nil
}

func GetVatNumberValidateHistory(ctx context.Context, merchantId int64, vatNumber string) (res *entity.MerchantVatNumberValicationHistory) {
	err := dao.MerchantVatNumberValicationHistory.Ctx(ctx).
		Where(entity.MerchantVatNumberValicationHistory{MerchantId: merchantId}).
		Where(entity.MerchantVatNumberValicationHistory{VatNumber: vatNumber}).OmitEmpty().Scan(&res)
	if err != nil {
		return nil
	}
	return res
}

func ValidateVatNumberByDefaultGateway(ctx context.Context, merchantId int64, vatNumber string, requestVatNumber string) (*base.ValidResult, error) {
	one := GetVatNumberValidateHistory(ctx, merchantId, vatNumber)
	if one != nil {
		var valid = false
		if one.Valid == 1 {
			valid = true
		}
		return &base.ValidResult{
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
		g.Log().Infof(ctx, "InitMerchantDefaultVatGateway merchant gateway data not setup merchantId:%d gateway:%s", merchantId, gateway.GetVatName())
		return nil, gerror.New("default vat gateway not setup")
	}
	result, validateError := gateway.ValidateVatNumber(vatNumber, requestVatNumber)
	if validateError != nil {
		return nil, validateError
	}
	var valid = 0
	if result.Valid {
		valid = 1
	}
	one = &entity.MerchantVatNumberValicationHistory{
		MerchantId:      merchantId,
		VatNumber:       vatNumber,
		Valid:           int64(valid),
		ValidateChannel: gateway.GetVatName(),
		CountryCode:     result.CountryCode,
		CompanyName:     result.CompanyName,
		CompanyAddress:  result.CompanyAddress,
		ValidateMessage: result.ValidateMessage,
	}
	_, err := dao.MerchantVatNumberValicationHistory.Ctx(ctx).Data(one).OmitEmpty().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`ValidateVatNumberByDefaultGateway record insert failure %s`, err)
	}
	return result, nil
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

func QueryVatCountryRateByMerchant(ctx context.Context, merchantId int64, countryCode string) (*VatCountryRate, error) {
	gateway := GetDefaultVatGateway(ctx, merchantId)
	if gateway == nil {
		g.Log().Infof(ctx, "MerchantCountryRateList merchant gateway data not setup merchantId:%d gateway:%s", merchantId, gateway.GetVatName())
		return nil, gerror.New("default vat gateway not setup")
	}
	var one *entity.CountryRate
	err := dao.CountryRate.Ctx(ctx).
		Where(dao.CountryRate.Columns().IsDeleted, 0).
		Where(dao.CountryRate.Columns().VatName, gateway.GetVatName()).
		Where(dao.CountryRate.Columns().CountryCode, countryCode).
		OmitEmpty().Scan(&one)
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
	return &VatCountryRate{
		CountryCode:           one.CountryCode,
		CountryName:           one.CountryName,
		VatSupport:            vatSupport,
		StandardTaxPercentage: one.StandardTaxPercentage,
	}, nil
}
