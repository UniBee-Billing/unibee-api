package vat

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/vat_gateway"
	"unibee/internal/query"
	"unibee/utility"
)

func GetUserCountryCode(ctx context.Context, userId uint64) (countryCode string, countryName string) {
	utility.Assert(userId > 0, "userId is nil")
	user := query.GetUserAccountById(ctx, userId)
	utility.Assert(user != nil, "GetUserCountryCode user not found")
	return user.CountryCode, user.CountryName
}

func GetUserTaxPercentage(ctx context.Context, userId uint64) (taxPercentage int64, countryCode string, vatNumber string, err error) {
	utility.Assert(userId > 0, "userId is nil")
	user := query.GetUserAccountById(ctx, userId)
	utility.Assert(user != nil, fmt.Sprintf("GetUserCountryCode user not found:%v", userId))
	gatewayId, _ := strconv.ParseUint(user.GatewayId, 10, 64)
	if vat_gateway.GetDefaultVatGateway(ctx, user.MerchantId) != nil {
		taxPercentage, _ = vat_gateway.ComputeMerchantVatPercentage(ctx, user.MerchantId, user.CountryCode, gatewayId, user.VATNumber)
		if taxPercentage != user.TaxPercentage && taxPercentage > 0 {
			_, _ = dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().TaxPercentage: taxPercentage,
				dao.UserAccount.Columns().GmtModify:     gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, user.Id).OmitNil().Update()
		}
		return taxPercentage, user.CountryCode, user.VATNumber, nil
	} else {
		return user.TaxPercentage, user.CountryCode, user.VATNumber, nil
	}
}
