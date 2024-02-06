package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	"unibee-api/internal/logic/email"
	"unibee-api/internal/logic/gateway/ro"
	"unibee-api/internal/logic/invoice/handler"
	"unibee-api/internal/logic/invoice/invoice_compute"
	subscription2 "unibee-api/internal/logic/subscription"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func HandleSubscriptionWebhookEvent(ctx context.Context, subscription *entity.Subscription, eventType string, details *ro.GatewayDetailSubscriptionInternalResp) error {
	//更新 Subscription
	return UpdateSubWithGatewayDetailBack(ctx, subscription, details)
}

func UpdateSubWithGatewayDetailBack(ctx context.Context, subscription *entity.Subscription, details *ro.GatewayDetailSubscriptionInternalResp) error {
	if subscription.Type != consts.SubTypeDefault {
		// not sync attribute from gateway
		return nil
	}
	var cancelAtPeriodEnd = 0
	if details.CancelAtPeriodEnd {
		cancelAtPeriodEnd = 1
	}
	var firstPaidAt int64 = 0
	if subscription.FirstPaidAt == 0 && details.Status == consts.SubStatusActive {
		firstPaidAt = gtime.Now().Timestamp()
	} else {
		firstPaidAt = subscription.FirstPaidAt
	}
	var gmtModify = subscription.GmtModify
	if subscription.Status != int(details.Status) ||
		subscription.CancelAtPeriodEnd != cancelAtPeriodEnd ||
		subscription.BillingCycleAnchor != details.BillingCycleAnchor ||
		subscription.TrialEnd != details.TrialEnd ||
		subscription.FirstPaidAt != firstPaidAt {
		gmtModify = gtime.Now()
	}

	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, utility.MaxInt64(details.TrialEnd, subscription.CurrentPeriodEnd), uint64(subscription.PlanId))
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:                      details.Status,
		dao.Subscription.Columns().GatewaySubscriptionId:       details.GatewaySubscriptionId,
		dao.Subscription.Columns().GatewayStatus:               details.GatewayStatus,
		dao.Subscription.Columns().GatewayItemData:             details.GatewayItemData,
		dao.Subscription.Columns().CancelAtPeriodEnd:           cancelAtPeriodEnd,
		dao.Subscription.Columns().BillingCycleAnchor:          details.BillingCycleAnchor,
		dao.Subscription.Columns().GatewayLatestInvoiceId:      details.GatewayLatestInvoiceId,
		dao.Subscription.Columns().GatewayDefaultPaymentMethod: details.GatewayDefaultPaymentMethod,
		dao.Subscription.Columns().TrialEnd:                    details.TrialEnd,
		dao.Subscription.Columns().DunningTime:                 dunningTime,
		dao.Subscription.Columns().GmtModify:                   gmtModify,
		dao.Subscription.Columns().FirstPaidAt:                 firstPaidAt,
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
		dao.Subscription.Columns().FirstPaidAt:            payment.PaidAt,
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

