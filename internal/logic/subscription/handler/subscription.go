package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	v1 "go-oversea-pay/api/open/payment"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/channel/ro"
	"go-oversea-pay/internal/logic/email"
	"go-oversea-pay/internal/logic/invoice/handler"
	"go-oversea-pay/internal/logic/invoice/invoice_compute"
	"go-oversea-pay/internal/logic/payment/service"
	subscription2 "go-oversea-pay/internal/logic/subscription"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strconv"
)

func CreateSubInvoicePayment(ctx context.Context, sub *entity.Subscription, invoice *ro.InvoiceDetailSimplify, billingReason string) (channelInternalPayResult *ro.CreatePayInternalResp, err error) {
	user := query.GetUserAccountById(ctx, uint64(sub.UserId))
	var mobile = ""
	var firstName = ""
	var lastName = ""
	var gender = ""
	var email = ""
	if user != nil {
		mobile = user.Mobile
		firstName = user.FirstName
		lastName = user.LastName
		gender = user.Gender
		email = user.Email
	}
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, sub.ChannelId)
	if payChannel == nil {
		return nil, gerror.New("SubscriptionBillingCycleDunningInvoice pay channel not found")
	}
	merchantInfo := query.GetMerchantInfoById(ctx, sub.MerchantId)
	if merchantInfo == nil {
		return nil, gerror.New("SubscriptionBillingCycleDunningInvoice merchantInfo not found")
	}
	return service.DoChannelPay(ctx, &ro.CreatePayContext{
		PayChannel: payChannel,
		Pay: &entity.Payment{
			SubscriptionId:  sub.SubscriptionId,
			BizId:           sub.SubscriptionId,
			BizType:         consts.BIZ_TYPE_SUBSCRIPTION,
			AuthorizeStatus: consts.AUTHORIZED,
			UserId:          sub.UserId,
			ChannelId:       int64(payChannel.Id),
			TotalAmount:     invoice.TotalAmount,
			Currency:        invoice.Currency,
			CountryCode:     sub.CountryCode,
			MerchantId:      sub.MerchantId,
			CompanyId:       merchantInfo.CompanyId,
			BillingReason:   billingReason,
		},
		Platform:      "WEB",
		DeviceType:    "Web",
		ShopperUserId: strconv.FormatInt(sub.UserId, 10),
		ShopperEmail:  email,
		ShopperLocale: "en",
		Mobile:        mobile,
		Invoice:       invoice,
		ShopperName: &v1.OutShopperName{
			FirstName: firstName,
			LastName:  lastName,
			Gender:    gender,
		},
		MediaData:              map[string]string{"BillingReason": billingReason},
		MerchantOrderReference: sub.SubscriptionId,
		PayMethod:              1, //automatic
		DaysUtilDue:            5, // todo mark
		ChannelPaymentMethod:   sub.ChannelDefaultPaymentMethod,
	})
}

func HandleSubscriptionWebhookEvent(ctx context.Context, subscription *entity.Subscription, eventType string, details *ro.ChannelDetailSubscriptionInternalResp) error {
	//更新 Subscription
	return UpdateSubWithChannelDetailBack(ctx, subscription, details)
}

func UpdateSubWithChannelDetailBack(ctx context.Context, subscription *entity.Subscription, details *ro.ChannelDetailSubscriptionInternalResp) error {
	if subscription.Type != consts.SubTypeDefault {
		// not sync attribute from channel
		return nil
	}
	var cancelAtPeriodEnd = 0
	if details.CancelAtPeriodEnd {
		cancelAtPeriodEnd = 1
	}
	var firstPayTime *gtime.Time
	if subscription.FirstPayTime == nil && details.Status == consts.SubStatusActive {
		firstPayTime = gtime.Now()
	}
	var gmtModify = subscription.GmtModify
	if subscription.Status != int(details.Status) ||
		subscription.CancelAtPeriodEnd != cancelAtPeriodEnd ||
		subscription.BillingCycleAnchor != details.BillingCycleAnchor ||
		subscription.TrialEnd != details.TrialEnd ||
		subscription.FirstPayTime != firstPayTime {
		gmtModify = gtime.Now()
	}

	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, utility.MaxInt64(details.TrialEnd, subscription.CurrentPeriodEnd), uint64(subscription.PlanId))
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:                      details.Status,
		dao.Subscription.Columns().ChannelSubscriptionId:       details.ChannelSubscriptionId,
		dao.Subscription.Columns().ChannelStatus:               details.ChannelStatus,
		dao.Subscription.Columns().ChannelItemData:             details.ChannelItemData,
		dao.Subscription.Columns().CancelAtPeriodEnd:           cancelAtPeriodEnd,
		dao.Subscription.Columns().BillingCycleAnchor:          details.BillingCycleAnchor,
		dao.Subscription.Columns().ChannelLatestInvoiceId:      details.ChannelLatestInvoiceId,
		dao.Subscription.Columns().ChannelDefaultPaymentMethod: details.ChannelDefaultPaymentMethod,
		dao.Subscription.Columns().TrialEnd:                    details.TrialEnd,
		dao.Subscription.Columns().DunningTime:                 dunningTime,
		dao.Subscription.Columns().GmtModify:                   gmtModify,
		dao.Subscription.Columns().FirstPayTime:                firstPayTime,
	}).Where(dao.Subscription.Columns().Id, subscription.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}

