package setup

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/merchant_config/update"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func SetupMerchantVatConfig(ctx context.Context, merchantId uint64, vatName string, data string, isDefault bool) error {
	utility.Assert(strings.Contains(vat_gateway.VAT_IMPLEMENT_NAMES, vatName), "gateway not support, should be "+vat_gateway.VAT_IMPLEMENT_NAMES)
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
		err = update.SetMerchantConfig(ctx, merchantId, vat_gateway.KeyMerchantVatName, vatName)
	}
	return err
}

func CleanMerchantDefaultVatConfig(ctx context.Context, merchantId uint64) error {
	return update.SetMerchantConfig(ctx, merchantId, vat_gateway.KeyMerchantVatName, "")
}

func InitMerchantDefaultVatGateway(ctx context.Context, merchantId uint64) error {
	gateway := vat_gateway.GetDefaultVatGateway(ctx, merchantId)
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
