package user

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/payment/method"
	"unibee/internal/query"
	"unibee/utility"
)

func UpdateUserDefaultSubscriptionForUpdate(ctx context.Context, userId uint64, subscriptionId string) {
	if userId > 0 && len(subscriptionId) > 0 {
		one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
		user := query.GetUserAccountById(ctx, userId)
		var subName = ""
		if one != nil && user != nil && user.SubscriptionId == subscriptionId {
			plan := query.GetPlanById(ctx, one.PlanId)
			if plan != nil {
				subName = plan.PlanName
			}
			_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().SubscriptionId:     subscriptionId,
				dao.UserAccount.Columns().SubscriptionStatus: one.Status,
				dao.UserAccount.Columns().SubscriptionName:   subName,
				dao.UserAccount.Columns().GmtModify:          gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, userId).OmitNil().Update()
			if err != nil {
				g.Log().Errorf(ctx, "UpdateUserDefaultSubscriptionForPaymentSuccess err:%s", err.Error())
			}
		}
	}
}

func UpdateUserDefaultSubscriptionForPaymentSuccess(ctx context.Context, userId uint64, subscriptionId string) {
	if userId > 0 && len(subscriptionId) > 0 {
		one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
		user := query.GetUserAccountById(ctx, userId)
		var subName = ""
		if one != nil && user != nil {
			plan := query.GetPlanById(ctx, one.PlanId)
			if plan != nil {
				subName = plan.PlanName
			}
			_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().SubscriptionId:     subscriptionId,
				dao.UserAccount.Columns().SubscriptionStatus: one.Status,
				dao.UserAccount.Columns().SubscriptionName:   subName,
				dao.UserAccount.Columns().BillingType:        1,
				dao.UserAccount.Columns().RecurringAmount:    user.RecurringAmount + one.Amount,
				dao.UserAccount.Columns().GmtModify:          gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, userId).OmitNil().Update()
			if err != nil {
				g.Log().Errorf(ctx, "UpdateUserDefaultSubscriptionForPaymentSuccess err:%s", err.Error())
			}
		}
	}
}

func UpdateUserDefaultVatNumber(ctx context.Context, userId uint64, vatNumber string) {
	if userId > 0 && len(vatNumber) > 0 {
		_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
			dao.UserAccount.Columns().VATNumber: vatNumber,
			dao.UserAccount.Columns().GmtModify: gtime.Now(),
		}).Where(dao.UserAccount.Columns().Id, userId).OmitNil().Update()
		if err != nil {
			g.Log().Errorf(ctx, "UpdateUserDefaultVatNumber err:%s", err.Error())
		}
	}
}

func UpdateUserDefaultGatewayPaymentMethod(ctx context.Context, userId uint64, gatewayId uint64, paymentMethodId string) {
	utility.Assert(userId > 0, "userId is nil")
	utility.Assert(gatewayId > 0, "gatewayId is nil")
	user := query.GetUserAccountById(ctx, userId)
	utility.Assert(user != nil, "UpdateUserDefaultGatewayPaymentMethod user not found")
	gateway := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(gateway.MerchantId == user.MerchantId, "merchant not match:"+strconv.FormatUint(gatewayId, 10))
	var newPaymentMethodId = ""
	if gateway.GatewayType == consts.GatewayTypeCard && len(paymentMethodId) > 0 {
		paymentMethod := method.QueryPaymentMethod(ctx, user.MerchantId, user.Id, gatewayId, paymentMethodId)
		utility.Assert(paymentMethod != nil, "card not found")
		newPaymentMethodId = paymentMethodId
	}
	_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
		dao.UserAccount.Columns().GatewayId:     gatewayId,
		dao.UserAccount.Columns().PaymentMethod: newPaymentMethodId,
		dao.UserAccount.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, user.Id).OmitNil().Update()
	if err != nil {
		g.Log().Errorf(ctx, "UpdateUserDefaultGatewayPaymentMethod userId:%d gatewayId:%d, paymentMethodId:%s error:%s", userId, gatewayId, paymentMethodId, err.Error())
	} else {
		g.Log().Errorf(ctx, "UpdateUserDefaultGatewayPaymentMethod userId:%d gatewayId:%d, paymentMethodId:%s success", userId, gatewayId, paymentMethodId)
	}
}