func HandleSubscriptionCreatePaymentSuccess(ctx context.Context, sub *entity.Subscription, payment *entity.Payment) error {
	utility.Assert(payment != nil, "HandleSubscriptionCreatePaymentSuccess payment is nil")
	utility.Assert(len(payment.SubscriptionId) > 0, "HandleSubscriptionCreatePaymentSuccess payment subId is nil")
	sub = query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
	utility.Assert(sub != nil, "HandleSubscriptionCreatePaymentSuccess sub not found")
	invoice := query.GetInvoiceByInvoiceId(ctx, payment.InvoiceId)
	utility.Assert(invoice != nil, "HandleSubscriptionCreatePaymentSuccess invoice not found payment:"+payment.PaymentId)
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, invoice.PeriodEnd, uint64(sub.PlanId))
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:                 consts.SubStatusActive,
		dao.Subscription.Columns().CurrentPeriodStart:     invoice.PeriodStart,
		dao.Subscription.Columns().CurrentPeriodEnd:       invoice.PeriodEnd,
		dao.Subscription.Columns().CurrentPeriodStartTime: gtime.NewFromTimeStamp(invoice.PeriodStart),
		dao.Subscription.Columns().CurrentPeriodEndTime:   gtime.NewFromTimeStamp(invoice.PeriodEnd),
		dao.Subscription.Columns().DunningTime:            dunningTime,
		dao.Subscription.Columns().GmtModify:              gtime.Now(),
		dao.Subscription.Columns().FirstPayTime:           payment.PaidTime,
	}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}

func FinishNextBillingCycleForSubscription(ctx context.Context, sub *entity.Subscription, payment *entity.Payment) error {
	// billing-cycle
	err := CreateOrUpdateSubscriptionTimeline(ctx, sub, fmt.Sprintf("cycle-paymentId-%s", payment.PaymentId))
	if err != nil {
		g.Log().Errorf(ctx, "FinishNextBillingCycleForSubscription error:%s", err.Error())
	}
	err = UpdateSubscriptionBillingCycleWithPayment(ctx, payment)
	return err
}

func ChangeTrialEnd(ctx context.Context, newTrialEnd int64, subscriptionId string) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status != consts.SubStatusExpired && sub.Status != consts.SubStatusCancelled, "sub cancelled or sub expired")

	var newBillingCycleAnchor = utility.MaxInt64(newTrialEnd, sub.CurrentPeriodEnd)
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, newBillingCycleAnchor, uint64(sub.PlanId))
	newStatus := sub.Status
	if newTrialEnd > gtime.Now().Timestamp() {
		//automatic change sub status to active
		newStatus = consts.SubStatusActive
	}
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:             newStatus,
		dao.Subscription.Columns().TrialEnd:           newTrialEnd,
		dao.Subscription.Columns().DunningTime:        dunningTime,
		dao.Subscription.Columns().BillingCycleAnchor: newBillingCycleAnchor,
		dao.Subscription.Columns().GmtModify:          gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}

func SubscriptionIncomplete(ctx context.Context, subscriptionId string) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd) < gtime.Now().Timestamp(), "subscription not incomplete")
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:    consts.SubStatusIncomplete,
		dao.Subscription.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}

