package sub_update

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	redismq2 "unibee/internal/cmd/redismq"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/vat_gateway"
	"unibee/internal/query"
	"unibee/utility"
)

func UpdateUserVatNumber(ctx context.Context, userId uint64, vatNumber string) {
	utility.Assert(userId > 0, "userId is nil")
	user := query.GetUserAccountById(ctx, userId)
	utility.Assert(user != nil, "UpdateUserCountryCode user not found")
	if user.VATNumber == vatNumber {
		return
	}
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
					Content:        fmt.Sprintf("UpdateVATNumber(%s)", vatNumber),
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
			Content:        "Clear(VatNumber)",
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
	}
}

func UpdateUserDefaultVatNumber(ctx context.Context, userId uint64, vatNumber string) {
	if userId > 0 && len(vatNumber) > 0 {
		user := query.GetUserAccountById(ctx, userId)
		if user == nil {
			return
		}
		_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
			dao.UserAccount.Columns().VATNumber: vatNumber,
			dao.UserAccount.Columns().GmtModify: gtime.Now(),
		}).Where(dao.UserAccount.Columns().Id, userId).OmitNil().Update()
		if err != nil {
			g.Log().Errorf(ctx, "UpdateUserDefaultVatNumber err:%s", err.Error())
		}

		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     user.MerchantId,
			Target:         fmt.Sprintf("User(%v)", user.Id),
			Content:        fmt.Sprintf("UpdateVATNumber(%s)", vatNumber),
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
	}
}
