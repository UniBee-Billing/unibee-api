package vat_gateway

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	config2 "unibee/internal/cmd/config"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func UpdateUserVatNumber(ctx context.Context, userId uint64, vatNumber string) {
	utility.Assert(userId > 0, "userId is nil")
	user := query.GetUserAccountById(ctx, userId)
	utility.Assert(user != nil, "UpdateUserVatNumber user not found")
	if len(vatNumber) > 0 {
		var taxPercentage = user.TaxPercentage
		if len(user.VATNumber) > 0 && !strings.Contains(config2.GetConfigInstance().VatConfig.NumberUnExemptionCountryCodes, user.CountryCode) {
			taxPercentage = 0
		}
		_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
			dao.UserAccount.Columns().VATNumber:     vatNumber,
			dao.UserAccount.Columns().TaxPercentage: taxPercentage,
			dao.UserAccount.Columns().GmtModify:     gtime.Now(),
		}).Where(dao.UserAccount.Columns().Id, user.Id).OmitNil().Update()
		if err != nil {
			g.Log().Errorf(ctx, "UpdateUserVatNumber userId:%d vatNumber:%d, error:%s", userId, vatNumber, err.Error())
		} else {
			g.Log().Errorf(ctx, "UpdateUserVatNumber userId:%d vatNumber:%d, success", userId, vatNumber)
		}
	}
}
