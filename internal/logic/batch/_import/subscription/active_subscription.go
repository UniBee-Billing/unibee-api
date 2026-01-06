package subscription

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/gateway/api"
	handler2 "unibee/internal/logic/invoice/handler"
	service3 "unibee/internal/logic/invoice/service"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/plan/period"
	"unibee/internal/logic/subscription/timeline"
	user2 "unibee/internal/logic/user"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type TaskActiveSubscriptionImport struct {
}

func (t TaskActiveSubscriptionImport) TemplateVersion() string {
	return "v1"
}

func (t TaskActiveSubscriptionImport) TaskName() string {
	return "ActiveSubscriptionImport"
}

func (t TaskActiveSubscriptionImport) TemplateHeader() interface{} {
	return &ImportActiveSubscriptionEntity{
		ExternalSubscriptionId: "exampleSubscriptionId",
		ExternalUserId:         "exampleUserId",
		ExternalPlanId:         "examplePlanId",
		ExpectedTotalAmount:    "10.00",
		Quantity:               "1",
		Gateway:                "stripe",
		CountryCode:            "EE",
		CurrentPeriodStart:     "2024-05-13 06:19:27",
		CurrentPeriodEnd:       "2024-06-13 06:19:27",
		BillingCycleAnchor:     "2024-05-13 06:19:27",
		FirstPaidTime:          "2024-05-13 06:19:27",
		CreateTime:             "2024-05-13 06:19:27",
		//StripeUserId:           "",
		//StripePaymentMethod:    "",
		//PaypalVaultId:          "",
		Features: "",
	}
}

