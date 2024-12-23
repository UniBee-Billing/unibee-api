package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"strings"
	"unibee/api/bean"
	"unibee/api/user/vat"
	config2 "unibee/internal/cmd/config"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/email"
	"unibee/internal/logic/invoice/invoice_compute"
	service3 "unibee/internal/logic/invoice/service"
	metric2 "unibee/internal/logic/metric"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/payment/service"
	subscription2 "unibee/internal/logic/subscription"
	"unibee/internal/logic/subscription/handler"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func checkAndListAddonsFromParams(ctx context.Context, addonParams []*bean.PlanAddonParam) []*bean.PlanAddonDetail {
	var addons []*bean.PlanAddonDetail
	var totalAddonIds []uint64
	if len(addonParams) > 0 {
		for _, s := range addonParams {
			totalAddonIds = append(totalAddonIds, s.AddonPlanId)
		}
	}
	var allAddonList []*entity.Plan
	if len(totalAddonIds) > 0 {
		//query all plan
		err := dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, totalAddonIds).OmitEmpty().Scan(&allAddonList)
		if err == nil {
			//add to list
			mapPlans := make(map[uint64]*entity.Plan)
			for _, pair := range allAddonList {
				key := pair.Id
				value := pair
				mapPlans[key] = value
			}
			for _, param := range addonParams {
				utility.Assert(mapPlans[param.AddonPlanId] != nil, fmt.Sprintf("AddonPlanId not found:%v", param.AddonPlanId))
				utility.Assert(mapPlans[param.AddonPlanId].Type == consts.PlanTypeRecurringAddon, fmt.Sprintf("Id:%v not Addon Type", param.AddonPlanId))
				utility.Assert(mapPlans[param.AddonPlanId].IsDeleted == 0, fmt.Sprintf("Addon Id:%v is Deleted", param.AddonPlanId))
				utility.Assert(param.Quantity > 0, fmt.Sprintf("Id:%v quantity invalid", param.AddonPlanId))
				addons = append(addons, &bean.PlanAddonDetail{
					Quantity:  param.Quantity,
					AddonPlan: bean.SimplifyPlan(mapPlans[param.AddonPlanId]),
				})
			}
		}
	}
	return addons
}

func VatNumberValidate(ctx context.Context, req *vat.NumberValidateReq) (*vat.NumberValidateRes, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(len(req.VatNumber) > 0, "vatNumber invalid")
	vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, _interface.GetMerchantId(ctx), 0, req.VatNumber, "")
	if err != nil {
		return nil, err
	}
	return &vat.NumberValidateRes{VatNumberValidate: vatNumberValidate}, nil
}

func GetSubscriptionZeroPaymentLink(returnUrl string, subId string) string {
	if returnUrl == "" {
		return returnUrl
	}
	if strings.Contains(returnUrl, "?") {
		return fmt.Sprintf("%s&subId=%s&success=true", returnUrl, subId)
	} else {
		return fmt.Sprintf("%s?subId=%s&success=true", returnUrl, subId)
	}
}

