package user

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/vat_gateway"
	"unibee/internal/query"
	"unibee/utility"
)

func UpdateUserVatNumber(ctx context.Context, userId uint64, vatNumber string) {
	utility.Assert(userId > 0, "userId is nil")
	user := query.GetUserAccountById(ctx, userId)
	utility.Assert(user != nil, "UpdateUserCountryCode user not found")
	if len(vatNumber) > 0 {
		if vat_gateway.GetDefaultVatGateway(ctx, user.MerchantId) != nil {
			gateway := vat_gateway.GetDefaultVatGateway(ctx, user.MerchantId)
			utility.Assert(gateway != nil, "Default Vat Gateway Need Setup")
			vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, user.MerchantId, user.Id, vatNumber, "")
			if err == nil && vatNumberValidate.Valid {
				_, err = dao.UserAccount.Ctx(ctx).Data(g.Map{
					dao.UserAccount.Columns().VATNumber: vatNumber,
					dao.UserAccount.Columns().GmtModify: gtime.Now(),
				}).Where(dao.UserAccount.Columns().Id, user.Id).OmitNil().Update()
				operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
					MerchantId:     user.MerchantId,
					Target:         fmt.Sprintf("User(%v)", user.Id),
					Content:        "Update(VATNumber)",
					UserId:         user.Id,
					SubscriptionId: "",
					InvoiceId:      "",
					PlanId:         0,
					DiscountCode:   "",
				}, nil)
				if err != nil {
					g.Log().Errorf(ctx, "UpdateUserVatNumber userId:%d vatNumber:%s, error:%s", userId, vatNumber, err.Error())
				} else {
					g.Log().Errorf(ctx, "UpdateUserVatNumber userId:%d vatNumber:%s, success", userId, vatNumber)
					UpdateUserCountryCode(ctx, userId, vatNumberValidate.CountryCode)
				}
			}
		}
	} else {
		_, _ = dao.UserAccount.Ctx(ctx).Data(g.Map{
			dao.UserAccount.Columns().VATNumber: vatNumber,
			dao.UserAccount.Columns().GmtModify: gtime.Now(),
		}).Where(dao.UserAccount.Columns().Id, user.Id).Update()
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     user.MerchantId,
			Target:         fmt.Sprintf("User(%v)", user.Id),
			Content:        "Clear(VATNumber)",
			UserId:         user.Id,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, nil)
	}
}

func UpdateUserCountryCode(ctx context.Context, userId uint64, countryCode string) {
	utility.Assert(userId > 0, "userId is nil")
	user := query.GetUserAccountById(ctx, userId)
	utility.Assert(user != nil, "UpdateUserCountryCode user not found")
	if len(countryCode) > 0 && strings.Compare(user.CountryCode, countryCode) != 0 {
		if vat_gateway.GetDefaultVatGateway(ctx, user.MerchantId) != nil {
			gatewayId, _ := strconv.ParseUint(user.GatewayId, 10, 64)
			taxPercentage, countryName := vat_gateway.ComputeMerchantVatPercentage(ctx, user.MerchantId, countryCode, gatewayId, user.VATNumber)
			_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().CountryCode:   countryCode,
				dao.UserAccount.Columns().CountryName:   countryName,
				dao.UserAccount.Columns().TaxPercentage: taxPercentage,
				dao.UserAccount.Columns().GmtModify:     gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, user.Id).OmitNil().Update()
			operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
				MerchantId:     user.MerchantId,
				Target:         fmt.Sprintf("User(%v)", user.Id),
				Content:        "Update(CountryCode&TaxPercentage)",
				UserId:         user.Id,
				SubscriptionId: "",
				InvoiceId:      "",
				PlanId:         0,
				DiscountCode:   "",
			}, nil)
			if err != nil {
				g.Log().Errorf(ctx, "UpdateUserCountryCode userId:%d CountryCode:%s, error:%s", userId, countryCode, err.Error())
			} else {
				g.Log().Infof(ctx, "UpdateUserCountryCode userId:%d CountryCode:%s, success", userId, countryCode)
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

func GetUserTaxPercentage(ctx context.Context, userId uint64) (taxPercentage int64, vatNumber string, err error) {
	utility.Assert(userId > 0, "userId is nil")
	user := query.GetUserAccountById(ctx, userId)
	utility.Assert(user != nil, fmt.Sprintf("GetUserCountryCode user not found:%v", userId))
	gatewayId, _ := strconv.ParseUint(user.GatewayId, 10, 64)
	if vat_gateway.GetDefaultVatGateway(ctx, user.MerchantId) != nil {
		taxPercentage, _ = vat_gateway.ComputeMerchantVatPercentage(ctx, user.MerchantId, user.CountryCode, gatewayId, user.VATNumber)
		return taxPercentage, user.VATNumber, nil
	} else {
		return user.TaxPercentage, user.VATNumber, nil
	}
}
