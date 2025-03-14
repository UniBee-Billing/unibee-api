package sub_update

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"strconv"
	"strings"
	redismq2 "unibee/internal/cmd/redismq"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/vat_gateway"
	"unibee/internal/query"
	"unibee/utility"
)

func UpdateUserTaxPercentageOnly(ctx context.Context, userId uint64, taxPercentage int64) {
	_, _ = dao.UserAccount.Ctx(ctx).Data(g.Map{
		dao.UserAccount.Columns().TaxPercentage: taxPercentage,
		dao.UserAccount.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, userId).Update()
}

func UpdateUserCountryCode(ctx context.Context, userId uint64, countryCode string) {
	utility.Assert(userId > 0, "userId is nil")
	user := query.GetUserAccountById(ctx, userId)
	utility.Assert(user != nil, "UpdateUserCountryCode user not found")
	if user.CountryCode == countryCode {
		return
	}
	if len(countryCode) > 0 && strings.Compare(user.CountryCode, countryCode) != 0 {
		taxPercentage := user.TaxPercentage
		countryName := user.CountryName
		gatewayId, _ := strconv.ParseUint(user.GatewayId, 10, 64)
		if vat_gateway.GetDefaultVatGateway(ctx, user.MerchantId).VatRatesEnabled() {
			taxPercentage, countryName = vat_gateway.ComputeMerchantVatPercentage(ctx, user.MerchantId, countryCode, gatewayId, user.VATNumber)
		} else {
			countryOne, _ := vat_gateway.QueryVatCountryRateByMerchant(ctx, user.MerchantId, countryCode)
			if countryOne != nil {
				countryName = countryOne.CountryName
			}
		}
		_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
			dao.UserAccount.Columns().CountryCode:   countryCode,
			dao.UserAccount.Columns().CountryName:   countryName,
			dao.UserAccount.Columns().TaxPercentage: taxPercentage,
			dao.UserAccount.Columns().GmtModify:     gtime.Now(),
		}).Where(dao.UserAccount.Columns().Id, user.Id).OmitNil().Update()
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     user.MerchantId,
			Target:         fmt.Sprintf("User(%v)", user.Id),
			Content:        fmt.Sprintf("UpdateCountryCode(%s-%v)", countryCode, taxPercentage),
			UserId:         user.Id,
			SubscriptionId: "",
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, nil)
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicUserAccountUpdate.Topic,
			Tag:        redismq2.TopicUserAccountUpdate.Tag,
			Body:       fmt.Sprintf("%d", user.Id),
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
		if err != nil {
			g.Log().Errorf(ctx, "UpdateUserCountryCode userId:%d CountryCode:%s, error:%s", userId, countryCode, err.Error())
		} else {
			g.Log().Infof(ctx, "UpdateUserCountryCode userId:%d CountryCode:%s, success", userId, countryCode)
		}
	}
}
