package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/customer"
	"strconv"
	"strings"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	currency2 "unibee/internal/logic/currency"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/gateway_bean"
	user2 "unibee/internal/logic/user"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type TaskActiveSubscriptionImport struct {
}

func (t TaskActiveSubscriptionImport) TaskName() string {
	return "ActiveSubscriptionImport"
}

func (t TaskActiveSubscriptionImport) TemplateHeader() interface{} {
	return &ImportActiveSubscriptionEntity{
		ExternalSubscriptionId: "exampleSubscriptionId",
		ExternalUserId:         "exampleUserId",
		ExternalPlanId:         "examplePlanId",
		Amount:                 "10.00",
		Currency:               "EUR",
		Quantity:               "1",
		Gateway:                "stripe",
		CurrentPeriodStart:     "2024-05-13 06:19:27",
		CurrentPeriodEnd:       "2024-06-13 06:19:27",
		BillingCycleAnchor:     "2024-05-13 06:19:27",
		FirstPaidTime:          "2024-05-13 06:19:27",
		CreateTime:             "2024-05-13 06:19:27",
		StripeUserId:           "",
		StripePaymentMethod:    "",
		PaypalVaultId:          "",
		Features:               "",
	}
}

func (t TaskActiveSubscriptionImport) ImportRow(ctx context.Context, task *entity.MerchantBatchTask, row map[string]string) (interface{}, error) {
	var err error
	target := &ImportActiveSubscriptionEntity{
		ExternalSubscriptionId: fmt.Sprintf("%s", row["ExternalSubscriptionId"]),
		ExternalUserId:         fmt.Sprintf("%s", row["ExternalUserId"]),
		ExternalPlanId:         fmt.Sprintf("%s", row["ExternalPlanId"]),
		Amount:                 fmt.Sprintf("%s", row["Amount"]),
		Currency:               fmt.Sprintf("%s", row["Currency"]),
		Quantity:               fmt.Sprintf("%s", row["Quantity"]),
		Gateway:                fmt.Sprintf("%s", row["Gateway(stripe|paypal|wire_transfer|changelly)"]),
		CurrentPeriodStart:     fmt.Sprintf("%s", row["CurrentPeriodStart(UTC)"]),
		CurrentPeriodEnd:       fmt.Sprintf("%s", row["CurrentPeriodEnd(UTC)"]),
		BillingCycleAnchor:     fmt.Sprintf("%s", row["BillingCycleAnchor(UTC)"]),
		FirstPaidTime:          fmt.Sprintf("%s", row["FirstPaidTime(UTC)"]),
		CreateTime:             fmt.Sprintf("%s", row["CreateTime(UTC)"]),
		StripeUserId:           fmt.Sprintf("%s", row["StripeUserId(Auto-Charge Required)"]),
		StripePaymentMethod:    fmt.Sprintf("%s", row["StripePaymentMethod(Auto-Charge Required)"]),
		PaypalVaultId:          fmt.Sprintf("%s", row["PaypalVaultId(Auto-Charge Required)"]),
		Features:               fmt.Sprintf("%s", row["Features(Json)"]),
	}
	tag := fmt.Sprintf("ImportBy%v", task.MemberId)
	if len(target.ExternalSubscriptionId) == 0 {
		return target, gerror.New("Error, ExternalSubscriptionId is blank")
	}

	// data prepare
	if len(target.ExternalUserId) == 0 {
		return target, gerror.New("Error, ExternalUserId is blank")
	}
	user := query.GetUserAccountByExternalUserId(ctx, task.MerchantId, target.ExternalUserId)
	if user == nil {
		return target, gerror.New("Error, can't find user by ExternalUserId")
	}
	if len(target.ExternalPlanId) == 0 {
		return target, gerror.New("Error, ExternalPlanId is blank")
	}
	plan := query.GetPlanByExternalPlanId(ctx, task.MerchantId, target.ExternalPlanId)
	if plan == nil {
		return target, gerror.New("Error, can't find plan by ExternalPlanId")
	}
	if len(target.Amount) == 0 {
		return target, gerror.New("Error, Amount is blank")
	}
	amountFloat, err := strconv.ParseFloat(target.Amount, 64)
	if err != nil {
		return target, gerror.Newf("Invalid Amount,error:", err.Error())
	}
	amount := int64(amountFloat * 100)
	if amount <= 0 {
		return target, gerror.New("Invalid Amount, should greater then 0")
	}
	if len(target.Currency) == 0 {
		return target, gerror.New("Error, Currency is blank")
	}
	currency := strings.TrimSpace(strings.ToUpper(target.Currency))
	if !currency2.IsCurrencySupport(currency) {
		return target, gerror.New("Error, invalid Currency")
	}
	if utility.IsNoCentCurrency(currency) {
		if amount%100 != 0 {
			return target, gerror.New("Error, this currency No decimals allowed，made it divisible by 100")
		}
	}
	if len(target.Gateway) == 0 {
		return target, gerror.New("Error, Gateway is blank")
	}
	var gatewayId uint64 = 0
	gatewayImpl := api.GatewayNameMapping[target.Gateway]
	if gatewayImpl == nil {
		return target, gerror.New("Error, Invalid Gateway, should be one of stripe|paypal|changelly|wire_transfer")
	}
	gateway := query.GetGatewayByGatewayName(ctx, task.MerchantId, target.Gateway)
	if gateway == nil {
		return target, gerror.New("Error, gateway need setup")
	}
	gatewayId = gateway.Id
	quantity, _ := strconv.ParseInt(target.Amount, 10, 64)
	if quantity == 0 {
		quantity = 1
	}
	if len(target.CurrentPeriodStart) == 0 {
		return target, gerror.New("Error, CurrentPeriodStart is blank")
	}
	currentPeriodStart := gtime.New(target.CurrentPeriodStart)
	if len(target.CurrentPeriodEnd) == 0 {
		return target, gerror.New("Error, CurrentPeriodEnd is blank")
	}
	currentPeriodEnd := gtime.New(target.CurrentPeriodEnd)

	if len(target.BillingCycleAnchor) == 0 {
		return target, gerror.New("Error, BillingCycleAnchor is blank")
	}
	billingCycleAnchor := gtime.New(target.BillingCycleAnchor)
	if len(target.FirstPaidTime) == 0 {
		return target, gerror.New("Error, FirstPaidTime is blank")
	}
	firstPaidTime := gtime.New(target.FirstPaidTime)
	if len(target.CreateTime) == 0 {
		return target, gerror.New("Error, CreateTime is blank")
	}
	createTime := gtime.New(target.CreateTime)
	// check gatewayPaymentMethod
	gatewayPaymentMethod := ""
	if len(target.PaypalVaultId) > 0 && len(target.StripePaymentMethod) > 0 {
		return target, gerror.New("Error, both PaypalVaultId and StripePaymentMethod provided")
	}
	if len(target.PaypalVaultId) > 0 && gateway.GatewayType == consts.GatewayTypePaypal {
		gatewayPaymentMethod = target.PaypalVaultId
	} else if len(target.StripePaymentMethod) > 0 && gateway.GatewayType == consts.GatewayTypeCard {
		if len(target.StripeUserId) == 0 {
			return target, gerror.New("Error, StripeUserId is blank while StripePaymentMethod is not")
		}
		listQuery, err := api.GetGatewayServiceProvider(ctx, gatewayId).GatewayUserPaymentMethodListQuery(ctx, gateway, &gateway_bean.GatewayUserPaymentMethodReq{
			UserId: user.Id,
		})
		if err != nil {
			g.Log().Errorf(ctx, "Get StripePayment MethodList error:%v", err.Error())
			return target, gerror.New("Error, can't get Stripe paymentMethod list from stripe")
		}
		found := false
		for _, method := range listQuery.PaymentMethods {
			if method.Id == target.StripePaymentMethod {
				found = true
			}
		}
		if !found {
			return target, gerror.New("Error, can't found user's paymentMethod provided from stripe ")
		}
		gatewayPaymentMethod = target.StripePaymentMethod
	}
	stripeUserId := ""
	// data verify
	{
		if currentPeriodStart.Timestamp() > gtime.Now().Timestamp() {
			return target, gerror.New("Error, CurrentPeriodStart should earlier then now")
		}
		if currentPeriodEnd.Timestamp() <= gtime.Now().Timestamp() {
			return target, gerror.New("Error, CurrentPeriodEnd should later then now")
		}
		if billingCycleAnchor.Timestamp() > gtime.Now().Timestamp() {
			return target, gerror.New("Error,BillingCycleAnchor should earlier then now")
		}
		if firstPaidTime.Timestamp() > gtime.Now().Timestamp() {
			return target, gerror.New("Error,FirstPaidTime should earlier then now")
		}
		if createTime.Timestamp() > gtime.Now().Timestamp() {
			return target, gerror.New("Error,CreateTime should earlier then now")
		}
		if len(target.StripeUserId) > 0 {
			stripeUserId = target.StripeUserId
			if gateway == nil || gateway.GatewayType != consts.GatewayTypeCard {
				return target, gerror.New("Error, gateway should be stripe while StripeUserId is not blank ")
			}
			gatewayUser := query.GetGatewayUser(ctx, user.Id, gatewayId)
			if gatewayUser != nil && gatewayUser.GatewayUserId != stripeUserId {
				return target, gerror.New("Error, There's another StripeUserId binding :" + gatewayUser.GatewayUserId)
			}
			if gatewayUser == nil {
				stripe.Key = gateway.GatewaySecret
				stripe.SetAppInfo(&stripe.AppInfo{
					Name:    "UniBee.api",
					Version: "1.0.0",
					URL:     "https://merchant.unibee.dev",
				})
				params := &stripe.CustomerParams{}
				response, err := customer.Get(stripeUserId, params)
				if err != nil {
					g.Log().Errorf(ctx, "Get StripeUserId error:%v", err.Error())
				}
				if err != nil || response == nil || len(response.ID) == 0 || response.ID != stripeUserId {
					return target, gerror.New("Error, can't get StripeUserId from stripe")
				}
				gatewayUser, err = query.CreateGatewayUser(ctx, user.Id, gatewayId, stripeUserId)
				if err != nil {
					return target, err
				}
			}
		}
	}
	one := query.GetSubscriptionByExternalSubscriptionId(ctx, target.ExternalSubscriptionId)
	if one != nil {
		if one.Data != tag {
			return target, gerror.New("Error, no permission to override," + one.Data)
		}
		_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().Amount:                      amount,
			dao.Subscription.Columns().Currency:                    currency,
			dao.Subscription.Columns().PlanId:                      plan.Id,
			dao.Subscription.Columns().Quantity:                    quantity,
			dao.Subscription.Columns().GatewayId:                   gatewayId,
			dao.Subscription.Columns().GatewayItemData:             target.Features,
			dao.Subscription.Columns().GatewayDefaultPaymentMethod: gatewayPaymentMethod,
			dao.Subscription.Columns().BillingCycleAnchor:          billingCycleAnchor.Timestamp(),
			dao.Subscription.Columns().CurrentPeriodStart:          currentPeriodStart.Timestamp(),
			dao.Subscription.Columns().CurrentPeriodEnd:            currentPeriodEnd.Timestamp(),
			dao.Subscription.Columns().CurrentPeriodStartTime:      currentPeriodStart,
			dao.Subscription.Columns().CurrentPeriodEndTime:        currentPeriodEnd,
			dao.Subscription.Columns().FirstPaidTime:               firstPaidTime.Timestamp(),
			dao.Subscription.Columns().CreateTime:                  createTime.Timestamp(),
		}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
	} else {
		one = &entity.Subscription{
			SubscriptionId:              utility.CreateSubscriptionId(),
			ExternalSubscriptionId:      target.ExternalSubscriptionId,
			UserId:                      user.Id,
			Amount:                      amount,
			Currency:                    currency,
			MerchantId:                  task.MerchantId,
			PlanId:                      plan.Id,
			Quantity:                    quantity,
			GatewayId:                   gatewayId,
			Status:                      consts.SubStatusActive,
			CurrentPeriodStart:          currentPeriodStart.Timestamp(),
			CurrentPeriodEnd:            currentPeriodEnd.Timestamp(),
			CurrentPeriodStartTime:      currentPeriodStart,
			CurrentPeriodEndTime:        currentPeriodEnd,
			BillingCycleAnchor:          billingCycleAnchor.Timestamp(),
			FirstPaidTime:               firstPaidTime.Timestamp(),
			CreateTime:                  createTime.Timestamp(),
			CountryCode:                 user.CountryCode,
			VatNumber:                   user.VATNumber,
			TaxPercentage:               user.TaxPercentage,
			GatewaySubscriptionId:       target.ExternalSubscriptionId,
			GatewayItemData:             target.Features,
			Data:                        tag,
			CurrentPeriodPaid:           1,
			GatewayDefaultPaymentMethod: gatewayPaymentMethod,
		}
		_, err = dao.Subscription.Ctx(ctx).Data(one).OmitNil().Insert(one)
	}
	if len(gatewayPaymentMethod) > 0 {
		user2.UpdateUserDefaultGatewayPaymentMethod(ctx, user.Id, gatewayId, gatewayPaymentMethod)
	}

	return target, err
}

type ImportActiveSubscriptionEntity struct {
	ExternalSubscriptionId string `json:"ExternalSubscriptionId"     `
	ExternalUserId         string `json:"ExternalUserId"     `
	ExternalPlanId         string `json:"ExternalPlanId"     `
	Amount                 string `json:"Amount"             `
	Currency               string `json:"Currency"           `
	Quantity               string `json:"Quantity"           `
	Gateway                string `json:"Gateway(stripe|paypal|wire_transfer|changelly)"            `
	CurrentPeriodStart     string `json:"CurrentPeriodStart(UTC)" `
	CurrentPeriodEnd       string `json:"CurrentPeriodEnd(UTC)"   `
	BillingCycleAnchor     string `json:"BillingCycleAnchor(UTC)" `
	FirstPaidTime          string `json:"FirstPaidTime(UTC)"      `
	CreateTime             string `json:"CreateTime(UTC)"         `
	StripeUserId           string `json:"StripeUserId(Auto-Charge Required)"             `
	StripePaymentMethod    string `json:"StripePaymentMethod(Auto-Charge Required)"      `
	PaypalVaultId          string `json:"PaypalVaultId(Auto-Charge Required)"      `
	Features               string `json:"Features(Json)"         `
}