func EndTrialManual(ctx context.Context, subscriptionId string) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.TrialEnd > gtime.Now().Timestamp(), "subscription not in trial period")
	newTrialEnd := sub.CurrentPeriodStart - 1
	var newBillingCycleAnchor = utility.MaxInt64(newTrialEnd, sub.CurrentPeriodEnd)
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, newBillingCycleAnchor, uint64(sub.PlanId))
	newStatus := sub.Status
	if gtime.Now().Timestamp() > sub.CurrentPeriodEnd {
		// todo mark has unfinished pending update
		newStatus = consts.SubStatusIncomplete
		// Payment Pending Enter Incomplete
		plan := query.GetPlanById(ctx, sub.PlanId)

		var nextPeriodStart = gtime.Now().Timestamp()
		var nextPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, nextPeriodStart, plan.Id)
		invoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
			Currency:      sub.Currency,
			PlanId:        sub.PlanId,
			Quantity:      sub.Quantity,
			AddonJsonData: sub.AddonData,
			TaxScale:      sub.TaxScale,
			PeriodStart:   nextPeriodStart,
			PeriodEnd:     nextPeriodEnd,
		})
		createRes, err := CreateSubInvoicePayment(ctx, sub, invoice, "SubscriptionCycle")
		if err != nil {
			g.Log().Print(ctx, "EndTrialManual CreateSubInvoicePayment err:", err.Error())
			return err
		}
		_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().CurrentPeriodStart: invoice.PeriodStart,
			dao.Subscription.Columns().CurrentPeriodEnd:   invoice.PeriodEnd,
			dao.Subscription.Columns().DunningTime:        dunningTime,
			dao.Subscription.Columns().BillingCycleAnchor: newBillingCycleAnchor,
			dao.Subscription.Columns().GmtModify:          gtime.Now(),
		}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
		if err != nil {
			return err
		}
		g.Log().Print(ctx, "EndTrialManual CreateSubInvoicePayment:", utility.MarshalToJsonString(createRes))
		err = SubscriptionIncomplete(ctx, sub.SubscriptionId)
		if err != nil {
			g.Log().Print(ctx, "EndTrialManual SubscriptionIncomplete err:", err.Error())
			return err
		}
	} else {
		_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().Status:             newStatus,
			dao.Subscription.Columns().TrialEnd:           newTrialEnd,
			dao.Subscription.Columns().DunningTime:        dunningTime,
			dao.Subscription.Columns().BillingCycleAnchor: newBillingCycleAnchor,
			dao.Subscription.Columns().GmtModify:          gtime.Now(),
		}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateSubscriptionBillingCycleWithPayment(ctx context.Context, payment *entity.Payment) error {
	utility.Assert(payment != nil, "UpdateSubscriptionBillingCycleWithPayment payment is nil")
	utility.Assert(len(payment.SubscriptionId) > 0, "UpdateSubscriptionBillingCycleWithPayment payment subId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
	utility.Assert(sub != nil, "UpdateSubscriptionBillingCycleWithPayment sub not found")
	invoice := query.GetInvoiceByInvoiceId(ctx, payment.InvoiceId)
	utility.Assert(invoice != nil, "UpdateSubscriptionBillingCycleWithPayment invoice not found payment:"+payment.PaymentId)
	var firstPayTime = sub.FirstPayTime
	if sub.FirstPayTime == nil && payment.Status == consts.PAY_SUCCESS {
		firstPayTime = payment.PaidTime
	}
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, utility.MaxInt64(invoice.PeriodEnd, sub.TrialEnd), uint64(sub.PlanId))
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:                 consts.SubStatusActive,
		dao.Subscription.Columns().CurrentPeriodStart:     invoice.PeriodStart,
		dao.Subscription.Columns().CurrentPeriodEnd:       invoice.PeriodEnd,
		dao.Subscription.Columns().CurrentPeriodStartTime: gtime.NewFromTimeStamp(invoice.PeriodStart),
		dao.Subscription.Columns().CurrentPeriodEndTime:   gtime.NewFromTimeStamp(invoice.PeriodEnd),
		dao.Subscription.Columns().DunningTime:            dunningTime,
		dao.Subscription.Columns().GmtModify:              gtime.Now(),
		dao.Subscription.Columns().FirstPayTime:           firstPayTime,
	}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}

type SubscriptionPaymentSuccessWebHookReq struct {
	Payment                     *entity.Payment                           `json:"payment" `
	ChannelSubscriptionDetail   *ro.ChannelDetailSubscriptionInternalResp `json:"channelSubscriptionDetail"`
	ChannelInvoiceDetail        *ro.ChannelDetailInvoiceInternalResp      `json:"channelInvoiceDetail"`
	ChannelPaymentId            string                                    `json:"channelPaymentId" `
	ChannelInvoiceId            string                                    `json:"channelInvoiceId"`
	ChannelSubscriptionId       string                                    `json:"channelSubscriptionId" `
	ChannelSubscriptionUpdateId string                                    `json:"channelSubscriptionUpdateId"`
	Status                      consts.SubscriptionStatusEnum             `json:"status"`
	ChannelStatus               string                                    `json:"channelStatus"                  `
	Data                        string                                    `json:"data"`
	ChannelItemData             string                                    `json:"channelItemData"`
	CancelAtPeriodEnd           bool                                      `json:"cancelAtPeriodEnd"`
	CurrentPeriodEnd            int64                                     `json:"currentPeriodEnd"`
	CurrentPeriodStart          int64                                     `json:"currentPeriodStart"`
	TrialEnd                    int64                                     `json:"trialEnd"`
}

