package sub_update

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"strconv"
	"strings"
	config2 "unibee/internal/cmd/config"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/payment/method"
	"unibee/internal/query"
	"unibee/utility"
)

func UpdateUserDefaultGatewayPaymentMethod(ctx context.Context, userId uint64, gatewayId uint64, paymentMethodId string) {
	g.Log().Infof(ctx, "UpdateUserDefaultGatewayPaymentMethod userId:%v gatewayId:%s paymentMethod:%s", userId, gatewayId, paymentMethodId)
	utility.Assert(userId > 0, "userId is nil")
	utility.Assert(gatewayId > 0, "gatewayId is nil")
	user := query.GetUserAccountById(ctx, userId)
	utility.Assert(user != nil, "UpdateUserDefaultGatewayPaymentMethod user not found")
	gateway := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(gateway.MerchantId == user.MerchantId, "merchant not match:"+strconv.FormatUint(gatewayId, 10))
	if user.GatewayId == fmt.Sprintf("%v", gatewayId) && user.PaymentMethod == paymentMethodId {
		return
	}
	var newPaymentMethodId = ""
	if gateway.GatewayType == consts.GatewayTypeCard && len(paymentMethodId) > 0 {
		paymentMethod := method.QueryPaymentMethod(ctx, user.MerchantId, user.Id, gatewayId, paymentMethodId)
		utility.Assert(paymentMethod != nil, "card not found")
		newPaymentMethodId = paymentMethodId
	} else if gateway.GatewayType == consts.GatewayTypePaypal && len(paymentMethodId) > 0 {
		newPaymentMethodId = paymentMethodId
	}
	gatewayUser := query.GetGatewayUser(ctx, userId, gatewayId)
	if len(newPaymentMethodId) == 0 && gatewayUser != nil && len(gatewayUser.GatewayDefaultPaymentMethod) > 0 {
		newPaymentMethodId = gatewayUser.GatewayDefaultPaymentMethod
	}
	_, err := dao.UserAccount.Ctx(ctx).Data(g.Map{
		dao.UserAccount.Columns().GatewayId:     gatewayId,
		dao.UserAccount.Columns().PaymentMethod: newPaymentMethodId,
		dao.UserAccount.Columns().GmtModify:     gtime.Now(),
	}).Where(dao.UserAccount.Columns().Id, user.Id).OmitNil().Update()
	if err != nil {
		g.Log().Errorf(ctx, "UpdateUserDefaultGatewayPaymentMethod userId:%d gatewayId:%d, paymentMethodId:%s error:%s", userId, gatewayId, paymentMethodId, err.Error())
	} else {
		g.Log().Infof(ctx, "UpdateUserDefaultGatewayPaymentMethod userId:%d gatewayId:%d, paymentMethodId:%s success", userId, gatewayId, paymentMethodId)
	}
	var oldGatewayId uint64 = 0
	if len(user.GatewayId) > 0 {
		oldGatewayId, _ = strconv.ParseUint(user.GatewayId, 10, 64)
	}
	if len(newPaymentMethodId) > 0 && gatewayUser != nil && strings.Compare(gatewayUser.GatewayDefaultPaymentMethod, newPaymentMethodId) != 0 {
		_, _ = query.CreateOrUpdateGatewayUser(ctx, userId, gatewayId, gatewayUser.GatewayUserId, newPaymentMethodId)
	}
	if oldGatewayId != gatewayId || user.PaymentMethod != paymentMethodId {
		_, _ = redismq.Send(&redismq.Message{
			Topic: redismq2.TopicUserPaymentMethodChanged.Topic,
			Tag:   redismq2.TopicUserPaymentMethodChanged.Tag,
			Body:  strconv.FormatUint(user.Id, 10),
		})
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     user.MerchantId,
		Target:         fmt.Sprintf("User(%v)", user.Id),
		Content:        fmt.Sprintf("ChangeGateway(%s)", gateway.GatewayName),
		UserId:         user.Id,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
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
		utility.Assert(reqGatewayId != nil, "gateway need specified while payment method not empty")
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
	if !config2.GetConfigInstance().IsProd() {
		if gatewayId > 0 && len(paymentMethodId) == 0 && len(subscriptionId) > 0 {
			sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
			if sub != nil && sub.GatewayId == gatewayId {
				paymentMethodId = sub.GatewayDefaultPaymentMethod
			}
		}
	}
	return
}