func (t TaskActiveSubscriptionImport) ImportRow(ctx context.Context, task *entity.MerchantBatchTask, row map[string]string) (interface{}, error) {
	var err error
	target := &ImportActiveSubscriptionEntity{
		ExternalSubscriptionId: fmt.Sprintf("%s", row["ExternalSubscriptionId"]),
		ExternalUserId:         fmt.Sprintf("%s", row["ExternalUserId"]),
		ExternalPlanId:         fmt.Sprintf("%s", row["ExternalPlanId"]),
		ExpectedTotalAmount:    fmt.Sprintf("%s", row["Amount"]),
		Quantity:               fmt.Sprintf("%s", row["Quantity"]),
		Gateway:                fmt.Sprintf("%s", row["Gateway"]),
		CurrentPeriodStart:     fmt.Sprintf("%s", row["CurrentPeriodStart"]),
		CurrentPeriodEnd:       fmt.Sprintf("%s", row["CurrentPeriodEnd"]),
		BillingCycleAnchor:     fmt.Sprintf("%s", row["BillingCycleAnchor"]),
		FirstPaidTime:          fmt.Sprintf("%s", row["FirstPaidTime"]),
		CreateTime:             fmt.Sprintf("%s", row["CreateTime"]),
		CountryCode:            fmt.Sprintf("%s", row["CountryCode"]),
		VatNumber:              fmt.Sprintf("%s", row["VatNumber"]),
		TaxPercentage:          fmt.Sprintf("%s", row["TaxPercentage"]),
		//StripeUserId:           fmt.Sprintf("%s", row["StripeUserId(Auto-Charge Required)"]),
		//StripePaymentMethod:    fmt.Sprintf("%s", row["StripePaymentMethod(Auto-Charge Required)"]),
		//PaypalVaultId:          fmt.Sprintf("%s", row["PaypalVaultId(Auto-Charge Required)"]),
		Features: fmt.Sprintf("%s", row["Features"]),
	}
	tag := fmt.Sprintf("ImportBy%v", task.MemberId)
	if len(target.ExternalSubscriptionId) == 0 {
		return target, gerror.New("Error, ExternalSubscriptionId is blank")
	}
	target.CountryCode = strings.ToUpper(target.CountryCode)
	if len(target.CountryCode) > 0 {
		err = utility.ValidateCountryCode(target.CountryCode)
		if err != nil {
			return target, gerror.New(fmt.Sprintf("Error, CountryCode is invalid, %s", err.Error()))
		}
	}
	// data prepare
	user, err := user2.QueryOrCreateUser(ctx, &user2.NewUserInternalReq{
		ExternalUserId: target.ExternalUserId,
		Email:          target.Email,
		CountryCode:    target.CountryCode,
		VATNumber:      target.VatNumber,
		MerchantId:     _interface.GetMerchantId(ctx),
	})
	if err != nil {
		return target, gerror.Newf("QueryOrCreateUser,error:%s", err.Error())
	}
	if user == nil {
		return target, gerror.New("Error, can't find user by ExternalUserId")
	}
	taxPercentage := user.TaxPercentage
	if !vat_gateway.GetDefaultVatGateway(ctx, user.MerchantId).VatRatesEnabled() && len(target.TaxPercentage) > 0 {
		taxPercentage = 0
		taxPercentageFloat, err := strconv.ParseFloat(target.TaxPercentage, 64)
		if err == nil {
			taxPercentage = int64(taxPercentageFloat * 10000)
		}
	}
	if len(target.ExternalPlanId) == 0 {
		return target, gerror.New("Error, ExternalPlanId is blank")
	}
	plan := query.GetPlanByExternalPlanId(ctx, task.MerchantId, target.ExternalPlanId)
	if plan == nil {
		return target, gerror.New("Error, can't find plan by ExternalPlanId")
	}
	utility.Assert(plan.Status != consts.PlanStatusEditable, "plan status should not editable")
	quantity, _ := strconv.ParseInt(target.Quantity, 10, 64)
	if quantity == 0 {
		quantity = 1
	}
	totalAmountExcludingTax := plan.Amount * quantity
	var taxAmount = int64(math.Round(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(taxPercentage)))
	totalAmount := totalAmountExcludingTax + taxAmount
	if len(target.ExpectedTotalAmount) == 0 {
		expectedTotalAmountFloat, err := strconv.ParseFloat(target.ExpectedTotalAmount, 64)
		if err != nil {
			return target, gerror.Newf("Invalid Amount,error:%s", err.Error())
		}
		expectedTotalAmount := int64(expectedTotalAmountFloat * 100)
		if expectedTotalAmount > 0 && expectedTotalAmount != totalAmount {
			return target, gerror.New(fmt.Sprintf("ExpectedTotalAmount Verify failed, calculated totalAmount %s != expectedTotalAmount %s", utility.ConvertCentToDollarStr(totalAmount, plan.Currency), utility.ConvertCentToDollarStr(expectedTotalAmount, plan.Currency)))
		}
	}
	if len(target.Gateway) == 0 {
		return target, gerror.New("Error, Gateway is blank")
	}
	var gatewayId uint64 = 0
	gatewayImpl := api.GatewayNameMapping[target.Gateway]
	if gatewayImpl == nil {
		return target, gerror.New("Error, Invalid Gateway, should be one of " + strings.Join(api.ExportGatewaySetupListKeys(), "|"))
	}
	gateway := query.GetDefaultGatewayByGatewayName(ctx, task.MerchantId, target.Gateway)
	if gateway == nil {
		return target, gerror.New("Error, gateway need setup")
	}
	gatewayId = gateway.Id

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
	//if len(target.PaypalVaultId) > 0 && len(target.StripePaymentMethod) > 0 {
	//	return target, gerror.New("Error, both PaypalVaultId and StripePaymentMethod provided")
	//}
	//if len(target.PaypalVaultId) > 0 && gateway.GatewayType == consts.GatewayTypePaypal {
	//	gatewayPaymentMethod = target.PaypalVaultId
	//	// todo mark check paypal vaultId
	//} else if len(target.StripePaymentMethod) > 0 && gateway.GatewayType == consts.GatewayTypeCard {
	//	if len(target.StripeUserId) == 0 {
	//		return target, gerror.New("Error, StripeUserId is blank while StripePaymentMethod is not")
	//	}
	//	listQuery, err := api.GetGatewayServiceProvider(ctx, gatewayId).GatewayUserPaymentMethodListQuery(ctx, gateway, &gateway_bean.GatewayUserPaymentMethodReq{
	//		UserId:        user.Id,
	//		GatewayUserId: target.StripeUserId,
	//	})
	//	if err != nil {
	//		g.Log().Errorf(ctx, "Get StripePayment MethodList error:%v", err.Error())
	//		return target, gerror.New("Error, can't get Stripe paymentMethod list from stripe")
	//	}
	//	found := false
	//	for _, method := range listQuery.PaymentMethods {
	//		if method.Id == target.StripePaymentMethod {
	//			found = true
	//		}
	//	}
	//	if !found {
	//		return target, gerror.New("Error, can't found user's paymentMethod provided from stripe ")
	//	}
	//	gatewayPaymentMethod = target.StripePaymentMethod
	//}
	//stripeUserId := ""
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
		if currentPeriodStart.Timestamp() < createTime.Timestamp() || currentPeriodStart.Timestamp() < billingCycleAnchor.Timestamp() {
			return target, gerror.New("Error,currentPeriodStart should later then createTime and billingCycleAnchor")
		}
		if currentPeriodEnd.Timestamp() <= currentPeriodStart.Timestamp() ||
			currentPeriodEnd.Timestamp() <= billingCycleAnchor.Timestamp() ||
			currentPeriodEnd.Timestamp() <= firstPaidTime.Timestamp() ||
			currentPeriodEnd.Timestamp() <= createTime.Timestamp() {
			return target, gerror.New("Error,currentPeriodEnd should later then currentPeriodStart,firstPaidTime,createTime and billingCycleAnchor")
		}

		//if len(target.StripeUserId) > 0 {
		//	stripeUserId = target.StripeUserId
		//	if gateway.GatewayType != consts.GatewayTypeCard {
		//		return target, gerror.New("Error, gateway should be stripe while StripeUserId is not blank ")
		//	}
		//	gatewayUser := util.GetGatewayUser(ctx, user.Id, gatewayId)
		//	if gatewayUser != nil && gatewayUser.GatewayUserId != stripeUserId {
		//		// todo mark may delete the old one
		//		return target, gerror.New("Error, There's another StripeUserId binding :" + gatewayUser.GatewayUserId)
		//	}
		//	if gatewayUser == nil {
		//		stripe.Key = gateway.GatewaySecret
		//		stripe.SetAppInfo(&stripe.AppInfo{
		//			Name:    "UniBee.api",
		//			Version: "1.0.0",
		//			URL:     "https://merchant.unibee.dev",
		//		})
		//		params := &stripe.CustomerParams{}
		//		response, err := customer.Get(stripeUserId, params)
		//		if err != nil {
		//			g.Log().Errorf(ctx, "Get StripeUserId error:%v", err.Error())
		//		}
		//		if err != nil || response == nil || len(response.ID) == 0 || response.ID != stripeUserId {
		//			return target, gerror.New("Error, can't get StripeUserId from stripe")
		//		}
		//		//// todo mark verify email from stripe
		//		//if response.Email != user.Email {
		//		//	return target, gerror.New("Error, stripe customer email not equal user's email")
		//		//}
		//		gatewayUser, err = util.CreateGatewayUser(ctx, user.Id, gatewayId, stripeUserId)
		//		if err != nil {
		//			return target, err
		//		}
		//	}
		//}
	}

	metadata := make(map[string]interface{})
	if _interface.Context().Get(ctx).MerchantMember != nil {
		tag = fmt.Sprintf("ImportByMember:%d", _interface.Context().Get(ctx).MerchantMember.Id)
		metadata["ImportFrom"] = tag
	} else {
		tag = fmt.Sprintf("ImportByOpenAPI")
		metadata["ImportFrom"] = tag
	}

	var dunningTime = period.GetDunningTimeFromEnd(ctx, utility.MaxInt64(currentPeriodEnd.Timestamp(), 0), plan.Id)
	one := query.GetSubscriptionByExternalSubscriptionId(ctx, target.ExternalSubscriptionId)
	override := false
	if one != nil {
		if one.Data != tag {
			return target, gerror.New("Error, no permission to override," + one.Data)
		}
		if one.UserId != user.Id {
			return target, gerror.New("Error, no permission to override, user not match")
		}
		_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().Type:                        consts.SubTypeUniBeeControl,
			dao.Subscription.Columns().Status:                      consts.SubStatusActive,
			dao.Subscription.Columns().Amount:                      totalAmount,
			dao.Subscription.Columns().Currency:                    plan.Currency,
			dao.Subscription.Columns().PlanId:                      plan.Id,
			dao.Subscription.Columns().Quantity:                    quantity,
			dao.Subscription.Columns().GatewayId:                   gatewayId,
			dao.Subscription.Columns().TaxPercentage:               taxPercentage,
			dao.Subscription.Columns().GatewayItemData:             target.Features,
			dao.Subscription.Columns().GatewayDefaultPaymentMethod: gatewayPaymentMethod,
			dao.Subscription.Columns().BillingCycleAnchor:          billingCycleAnchor.Timestamp(),
			dao.Subscription.Columns().CurrentPeriodStart:          currentPeriodStart.Timestamp(),
			dao.Subscription.Columns().CurrentPeriodEnd:            currentPeriodEnd.Timestamp(),
			dao.Subscription.Columns().DunningTime:                 dunningTime,
			dao.Subscription.Columns().CurrentPeriodStartTime:      currentPeriodStart,
			dao.Subscription.Columns().CurrentPeriodEndTime:        currentPeriodEnd,
			dao.Subscription.Columns().FirstPaidTime:               firstPaidTime.Timestamp(),
			dao.Subscription.Columns().CreateTime:                  createTime.Timestamp(),
			dao.Subscription.Columns().MetaData:                    utility.MarshalToJsonString(metadata),
		}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
		override = true
		utility.AssertError(err, "Override history error")

		{
			if len(one.LatestInvoiceId) == 0 {
				currentInvoice := &bean.Invoice{
					InvoiceName:                    "SubscriptionCreate",
					BizType:                        consts.BizTypeSubscription,
					ProductName:                    plan.PlanName,
					OriginAmount:                   0,
					TotalAmount:                    0,
					TotalAmountExcludingTax:        0,
					DiscountCode:                   "",
					DiscountAmount:                 0,
					Currency:                       one.Currency,
					TaxAmount:                      0,
					SubscriptionAmount:             0,
					SubscriptionAmountExcludingTax: 0,
					Lines: []*bean.InvoiceItemSimplify{{
						Currency:               one.Currency,
						OriginAmount:           0,
						Amount:                 0,
						DiscountAmount:         0,
						Tax:                    0,
						AmountExcludingTax:     0,
						TaxPercentage:          0,
						UnitAmountExcludingTax: 0,
						Name:                   plan.PlanName,
						Description:            plan.Description,
						Proration:              false,
						Quantity:               quantity,
						PeriodEnd:              currentPeriodEnd.Timestamp(),
						PeriodStart:            currentPeriodStart.Timestamp(),
						Plan:                   bean.SimplifyPlan(plan),
					}},
					ProrationDate: time.Now().Unix(),
					PeriodStart:   one.CurrentPeriodStart,
					PeriodEnd:     one.CurrentPeriodEnd,
					Metadata:      metadata,
					CountryCode:   target.CountryCode,
					VatNumber:     target.VatNumber,
					TaxPercentage: 0,
				}
				invoice, err := service3.CreateProcessingInvoiceForSub(ctx, &service3.CreateProcessingInvoiceForSubReq{
					PlanId:             plan.Id,
					Simplify:           currentInvoice,
					Sub:                one,
					GatewayId:          one.GatewayId,
					IsSubLatestInvoice: true,
					TimeNow:            currentInvoice.ProrationDate,
				})
				utility.AssertError(err, "Create Latest Invoice Error")
				invoice, err = handler2.MarkInvoiceAsPaidForZeroPayment(ctx, invoice.InvoiceId)
				utility.AssertError(err, "Create Latest Invoice Error")
				timeline.SubscriptionNewTimeline(ctx, invoice)
				sub_update.UpdateUserDefaultSubscriptionForPaymentSuccess(ctx, one.UserId, one.SubscriptionId)
			}
		}

		subscription3.SendMerchantSubscriptionWebhookBackground(one, -10000, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_IMPORT_OVERRIDE, map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()})

		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     one.MerchantId,
			Target:         fmt.Sprintf("Subscription(%s)", one.SubscriptionId),
			Content:        "ImportOverride",
			UserId:         one.UserId,
			SubscriptionId: one.SubscriptionId,
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)
	} else {
		one = &entity.Subscription{
			Type:                        consts.SubTypeUniBeeControl,
			SubscriptionId:              utility.CreateSubscriptionId(),
			ExternalSubscriptionId:      target.ExternalSubscriptionId,
			UserId:                      user.Id,
			Amount:                      totalAmount,
			Currency:                    plan.Currency,
			MerchantId:                  task.MerchantId,
			PlanId:                      plan.Id,
			Quantity:                    quantity,
			GatewayId:                   gatewayId,
			Status:                      consts.SubStatusActive,
			CurrentPeriodStart:          currentPeriodStart.Timestamp(),
			CurrentPeriodEnd:            currentPeriodEnd.Timestamp(),
			CurrentPeriodStartTime:      currentPeriodStart,
			CurrentPeriodEndTime:        currentPeriodEnd,
			DunningTime:                 dunningTime,
			BillingCycleAnchor:          billingCycleAnchor.Timestamp(),
			FirstPaidTime:               firstPaidTime.Timestamp(),
			CreateTime:                  createTime.Timestamp(),
			CountryCode:                 user.CountryCode,
			VatNumber:                   user.VATNumber,
			TaxPercentage:               taxPercentage,
			GatewaySubscriptionId:       target.ExternalSubscriptionId,
			GatewayItemData:             target.Features,
			Data:                        tag,
			CurrentPeriodPaid:           1,
			GatewayDefaultPaymentMethod: gatewayPaymentMethod,
			MetaData:                    utility.MarshalToJsonString(metadata),
		}
		result, err := dao.Subscription.Ctx(ctx).Data(one).OmitNil().Insert(one)
		utility.AssertError(err, "Save history error")
		id, err := result.LastInsertId()
		one.Id = uint64(id)

		{
			currentInvoice := &bean.Invoice{
				InvoiceName:                    "SubscriptionCreate",
				BizType:                        consts.BizTypeSubscription,
				ProductName:                    plan.PlanName,
				OriginAmount:                   0,
				TotalAmount:                    0,
				TotalAmountExcludingTax:        0,
				DiscountCode:                   "",
				DiscountAmount:                 0,
				Currency:                       one.Currency,
				TaxAmount:                      0,
				SubscriptionAmount:             0,
				SubscriptionAmountExcludingTax: 0,
				Lines: []*bean.InvoiceItemSimplify{{
					Currency:               one.Currency,
					OriginAmount:           0,
					Amount:                 0,
					DiscountAmount:         0,
					Tax:                    0,
					AmountExcludingTax:     0,
					TaxPercentage:          0,
					UnitAmountExcludingTax: 0,
					Name:                   plan.PlanName,
					Description:            plan.Description,
					Proration:              false,
					Quantity:               quantity,
					PeriodEnd:              currentPeriodEnd.Timestamp(),
					PeriodStart:            currentPeriodStart.Timestamp(),
					Plan:                   bean.SimplifyPlan(plan),
				}},
				ProrationDate: time.Now().Unix(),
				PeriodStart:   one.CurrentPeriodStart,
				PeriodEnd:     one.CurrentPeriodEnd,
				Metadata:      metadata,
				CountryCode:   target.CountryCode,
				VatNumber:     target.VatNumber,
				TaxPercentage: 0,
			}
			invoice, err := service3.CreateProcessingInvoiceForSub(ctx, &service3.CreateProcessingInvoiceForSubReq{
				PlanId:             plan.Id,
				Simplify:           currentInvoice,
				Sub:                one,
				GatewayId:          one.GatewayId,
				IsSubLatestInvoice: true,
				TimeNow:            currentInvoice.ProrationDate,
			})
			utility.AssertError(err, "Create Latest Invoice Error")
			invoice, err = handler2.MarkInvoiceAsPaidForZeroPayment(ctx, invoice.InvoiceId)
			utility.AssertError(err, "Create Latest Invoice Error")
			timeline.SubscriptionNewTimeline(ctx, invoice)
			sub_update.UpdateUserDefaultSubscriptionForPaymentSuccess(ctx, one.UserId, one.SubscriptionId)
		}

		subscription3.SendMerchantSubscriptionWebhookBackground(one, -10000, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_IMPORT_CREATED, map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()})
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     one.MerchantId,
			Target:         fmt.Sprintf("Subscription(%s)", one.SubscriptionId),
			Content:        "ImportNew",
			UserId:         one.UserId,
			SubscriptionId: one.SubscriptionId,
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)
	}
	//if len(gatewayPaymentMethod) > 0 {
	//	user2.UpdateUserDefaultGatewayPaymentMethod(ctx, user.Id, gatewayId, gatewayPaymentMethod)
	//}
	if err == nil && override {
		err = gerror.New("override success")
	}
	return target, err
}