func checkAndListAddonsFromParams(ctx context.Context, addonParams []*ro.SubscriptionPlanAddonParamRo, channelId int64) []*ro.SubscriptionPlanAddonRo {
	var addons []*ro.SubscriptionPlanAddonRo
	var totalAddonIds []int64
	if len(addonParams) > 0 {
		for _, s := range addonParams {
			totalAddonIds = append(totalAddonIds, s.AddonPlanId) // 添加到整数列表中
		}
	}
	var allAddonList []*entity.SubscriptionPlan
	if len(totalAddonIds) > 0 {
		//查询所有 Plan
		err := dao.SubscriptionPlan.Ctx(ctx).WhereIn(dao.SubscriptionPlan.Columns().Id, totalAddonIds).OmitEmpty().Scan(&allAddonList)
		if err == nil {
			//整合进列表
			mapPlans := make(map[int64]*entity.SubscriptionPlan)
			for _, pair := range allAddonList {
				key := int64(pair.Id)
				value := pair
				mapPlans[key] = value
			}
			for _, param := range addonParams {
				//所有 Addon 项目必须要能查到
				//类型是 Addon
				//未删除
				//数量大于 0
				utility.Assert(mapPlans[param.AddonPlanId] != nil, fmt.Sprintf("AddonPlanId not found:%v", param.AddonPlanId))
				utility.Assert(mapPlans[param.AddonPlanId].Type == consts.PlanTypeAddon, fmt.Sprintf("Id:%v not Addon Type", param.AddonPlanId))
				utility.Assert(mapPlans[param.AddonPlanId].IsDeleted == 0, fmt.Sprintf("Addon Id:%v is Deleted", param.AddonPlanId))
				utility.Assert(param.Quantity > 0, fmt.Sprintf("Id:%v quantity invalid", param.AddonPlanId))
				planChannel := query.GetPlanChannel(ctx, int64(mapPlans[param.AddonPlanId].Id), channelId)
				utility.Assert(len(planChannel.ChannelPlanId) > 0, fmt.Sprintf("internal error PlanId:%v ChannelId:%v channelPlanId invalid", param.AddonPlanId, channelId))
				utility.Assert(planChannel.Status == consts.PlanChannelStatusActive, fmt.Sprintf("internal error PlanId:%v ChannelId:%v channelPlanStatus not active", param.AddonPlanId, channelId))
				addons = append(addons, &ro.SubscriptionPlanAddonRo{
					Quantity:         param.Quantity,
					AddonPlan:        mapPlans[param.AddonPlanId],
					AddonPlanChannel: planChannel,
				})
			}
		}
	}
	return addons
}