func SubscriptionCancel(ctx context.Context, subscriptionId string, proration bool, invoiceNow bool, reason string) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	if sub.Status == consts.SubStatusCancelled || sub.Status == consts.SubStatusExpired {
		g.Log().Infof(ctx, "SubscriptionCancel, subscription already cancelled or expired")
		return nil
	}
	plan := query.GetPlanById(ctx, sub.PlanId)
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	if !config2.GetConfigInstance().IsServerDev() || !config2.GetConfigInstance().IsLocal() {
		// todo mark will support proration invoiceNow later
		invoiceNow = false
		proration = false
		// todo mark will support proration invoiceNow later
		// only local env can cancel immediately invoice_compute proration invoice
		utility.Assert(invoiceNow == false && proration == false, "cancel subscription with proration invoice immediate not support for this version")
	}
	var nextStatus = consts.SubStatusCancelled
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:         nextStatus,
		dao.Subscription.Columns().CancelReason:   reason,
		dao.Subscription.Columns().TrialEnd:       sub.CurrentPeriodStart - 1,
		dao.Subscription.Columns().GmtModify:      gtime.Now(),
		dao.Subscription.Columns().LastUpdateTime: gtime.Now().Timestamp(),
	}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	service3.TryCancelSubscriptionLatestInvoice(ctx, sub)

	{
		user := query.GetUserAccountById(ctx, sub.UserId)
		if user != nil {
			merchant := query.GetMerchantById(ctx, sub.MerchantId)
			if merchant != nil {
				var template = email.TemplateSubscriptionImmediateCancel
				if (sub.Status == consts.SubStatusIncomplete || sub.Status == consts.SubStatusActive) && sub.TrialEnd >= sub.CurrentPeriodEnd {
					//first trial period without payment
					template = email.TemplateSubscriptionCancelledByTrialEnd
				}
				if strings.Compare(reason, "CancelledByAnotherCreation") != 0 {
					err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, user.Language, template, "", &email.TemplateVariable{
						UserName:              user.FirstName + " " + user.LastName,
						MerchantProductName:   plan.PlanName,
						MerchantCustomerEmail: merchant.Email,
						MerchantName:          query.GetMerchantCountryConfigName(ctx, merchant.Id, user.CountryCode),
						PeriodEnd:             gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
					})
					if err != nil {
						g.Log().Errorf(ctx, "SendTemplateEmail SubscriptionCancel:%s", err.Error())
					}
				}
			}
		}

	}

	_, _ = redismq.Send(&redismq.Message{
		Topic:      redismq2.TopicSubscriptionCancel.Topic,
		Tag:        redismq2.TopicSubscriptionCancel.Tag,
		Body:       sub.SubscriptionId,
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     sub.MerchantId,
		Target:         fmt.Sprintf("Subscription(%v)", sub.SubscriptionId),
		Content:        fmt.Sprintf("Cancel(%s)", reason),
		UserId:         sub.UserId,
		SubscriptionId: sub.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return nil
}

func SubscriptionCancelAtPeriodEnd(ctx context.Context, subscriptionId string, proration bool, merchantMemberId int64) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	if sub.CancelAtPeriodEnd == 1 {
		return nil
	}

	plan := query.GetPlanById(ctx, sub.PlanId)
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().CancelAtPeriodEnd: 1,
		dao.Subscription.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}

	user := query.GetUserAccountById(ctx, sub.UserId)
	merchant := query.GetMerchantById(ctx, sub.MerchantId)
	// SendEmail
	if merchantMemberId > 0 {
		//merchant Cancel
		err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, user.Language, email.TemplateSubscriptionCancelledAtPeriodEndByMerchantAdmin, "", &email.TemplateVariable{
			UserName:              user.FirstName + " " + user.LastName,
			MerchantProductName:   plan.PlanName,
			MerchantCustomerEmail: merchant.Email,
			MerchantName:          query.GetMerchantCountryConfigName(ctx, merchant.Id, user.CountryCode),
			PeriodEnd:             gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
		})
		if err != nil {
			g.Log().Errorf(ctx, "SendTemplateEmail SubscriptionCancelAtPeriodEnd:%s", err.Error())
		}
	} else {
		err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, user.Language, email.TemplateSubscriptionCancelledAtPeriodEndByUser, "", &email.TemplateVariable{
			UserName:              user.FirstName + " " + user.LastName,
			MerchantProductName:   plan.PlanName,
			MerchantCustomerEmail: merchant.Email,
			MerchantName:          query.GetMerchantCountryConfigName(ctx, merchant.Id, user.CountryCode),
			PeriodEnd:             gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
		})
		if err != nil {
			g.Log().Errorf(ctx, "SendTemplateEmail SubscriptionCancelAtPeriodEnd:%s", err.Error())
		}
	}
	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicUserMetricUpdate.Topic,
		Tag:   redismq2.TopicUserMetricUpdate.Tag,
		Body: utility.MarshalToJsonString(&metric2.UserMetricUpdateMessage{
			UserId:         sub.UserId,
			SubscriptionId: sub.SubscriptionId,
			Description:    "SubscriptionCancelAtPeriodEnd",
		}),
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     sub.MerchantId,
		Target:         fmt.Sprintf("Subscription(%v)", sub.SubscriptionId),
		Content:        "EnableCancelAtPeriodEnd",
		UserId:         sub.UserId,
		SubscriptionId: sub.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return nil
}

func SubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, subscriptionId string, proration bool) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	if sub.CancelAtPeriodEnd == 0 {
		return nil
	}

	plan := query.GetPlanById(ctx, sub.PlanId)
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")

	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().CancelAtPeriodEnd: 0,
		dao.Subscription.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	user := query.GetUserAccountById(ctx, sub.UserId)
	merchant := query.GetMerchantById(ctx, sub.MerchantId)
	err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, user.Language, email.TemplateSubscriptionCancelLastCancelledAtPeriodEnd, "", &email.TemplateVariable{
		UserName:              user.FirstName + " " + user.LastName,
		MerchantProductName:   plan.PlanName,
		MerchantCustomerEmail: merchant.Email,
		MerchantName:          query.GetMerchantCountryConfigName(ctx, merchant.Id, user.CountryCode),
		PeriodEnd:             gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
	})
	if err != nil {
		g.Log().Errorf(ctx, "SendTemplateEmail SubscriptionCancelLastCancelAtPeriodEnd:%s", err.Error())
	}
	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicUserMetricUpdate.Topic,
		Tag:   redismq2.TopicUserMetricUpdate.Tag,
		Body: utility.MarshalToJsonString(&metric2.UserMetricUpdateMessage{
			UserId:         sub.UserId,
			SubscriptionId: sub.SubscriptionId,
			Description:    "SubscriptionCancelLastCancelAtPeriodEnd",
		}),
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     sub.MerchantId,
		Target:         fmt.Sprintf("Subscription(%v)", sub.SubscriptionId),
		Content:        "DisableCancelAtPeriodEnd",
		UserId:         sub.UserId,
		SubscriptionId: sub.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return nil
}

