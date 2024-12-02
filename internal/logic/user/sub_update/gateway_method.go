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
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	metric2 "unibee/internal/logic/metric"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/payment/method"
	"unibee/internal/query"
	"unibee/utility"
)

func ClearUserDefaultGatewayMethodForAutoCharge(ctx context.Context, userId uint64) {
	if userId > 0 {
		user := query.GetUserAccountById(ctx, userId)
		if user != nil && len(user.GatewayId) > 0 && len(user.PaymentMethod) > 0 {
			_, _ = dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().PaymentMethod: "",
				dao.UserAccount.Columns().GmtModify:     gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, user.Id).OmitNil().Update()
			subs := query.GetUserAllActiveOrIncompleteSubscriptions(ctx, user.Id, user.MerchantId)
			for _, sub := range subs {
				_, _ = dao.Subscription.Ctx(ctx).Data(g.Map{
					dao.Subscription.Columns().GatewayDefaultPaymentMethod: "",
					dao.Subscription.Columns().GmtModify:                   gtime.Now(),
				}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
			}
		}
	}
}

func UpdateUserDefaultGatewayForCheckout(ctx context.Context, userId uint64, gatewayId uint64) {
	if userId > 0 && gatewayId > 0 {
		user := query.GetUserAccountById(ctx, userId)
		if user != nil && len(user.GatewayId) == 0 {
			_, _ = dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().GatewayId: gatewayId,
				dao.UserAccount.Columns().GmtModify: gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, user.Id).OmitNil().Update()
		}
	}
}

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
			Topic: redismq2.TopicUserPaymentMethodChange.Topic,
			Tag:   redismq2.TopicUserPaymentMethodChange.Tag,
			Body:  strconv.FormatUint(user.Id, 10),
		})
	}
	if len(user.SubscriptionId) > 0 {
		// change user's sub gateway immediately
		_, _ = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().GatewayId:                   gatewayId,
			dao.Subscription.Columns().GatewayDefaultPaymentMethod: newPaymentMethodId,
			dao.Subscription.Columns().GmtModify:                   gtime.Now(),
		}).Where(dao.Subscription.Columns().SubscriptionId, user.SubscriptionId).OmitNil().Update()
		_, _ = redismq.Send(&redismq.Message{
			Topic: redismq2.TopicUserMetricUpdate.Topic,
			Tag:   redismq2.TopicUserMetricUpdate.Tag,
			Body: utility.MarshalToJsonString(&metric2.UserMetricUpdateMessage{
				UserId:         user.Id,
				SubscriptionId: user.SubscriptionId,
				Description:    "UpdateUserGateway",
			}),
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
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
	_, _ = redismq.Send(&redismq.Message{
		Topic:      redismq2.TopicUserAccountUpdate.Topic,
		Tag:        redismq2.TopicUserAccountUpdate.Tag,
		Body:       fmt.Sprintf("%d", user.Id),
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
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
	if reqGatewayId != nil {
		gatewayId = *reqGatewayId
		if gatewayId == userDefaultGatewayId && len(reqPaymentMethodId) == 0 {
			paymentMethodId = user.PaymentMethod
		} else {
			paymentMethodId = reqPaymentMethodId
		}
	} else if userDefaultGatewayId > 0 {
		gatewayId = userDefaultGatewayId
		paymentMethodId = user.PaymentMethod
	}
	if gatewayId <= 0 && len(subscriptionId) > 0 {
		sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
		if sub != nil {
			gatewayId = sub.GatewayId
			paymentMethodId = sub.GatewayDefaultPaymentMethod
			UpdateUserDefaultGatewayForCheckout(ctx, userId, gatewayId)
		}
	}

	return
}
