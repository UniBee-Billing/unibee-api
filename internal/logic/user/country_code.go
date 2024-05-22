package user

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	config2 "unibee/internal/cmd/config"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/vat_gateway"
	"unibee/internal/query"
	"unibee/utility"
)

func UpdateUserCountryCode(ctx context.Context, userId uint64, countryCode string) {
	utility.Assert(userId > 0, "userId is nil")
	user := query.GetUserAccountById(ctx, userId)
	utility.Assert(user != nil, "UpdateUserCountryCode user not found")
	if len(countryCode) > 0 && strings.Compare(user.CountryCode, countryCode) != 0 {
		if vat_gateway.GetDefaultVatGateway(ctx, user.MerchantId) != nil {
			vatCountryRate, err := vat_gateway.QueryVatCountryRateByMerchant(ctx, user.MerchantId, countryCode)
			var taxPercentage int64 = 0
			if err == nil && vatCountryRate != nil {
				if len(user.VATNumber) > 0 && !strings.Contains(config2.GetConfigInstance().VatConfig.NumberUnExemptionCountryCodes, countryCode) {
					taxPercentage = 0
				} else if vatCountryRate.StandardTaxPercentage > 0 {
					taxPercentage = vatCountryRate.StandardTaxPercentage
				}
			}
			if err == nil && vatCountryRate != nil {
				_, err = dao.UserAccount.Ctx(ctx).Data(g.Map{
					dao.UserAccount.Columns().CountryCode:   countryCode,
					dao.UserAccount.Columns().CountryName:   vatCountryRate.CountryName,
					dao.UserAccount.Columns().TaxPercentage: taxPercentage,
					dao.UserAccount.Columns().GmtModify:     gtime.Now(),
				}).Where(dao.UserAccount.Columns().Id, user.Id).OmitNil().Update()
				if err != nil {
					g.Log().Errorf(ctx, "UpdateUserCountryCode userId:%d CountryCode:%d, error:%s", userId, countryCode, err.Error())
				} else {
					g.Log().Errorf(ctx, "UpdateUserCountryCode userId:%d CountryCode:%d, success", userId, countryCode)
				}
			}
		}
	}
}

func GetUserCountryCode(ctx context.Context, userId uint64) (countryCode string, countryName string) {
	utility.Assert(userId > 0, "userId is nil")
	user := query.GetUserAccountById(ctx, userId)
	utility.Assert(user != nil, "GetUserCountryCode user not found")
	return user.CountryCode, user.CountryName
}

func GetUserTaxPercentage(ctx context.Context, userId uint64) (int64, error) {
	utility.Assert(userId > 0, "userId is nil")
	user := query.GetUserAccountById(ctx, userId)
	utility.Assert(user != nil, "GetUserCountryCode user not found")
	if vat_gateway.GetDefaultVatGateway(ctx, user.MerchantId) != nil {
		vatCountryRate, err := vat_gateway.QueryVatCountryRateByMerchant(ctx, user.MerchantId, user.CountryCode)
		var taxPercentage int64 = 0
		if err == nil && vatCountryRate != nil {
			if len(user.VATNumber) > 0 && !strings.Contains(config2.GetConfigInstance().VatConfig.NumberUnExemptionCountryCodes, user.CountryCode) {
				taxPercentage = 0
			} else if vatCountryRate.StandardTaxPercentage > 0 {
				taxPercentage = vatCountryRate.StandardTaxPercentage
			}
		}
		return taxPercentage, nil
	}
	return -1, gerror.New("Default Vat Gateway Need Setup")
}