func SubscriptionAddNewTrialEnd(ctx context.Context, subscriptionId string, AppendNewTrialEndByHour int64) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	//utility.Assert(sub.Status != consts.SubStatusExpired && sub.Status != consts.SubStatusCancelled, "sub cancelled or sub expired")
	//utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	plan := query.GetPlanById(ctx, sub.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")

	utility.Assert(AppendNewTrialEndByHour > 0, "invalid AppendNewTrialEndByHour , should > 0")
	newTrialEnd := sub.CurrentPeriodEnd + AppendNewTrialEndByHour*3600

	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, utility.MaxInt64(newTrialEnd, sub.CurrentPeriodEnd), sub.PlanId)
	newStatus := sub.Status
	if newTrialEnd > gtime.Now().Timestamp() {
		//automatic change sub status to active
		newStatus = consts.SubStatusActive
		if sub.Status != consts.SubStatusActive {
			service3.TryCancelSubscriptionLatestInvoice(ctx, sub)
		}
	}
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:             newStatus,
		dao.Subscription.Columns().TrialEnd:           newTrialEnd,
		dao.Subscription.Columns().BillingCycleAnchor: newTrialEnd,
		dao.Subscription.Columns().DunningTime:        dunningTime,
		dao.Subscription.Columns().GmtModify:          gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	if sub.Status != consts.SubStatusActive && newStatus == consts.SubStatusActive {
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicSubscriptionActive.Topic,
			Tag:        redismq2.TopicSubscriptionActive.Tag,
			Body:       sub.SubscriptionId,
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
	}
	_, _ = redismq.Send(&redismq.Message{
		Topic:      redismq2.TopicSubscriptionUpdate.Topic,
		Tag:        redismq2.TopicSubscriptionUpdate.Tag,
		Body:       sub.SubscriptionId,
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicUserMetricUpdate.Topic,
		Tag:   redismq2.TopicUserMetricUpdate.Tag,
		Body: utility.MarshalToJsonString(&metric2.UserMetricUpdateMessage{
			UserId:         sub.UserId,
			SubscriptionId: sub.SubscriptionId,
			Description:    "SubscriptionAddNewTrialEnd",
		}),
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     sub.MerchantId,
		Target:         fmt.Sprintf("Subscription(%v)", sub.SubscriptionId),
		Content:        fmt.Sprintf("AddNewTrialEnd(%d)", newTrialEnd),
		UserId:         sub.UserId,
		SubscriptionId: sub.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return nil
}

func SubscriptionActiveTemporarily(ctx context.Context, subscriptionId string, expireTime int64) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusPending || sub.Status == consts.SubStatusProcessing, "subscription not in pending or processing status")
	utility.Assert(sub.CurrentPeriodStart < expireTime, "expireTime should greater than subscription's period start time")
	utility.Assert(sub.CurrentPeriodEnd >= expireTime, "expireTime should lower than subscription's period end time")
	utility.Assert(len(sub.LatestInvoiceId) > 0, "sub latest invoice not found")
	invoice := query.GetInvoiceByInvoiceId(ctx, sub.LatestInvoiceId)
	utility.Assert(invoice != nil, "sub latest invoice not found")
	utility.Assert(invoice.Status == consts.InvoiceStatusProcessing, "sub latest invoice not in processing")
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().CurrentPeriodPaid: expireTime,
		dao.Subscription.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	utility.AssertError(err, "Subscription Active Temporarily")

	_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().DayUtilDue: ((expireTime - invoice.FinishTime) / 86400) + 3,
		dao.Invoice.Columns().GmtModify:  gtime.Now(),
	}).Where(dao.Invoice.Columns().InvoiceId, invoice.InvoiceId).OmitNil().Update()
	utility.AssertError(err, "Subscription Active Temporarily")

	err = handler.MakeSubscriptionIncomplete(ctx, subscriptionId)
	if err != nil {
		return err
	}
	if sub.TrialEnd > 0 && sub.TrialEnd > sub.CurrentPeriodStart {
		// trial start
		oneUser := query.GetUserAccountById(ctx, sub.UserId)
		plan := query.GetPlanById(ctx, sub.PlanId)
		merchant := query.GetMerchantById(ctx, sub.MerchantId)
		if oneUser != nil && plan != nil && merchant != nil {
			err := email.SendTemplateEmail(ctx, sub.MerchantId, oneUser.Email, oneUser.TimeZone, oneUser.Language, email.TemplateSubscriptionTrialStart, "", &email.TemplateVariable{
				UserName:              oneUser.FirstName + " " + oneUser.LastName,
				MerchantProductName:   plan.PlanName,
				MerchantCustomerEmail: merchant.Email,
				MerchantName:          query.GetMerchantCountryConfigName(ctx, sub.MerchantId, oneUser.CountryCode),
			})
			if err != nil {
				g.Log().Errorf(ctx, "SendTemplateEmail TemplateSubscriptionTrialStart:%s", err.Error())
			}
		}
	}
	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicUserMetricUpdate.Topic,
		Tag:   redismq2.TopicUserMetricUpdate.Tag,
		Body: utility.MarshalToJsonString(&metric2.UserMetricUpdateMessage{
			UserId:         sub.UserId,
			SubscriptionId: sub.SubscriptionId,
			Description:    "SubscriptionActivateTemporarily",
		}),
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     sub.MerchantId,
		Target:         fmt.Sprintf("Subscription(%v)", sub.SubscriptionId),
		Content:        "ActiveTemporarily",
		UserId:         sub.UserId,
		SubscriptionId: sub.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)

	return nil
}

func SubscriptionEndTrial(ctx context.Context, subscriptionId string) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	plan := query.GetPlanById(ctx, sub.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	utility.Assert(sub.TrialEnd > gtime.Now().Timestamp(), "subscription not trialed")
	err := EndTrialManual(ctx, sub.SubscriptionId)
	if err != nil {
		return err
	}

	return nil
}

