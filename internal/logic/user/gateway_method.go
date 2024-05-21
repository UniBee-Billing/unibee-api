package user

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"strconv"
	config2 "unibee/internal/cmd/config"
	"unibee/internal/query"
	"unibee/utility"
)

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