func UpdateSubscriptionBillingCycleWithPayment(ctx context.Context, payment *entity.Payment) error {
	utility.Assert(payment != nil, "UpdateSubscriptionBillingCycleWithPayment payment is nil")
	utility.Assert(len(payment.SubscriptionId) > 0, "UpdateSubscriptionBillingCycleWithPayment payment subId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, payment.SubscriptionId)
	utility.Assert(sub != nil, "UpdateSubscriptionBillingCycleWithPayment sub not found")
	invoice := query.GetInvoiceByInvoiceId(ctx, payment.InvoiceId)
	utility.Assert(invoice != nil, "UpdateSubscriptionBillingCycleWithPayment invoice not found payment:"+payment.PaymentId)
	var firstPaidAt int64 = 0
	if sub.FirstPaidAt == 0 && payment.Status == consts.PAY_SUCCESS {
		firstPaidAt = payment.PaidAt
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
		dao.Subscription.Columns().FirstPaidAt:            firstPaidAt,
	}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}

type SubscriptionPaymentSuccessWebHookReq struct {
	Payment                     *entity.Payment                           `json:"payment" `
	GatewaySubscriptionDetail   *ro.GatewayDetailSubscriptionInternalResp `json:"gatewaySubscriptionDetail"`
	GatewayInvoiceDetail        *ro.GatewayDetailInvoiceInternalResp      `json:"gatewayInvoiceDetail"`
	GatewayPaymentId            string                                    `json:"gatewayPaymentId" `
	GatewayInvoiceId            string                                    `json:"gatewayInvoiceId"`
	GatewaySubscriptionId       string                                    `json:"gatewaySubscriptionId" `
	GatewaySubscriptionUpdateId string                                    `json:"gatewaySubscriptionUpdateId"`
	Status                      consts.SubscriptionStatusEnum             `json:"status"`
	GatewayStatus               string                                    `json:"gatewayStatus"                  `
	Data                        string                                    `json:"data"`
	GatewayItemData             string                                    `json:"gatewayItemData"`
	CancelAtPeriodEnd           bool                                      `json:"cancelAtPeriodEnd"`
	CurrentPeriodEnd            int64                                     `json:"currentPeriodEnd"`
	CurrentPeriodStart          int64                                     `json:"currentPeriodStart"`
	TrialEnd                    int64                                     `json:"trialEnd"`
}

func checkAndListAddonsFromParams(ctx context.Context, addonParams []*ro.SubscriptionPlanAddonParamRo, gatewayId int64) []*ro.SubscriptionPlanAddonRo {
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
				gatewayPlan := query.GetGatewayPlan(ctx, int64(mapPlans[param.AddonPlanId].Id), gatewayId)
				utility.Assert(len(gatewayPlan.GatewayPlanId) > 0, fmt.Sprintf("internal error PlanId:%v GatewayId:%v GatewayPlanId invalid", param.AddonPlanId, gatewayId))
				utility.Assert(gatewayPlan.Status == consts.GatewayPlanStatusActive, fmt.Sprintf("internal error PlanId:%v GatewayId:%v gatewayPlanStatus not active", param.AddonPlanId, gatewayId))
				addons = append(addons, &ro.SubscriptionPlanAddonRo{
					Quantity:         param.Quantity,
					AddonPlan:        mapPlans[param.AddonPlanId],
					AddonGatewayPlan: gatewayPlan,
				})
			}
		}
	}
	return addons
}

func HandleSubscriptionPaymentUpdate(ctx context.Context, req *SubscriptionPaymentSuccessWebHookReq) error {
	sub := query.GetSubscriptionByGatewaySubscriptionId(ctx, req.GatewaySubscriptionId)
	if sub == nil {
		return gerror.Newf("HandleSubscriptionPaymentUpdate sub not found %s", req.GatewaySubscriptionId)
	}
	if sub.Type != consts.SubTypeDefault {
		return gerror.Newf("HandleSubscriptionPaymentUpdate not gateway subscription %s", sub.SubscriptionId)
	}
	eiPendingSubUpdate := query.GetUnfinishedEffectImmediateSubscriptionPendingUpdateByGatewayUpdateId(ctx, req.Payment.PaymentId)
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
					sub = query.GetSubscriptionByGatewaySubscriptionId(ctx, req.GatewaySubscriptionId)

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
				sub = query.GetSubscriptionByGatewaySubscriptionId(ctx, req.GatewaySubscriptionId)
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

	err := email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, template, "", &email.TemplateVariable{
		UserName:            user.FirstName + " " + user.LastName,
		MerchantProductName: merchantProductName,
		MerchantCustomEmail: merchant.Email,
		MerchantName:        merchant.Name,
		DateNow:             gtime.Now(),
		PeriodEnd:           gtime.NewFromTimeStamp(one.CurrentPeriodEnd),
	})
	if err != nil {
		return err
	}

	return nil
}