type ImportActiveSubscriptionEntity struct {
	ExternalSubscriptionId string `json:"ExternalSubscriptionId"    comment:"Required, The external id of subscription"     `
	ExternalPlanId         string `json:"ExternalPlanId"   comment:"Required, The external id of plan, plan should created at first"   `
	Email                  string `json:"Email"  comment:"The email of user, one of Email or ExternalUserId is required" `
	ExternalUserId         string `json:"ExternalUserId"    comment:"The external id of user, one of Email or ExternalUserId is required "    `
	ExpectedTotalAmount    string `json:"ExpectedTotalAmount" comment:"Optional. If greater than 0, the system will verify the calculated total amount against this value (e.g., 19.99 = 19.99 USD)"`
	Quantity               string `json:"Quantity"      comment:"the quantity of plan, default 1 if not provided "        `
	CountryCode            string `json:"CountryCode"    comment:"Required, The country code of subscription, Tax is applied based on this country code."  `
	VatNumber              string `json:"VatNumber"    comment:"The Vat Number of user"  `
	TaxPercentage          string `json:"TaxPercentage" comment:"The TaxPercentage of user, valid when system vat gateway enabled, em. 10 = 10%"`
	Gateway                string `json:"Gateway" comment:"Required, should one of stripe|paypal|wire_transfer|changelly "           `
	CurrentPeriodStart     string `json:"CurrentPeriodStart" comment:"Required, UTC time, the current period start time of subscription, format '2006-01-02 15:04:05'"`
	CurrentPeriodEnd       string `json:"CurrentPeriodEnd"   comment:"Required, UTC time, the current period end time of subscription, format '2006-01-02 15:04:05'"`
	BillingCycleAnchor     string `json:"BillingCycleAnchor"   comment:"Required, UTC time, The reference point that aligns future billing cycle dates. It sets the day of week for week intervals, the day of month for month and year intervals, and the month of year for year intervals, format '2006-01-02 15:04:05'"`
	FirstPaidTime          string `json:"FirstPaidTime"   comment:"UTC time, the first payment success time of subscription, format '2006-01-02 15:04:05'"   `
	CreateTime             string `json:"CreateTime"      comment:"Required, UTC time, the creation time of subscription, format '2006-01-02 15:04:05'"   `
	//StripeUserId           string `json:"StripeUserId(Auto-Charge Required)"      comment:"The id of user get from stripe, required if stripe auto-charge needed"       `
	//StripePaymentMethod    string `json:"StripePaymentMethod(Auto-Charge Required)"     comment:"The payment method id which user attached, get from stripe, required if stripe auto-charge needed"    `
	//PaypalVaultId          string `json:"PaypalVaultId(Auto-Charge Required)"    comment:"The vault id of user get from paypal, required if paypal auto-charge needed"   `
	Features string `json:"Features"    comment:"In json format, additional features data of subscription, will join user's metric data in user api if provided'"     `
}