func EndTrialManual(ctx context.Context, subscriptionId string) (err error) {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.TrialEnd > gtime.Now().Timestamp(), "subscription not in trial period")
	newTrialEnd := sub.CurrentPeriodStart - 1
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, utility.MaxInt64(newTrialEnd, sub.CurrentPeriodEnd), sub.PlanId)
	newStatus := sub.Status
	if gtime.Now().Timestamp() > sub.CurrentPeriodEnd {
		// todo mark has unfinished pending update
		newStatus = consts.SubStatusIncomplete
		// Payment Pending Enter Incomplete
		plan := query.GetPlanById(ctx, sub.PlanId)

		var nextPeriodStart = gtime.Now().Timestamp()
		var nextPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, nextPeriodStart, nextPeriodStart, plan.Id)
		invoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
			UserId:             sub.UserId,
			Currency:           sub.Currency,
			PlanId:             sub.PlanId,
			Quantity:           sub.Quantity,
			AddonJsonData:      sub.AddonData,
			TaxPercentage:      sub.TaxPercentage,
			PeriodStart:        nextPeriodStart,
			PeriodEnd:          nextPeriodEnd,
			InvoiceName:        "SubscriptionCycle",
			FinishTime:         nextPeriodStart,
			BillingCycleAnchor: nextPeriodStart,
			VatNumber:          sub.VatNumber,
			ApplyPromoCredit:   false,
		})
		gatewayId, paymentMethodId := sub_update.VerifyPaymentGatewayMethod(ctx, sub.UserId, nil, "", sub.SubscriptionId)
		utility.Assert(gatewayId > 0, "gateway need specified")
		one, err := service3.CreateProcessingInvoiceForSub(ctx, sub.PlanId, invoice, sub, gatewayId, paymentMethodId, true, gtime.Now().Timestamp())
		if err != nil {
			g.Log().Print(ctx, "EndTrialManual CreateProcessingInvoiceForSub err:", err.Error())
			return err
		}
		createRes, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, one, false, "", "", "SubscriptionEndTrialManual", 0)
		if err != nil {
			g.Log().Print(ctx, "EndTrialManual CreateSubInvoicePaymentDefaultAutomatic err:", err.Error())
			return err
		}
		_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().CurrentPeriodStart: invoice.PeriodStart,
			dao.Subscription.Columns().CurrentPeriodEnd:   invoice.PeriodEnd,
			dao.Subscription.Columns().DunningTime:        dunningTime,
			dao.Subscription.Columns().GmtModify:          gtime.Now(),
		}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
		if err != nil {
			return err
		}
		g.Log().Print(ctx, "EndTrialManual CreateSubInvoicePaymentDefaultAutomatic:", utility.MarshalToJsonString(createRes))
		err = handler.HandleSubscriptionIncomplete(ctx, sub.SubscriptionId, gtime.Now().Timestamp())

	} else {
		_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().Status:         newStatus,
			dao.Subscription.Columns().TrialEnd:       newTrialEnd,
			dao.Subscription.Columns().DunningTime:    dunningTime,
			dao.Subscription.Columns().GmtModify:      gtime.Now(),
			dao.Subscription.Columns().LastUpdateTime: gtime.Now().Timestamp(),
		}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	}
	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicUserMetricUpdate.Topic,
		Tag:   redismq2.TopicUserMetricUpdate.Tag,
		Body: utility.MarshalToJsonString(&metric2.UserMetricUpdateMessage{
			UserId:         sub.UserId,
			SubscriptionId: sub.SubscriptionId,
			Description:    "SubscriptionEndTrialManual",
		}),
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     sub.MerchantId,
		Target:         fmt.Sprintf("Subscription(%s)", sub.SubscriptionId),
		Content:        "EndTrial",
		UserId:         sub.UserId,
		SubscriptionId: sub.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	if err != nil {
		g.Log().Print(ctx, "EndTrialManual err:", err.Error())
		return err
	}
	return nil
}

func MarkSubscriptionProcessed(ctx context.Context, subscriptionId string) error {
	utility.Assert(len(subscriptionId) > 0, "invalid subscriptionId")
	one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(one != nil, "subscription not found")
	utility.Assert(one.Status == consts.SubStatusPending, "sub not pending status")
	gateway := query.GetGatewayById(ctx, one.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	utility.Assert(gateway.GatewayType == consts.GatewayTypeWireTransfer, "not wire transfer type of subscription")
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:         consts.SubStatusProcessing,
		dao.Subscription.Columns().GmtModify:      gtime.Now(),
		dao.Subscription.Columns().LastUpdateTime: gtime.Now().Timestamp(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Subscription(%s)", one.SubscriptionId),
		Content:        "MarkSubscriptionProcessed",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	if err != nil {
		return err
	}
	return nil
}