func HandleSubscriptionPaymentUpdate(ctx context.Context, req *SubscriptionPaymentSuccessWebHookReq) error {
	sub := query.GetSubscriptionByChannelSubscriptionId(ctx, req.ChannelSubscriptionId)
	if sub == nil {
		return gerror.Newf("HandleSubscriptionPaymentUpdate sub not found %s", req.ChannelSubscriptionId)
	}
	if sub.Type != consts.SubTypeDefault {
		return gerror.Newf("HandleSubscriptionPaymentUpdate not channel subscription %s", sub.SubscriptionId)
	}
	eiPendingSubUpdate := query.GetUnfinishedEffectImmediateSubscriptionPendingUpdateByChannelUpdateId(ctx, req.Payment.PaymentId)
	if eiPendingSubUpdate != nil {
		// subscription_update
		//更新单支付成功, EffectImmediate=true 需要用户 3DS 验证等场景
		if req.Payment.Status == consts.PAY_SUCCESS {
			_, err := FinishPendingUpdateForSubscription(ctx, sub, eiPendingSubUpdate.UpdateSubscriptionId)
			if err != nil {
				return err
			}
		}
	} else {
		// subscription_cycle
		var byUpdate = false
		if len(sub.PendingUpdateId) > 0 {
			//有 pending 的更新单存在，检查支付是否对应更新单
			pendingSubUpdate := query.GetUnfinishedSubscriptionPendingUpdateByPendingUpdateId(ctx, sub.PendingUpdateId)
			if pendingSubUpdate.UpdateAmount == req.Payment.TotalAmount {
				// compensate next pending update billing cycle invoice
				one := query.GetInvoiceByPaymentId(ctx, req.Payment.PaymentId)
				if one == nil {
					sub = query.GetSubscriptionByChannelSubscriptionId(ctx, req.ChannelSubscriptionId)

					plan := query.GetPlanById(ctx, pendingSubUpdate.UpdatePlanId)
					var nextPeriodStart = sub.CurrentPeriodEnd
					if sub.TrialEnd > sub.CurrentPeriodEnd {
						nextPeriodStart = sub.TrialEnd
					}
					var nextPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, nextPeriodStart, plan.Id)

					invoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
						Currency:      pendingSubUpdate.UpdateCurrency,
						PlanId:        pendingSubUpdate.UpdatePlanId,
						Quantity:      pendingSubUpdate.UpdateQuantity,
						AddonJsonData: pendingSubUpdate.UpdateAddonData,
						TaxScale:      sub.TaxScale,
						PeriodStart:   nextPeriodStart,
						PeriodEnd:     nextPeriodEnd,
					})
					_, _ = handler.CreateOrUpdateInvoiceFromPayment(ctx, invoice, req.Payment)
				}

				if req.Payment.Status == consts.PAY_SUCCESS {
					_, err := FinishPendingUpdateForSubscription(ctx, sub, pendingSubUpdate.UpdateSubscriptionId)
					if err != nil {
						return err
					}

					err = FinishNextBillingCycleForSubscription(ctx, sub, req.Payment)
					if err != nil {
						return err
					}
				}
				byUpdate = true
			}
		}
		if !byUpdate && req.Payment.TotalAmount == sub.Amount {
			// compensate billing cycle invoice
			one := query.GetInvoiceByPaymentId(ctx, req.Payment.PaymentId)
			if one == nil {
				sub = query.GetSubscriptionByChannelSubscriptionId(ctx, req.ChannelSubscriptionId)
				plan := query.GetPlanById(ctx, sub.PlanId)

				var nextPeriodStart = sub.CurrentPeriodEnd
				if sub.TrialEnd > sub.CurrentPeriodEnd {
					nextPeriodStart = sub.TrialEnd
				}
				var nextPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, nextPeriodStart, plan.Id)

				invoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
					Currency:      sub.Currency,
					PlanId:        sub.PlanId,
					Quantity:      sub.Quantity,
					AddonJsonData: sub.AddonData,
					TaxScale:      sub.TaxScale,
					PeriodStart:   nextPeriodStart,
					PeriodEnd:     nextPeriodEnd,
				})
				_, _ = handler.CreateOrUpdateInvoiceFromPayment(ctx, invoice, req.Payment)
			}
			if req.Payment.Status == consts.PAY_SUCCESS {
				err := FinishNextBillingCycleForSubscription(ctx, sub, req.Payment)
				if err != nil {
					return err
				}
			}
		} else if !byUpdate {
			// other situation payment need implementation todo mark

		}
	}
	return nil
}

func SendSubscriptionEmailToUser(ctx context.Context, subscriptionId string, template string) error {
	one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(one != nil, "invoice not found")
	utility.Assert(one.UserId > 0, "invoice userId not found")
	utility.Assert(one.MerchantId > 0, "invoice merchantId not found")
	user := query.GetUserAccountById(ctx, uint64(one.UserId))
	merchant := query.GetMerchantInfoById(ctx, one.MerchantId)
	var merchantProductName = ""
	sub := query.GetSubscriptionBySubscriptionId(ctx, one.SubscriptionId)
	if sub != nil {
		plan := query.GetPlanById(ctx, sub.PlanId)
		merchantProductName = plan.PlanName
	}

	err := email.SendTemplateEmail(ctx, merchant.Id, user.Email, template, "", &email.TemplateVariable{
		UserName:            user.FirstName + " " + user.LastName,
		MerchantProductName: merchantProductName,
		MerchantCustomEmail: merchant.Email,
		MerchantName:        merchant.Name,
		DateNow:             gtime.Now().Layout(`2006-01-02`),
		PeriodEnd:           gtime.NewFromTimeStamp(one.CurrentPeriodEnd).Layout(`2006-01-02`),
	})
	if err != nil {
		return err
	}

	return nil
}
