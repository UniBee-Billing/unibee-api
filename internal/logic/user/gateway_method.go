package user

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	config2 "unibee/internal/cmd/config"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/payment/method"
	"unibee/internal/query"
	"unibee/utility"
)

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

func VerifyPaymentGatewayMethod(ctx context.Context, userId uint64, reqGatewayId *uint64, reqPaymentMethodId string, subscriptionId string) (gatewayId uint64, paymentMethodId string) {
	user := query.GetUserAccountById(ctx, userId)
	utility.Assert(user != nil, fmt.Sprintf("user not found:%d", userId))
	var userDefaultGatewayId uint64 = 0
	var err error = nil
	if len(user.GatewayId) > 0 {
		userDefaultGatewayId, err = strconv.ParseUint(user.GatewayId, 10, 64)
		if err != nil {
			g.Log().Errorf(ctx, "ParseUserDefaultMethod:%d", user.GatewayId)
			return
		}
	}
	if len(reqPaymentMethodId) > 0 {
		utility.Assert(reqGatewayId != nil, "gateway need specified")
		// todo mark check reqPaymentMethodId valid
	}
	if userDefaultGatewayId > 0 && reqGatewayId == nil {
		gatewayId = userDefaultGatewayId
		paymentMethodId = user.PaymentMethod
	} else if reqGatewayId != nil {
		gatewayId = *reqGatewayId
		if gatewayId == userDefaultGatewayId && len(reqPaymentMethodId) == 0 {
			paymentMethodId = user.PaymentMethod
		} else {
			paymentMethodId = reqPaymentMethodId
		}
	}
	utility.Assert(gatewayId > 0, "gateway need specified")
	if !config2.GetConfigInstance().IsProd() {
		if len(paymentMethodId) == 0 && len(subscriptionId) > 0 {
			sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
			if sub != nil && sub.GatewayId == gatewayId {
				paymentMethodId = sub.GatewayDefaultPaymentMethod
			}
		}
	}
	return
}