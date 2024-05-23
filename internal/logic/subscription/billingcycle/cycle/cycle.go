package cycle

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	config2 "unibee/internal/cmd/config"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/discount"
	"unibee/internal/logic/invoice/invoice_compute"
	handler2 "unibee/internal/logic/invoice/service"
	"unibee/internal/logic/payment/service"
	subscription2 "unibee/internal/logic/subscription"
	"unibee/internal/logic/subscription/billingcycle/expire"
	"unibee/internal/logic/subscription/config"
	"unibee/internal/logic/subscription/handler"
	service2 "unibee/internal/logic/subscription/service"
	"unibee/internal/logic/user"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type BillingCycleWalkRes struct {
	WalkUnfinished bool
	Message        string
}

func SubPipeBillingCycleWalk(ctx context.Context, subId string, timeNow int64, source string) (*BillingCycleWalkRes, error) {
	if len(subId) == 0 {
		return &BillingCycleWalkRes{Message: "SubId Is Nil"}, nil
	}
	sub := query.GetSubscriptionBySubscriptionId(ctx, subId)
	if sub == nil {
		return &BillingCycleWalkRes{Message: "Sub Not Found"}, nil
	}
	key := fmt.Sprintf("SubscriptionCycleWalk-%s", sub.SubscriptionId)
	if utility.TryLock(ctx, key, 60) {
		g.Log().Print(ctx, source, "GetLock 60s", key)
		defer func() {
			utility.ReleaseLock(ctx, key)
			g.Log().Print(ctx, source, "ReleaseLock", key)
		}()
		_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().TaskTime: gtime.Now(),
		}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
		if err != nil {
			g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice Update TaskTime err:", err.Error())
		}
		// generate invoice and payment ahead
		latestInvoice := query.GetInvoiceByInvoiceId(ctx, sub.LatestInvoiceId)
		// todo mark
		//if latestInvoice.FinishTime > 0 && latestInvoice.FinishTime+(latestInvoice.DayUtilDue*86400) < timeNow {
		//	//invoice has expired
		//	service3.CancelInvoiceForSubscription(ctx, sub)
		//	return &BillingCycleWalkRes{WalkUnfinished: true, Message: "Subscription LatestInvoice Expired"}, nil
		//}

		trackForSubscription(ctx, sub, timeNow)

		var needInvoiceGenerate = true
		var needInvoiceFirstTryPayment = false
		if latestInvoice != nil && (latestInvoice.Status == consts.InvoiceStatusProcessing) {
			needInvoiceGenerate = false
			if timeNow > utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)-config.GetMerchantSubscriptionConfig(ctx, sub.MerchantId).TryAutomaticPaymentBeforePeriodEnd {
				needInvoiceFirstTryPayment = true
			}
		} else if latestInvoice != nil && latestInvoice.Status == consts.InvoiceStatusPaid && timeNow < latestInvoice.PeriodStart {
			needInvoiceGenerate = false
		} else if timeNow < sub.DunningTime {
			needInvoiceGenerate = false
		}
		if sub.Status == consts.SubStatusExpired || sub.Status == consts.SubStatusCancelled {
			return &BillingCycleWalkRes{WalkUnfinished: false, Message: "Nothing Todo As Sub Cancelled Or Expired"}, nil
		} else if sub.Status == consts.SubStatusPending || sub.Status == consts.SubStatusProcessing {
			if utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)+(2*24*60*60) < timeNow {
				// first time create sub expired
				err = expire.SubscriptionExpire(ctx, sub, "NotPayAfter48Hours")
				if err != nil {
					g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice SubscriptionExpire SubStatus:Created", err.Error())
					return nil, err
				} else {
					return &BillingCycleWalkRes{WalkUnfinished: true, Message: "SubscriptionExpire From Create Status As Payment Out Of 2 Days"}, nil
				}
			} else {
				return &BillingCycleWalkRes{WalkUnfinished: false, Message: "Nothing Todo As Sub At Create Status NotPayBefore48Hours"}, nil
			}
		} else if timeNow > utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd) && sub.CancelAtPeriodEnd == 1 && sub.Status != consts.SubStatusCancelled {
			// sub set cancelAtPeriodEnd, need cancel by system
			needInvoiceGenerate = false
			err = service2.SubscriptionCancel(ctx, sub.SubscriptionId, false, false, "CancelAtPeriodEndBySystem")
			if err != nil {
				g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice SubscriptionCancel err:", err.Error())
				return nil, err
			} else {
				return &BillingCycleWalkRes{WalkUnfinished: true, Message: "SubscriptionCancel At Billing Cycle End By CurrentPeriodEnd Set"}, nil
			}
		} else if !needInvoiceGenerate && !needInvoiceFirstTryPayment && isSubscriptionExpireExcludePending(ctx, sub, timeNow) {
			// invoice not generate and sub out of time, need expired by system
			err = expire.SubscriptionExpire(ctx, sub, "CycleExpireWithoutPay")
			if err != nil {
				g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice SubscriptionExpire", err.Error())
				return nil, err
			} else {
				return &BillingCycleWalkRes{WalkUnfinished: true, Message: "SubscriptionExpire From Billing Cycle As Payment Out Of Permission Days"}, nil
			}
		} else {
			if sub.CancelAtPeriodEnd == 1 {
				return &BillingCycleWalkRes{WalkUnfinished: false, Message: "Nothing Todo As CancelPeriodEnd Set"}, nil
			}
			// Unpaid after period end or trial end
			if utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd) < timeNow && sub.Status != consts.SubStatusIncomplete {
				err = handler.HandleSubscriptionIncomplete(ctx, sub.SubscriptionId, timeNow)
				if err != nil {
					g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice HandleSubscriptionIncomplete err:", err.Error())
					return nil, err
				} else {
					return &BillingCycleWalkRes{WalkUnfinished: true, Message: "HandleSubscriptionIncomplete As Not Paid After CurrentPeriodEnd Or TrialEnd"}, nil
				}
			}
			var taxPercentage = sub.TaxPercentage
			percentage, err := user.GetUserTaxPercentage(ctx, sub.UserId)
			if err == nil {
				taxPercentage = percentage
			}

			if needInvoiceGenerate {
				var invoice *bean.InvoiceSimplify
				var discountCode = ""
				canApply, isRecurring, _ := discount.UserDiscountApplyPreview(ctx, &discount.UserDiscountApplyReq{
					MerchantId:     sub.MerchantId,
					UserId:         sub.UserId,
					DiscountCode:   sub.DiscountCode,
					Currency:       sub.Currency,
					SubscriptionId: sub.SubscriptionId,
					PLanId:         sub.PlanId,
				})
				if canApply && isRecurring {
					discountCode = sub.DiscountCode
				}
				pendingUpdate := query.GetUnfinishedSubscriptionPendingUpdateByPendingUpdateId(ctx, sub.PendingUpdateId)
				if pendingUpdate != nil {
					//generate PendingUpdate cycle invoice
					plan := query.GetPlanById(ctx, pendingUpdate.UpdatePlanId)
					var nextPeriodStart = utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)
					var nextPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, nextPeriodStart, plan.Id)

					invoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
						Currency:      pendingUpdate.UpdateCurrency,
						DiscountCode:  discountCode,
						TimeNow:       timeNow,
						PlanId:        pendingUpdate.UpdatePlanId,
						Quantity:      pendingUpdate.UpdateQuantity,
						AddonJsonData: pendingUpdate.UpdateAddonData,
						TaxPercentage: taxPercentage,
						PeriodStart:   nextPeriodStart,
						PeriodEnd:     nextPeriodEnd,
						InvoiceName:   "SubscriptionDowngrade",
						FinishTime:    timeNow,
					})
				} else {
					//generate cycle invoice from sub
					plan := query.GetPlanById(ctx, sub.PlanId)

					var nextPeriodStart = utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)
					var nextPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, nextPeriodStart, plan.Id)

					invoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
						Currency:      sub.Currency,
						DiscountCode:  discountCode,
						TimeNow:       timeNow,
						PlanId:        sub.PlanId,
						Quantity:      sub.Quantity,
						AddonJsonData: sub.AddonData,
						TaxPercentage: taxPercentage,
						PeriodStart:   nextPeriodStart,
						PeriodEnd:     nextPeriodEnd,
						InvoiceName:   "SubscriptionCycle",
						FinishTime:    timeNow,
					})
				}
				if sub.TrialEnd > 0 && sub.TrialEnd == sub.CurrentPeriodEnd {
					invoice.TrialEnd = -2 // mark this invoice is the first invoice after trial
				}
				one, err := handler2.CreateProcessingInvoiceForSub(ctx, invoice, sub)
				if err != nil {
					g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice CreateProcessingInvoiceForSub err:", err.Error())
					return nil, err
				}
				if pendingUpdate != nil {
					_, err = dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
						dao.SubscriptionPendingUpdate.Columns().GmtModify: gtime.Now(),
						dao.SubscriptionPendingUpdate.Columns().InvoiceId: one.InvoiceId,
					}).Where(dao.SubscriptionPendingUpdate.Columns().PendingUpdateId, pendingUpdate.PendingUpdateId).OmitNil().Update()
					if err != nil {
						return nil, err
					}
				}
				g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice CreateProcessingInvoiceForSub:", utility.MarshalToJsonString(one))
				return &BillingCycleWalkRes{WalkUnfinished: true, Message: fmt.Sprintf("Subscription Generate Invoice Result:%s", utility.MarshalToJsonString(one))}, nil
			} else {
				if latestInvoice != nil && latestInvoice.Status == consts.InvoiceStatusProcessing {
					trackForSubscriptionLatestProcessInvoice(ctx, sub, timeNow)
				}
				var lastAutomaticTryTime int64 = 0
				var lastPayment *entity.Payment
				if len(latestInvoice.PaymentId) > 0 {
					lastPayment = query.GetPaymentByPaymentId(ctx, latestInvoice.PaymentId)
					if lastPayment != nil {
						lastAutomaticTryTime = lastPayment.CreateTime
					}
				}
				if latestInvoice != nil && (timeNow-lastAutomaticTryTime) > 86400 && latestInvoice.Status == consts.InvoiceStatusProcessing && needInvoiceFirstTryPayment {
					// finish the payment
					gatewayId, paymentMethodId := user.VerifyPaymentGatewayMethod(ctx, sub.UserId, nil, "", sub.SubscriptionId)
					createRes, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, paymentMethodId, latestInvoice, gatewayId, false, "", "SubscriptionBillingCycle")
					if err != nil {
						g.Log().Print(ctx, "AutomaticPaymentByCycle CreateSubInvoicePaymentDefaultAutomatic err:", err.Error())
						return nil, err
					}

					if createRes.Payment != nil && createRes.Status == consts.PaymentSuccess {
						_ = handler.HandleSubscriptionNextBillingCyclePaymentSuccess(ctx, sub, createRes.Payment)
					}
					return &BillingCycleWalkRes{WalkUnfinished: true, Message: fmt.Sprintf("Subscription Finish Invoice Payment Result:%s", utility.MarshalToJsonString(createRes))}, nil
				} else {
					return &BillingCycleWalkRes{WalkUnfinished: false, Message: "Nothing Todo, Seems Invoice Does not Need Generate"}, nil
				}
			}
		}
	} else {
		g.Log().Print(ctx, source, "GetLock Failure", key)
		return &BillingCycleWalkRes{WalkUnfinished: false, Message: "Sub Get Lock Failure"}, nil
	}
}

// trackForSubscriptionLatestProcessInvoice dunning system for subscription invoice
func trackForSubscriptionLatestProcessInvoice(ctx context.Context, sub *entity.Subscription, timeNow int64) {
	g.Log().Infof(ctx, "trackForSubscriptionLatestProcessInvoice sub:%s", sub.SubscriptionId)
	one := query.GetInvoiceByInvoiceId(ctx, sub.LatestInvoiceId)
	g.Log().Infof(ctx, "trackForSubscriptionLatestProcessInvoice invoiceId:%s", one.InvoiceId)
	if one != nil && one.Status == consts.InvoiceStatusProcessing && one.LastTrackTime+86400 < timeNow {
		// dunning: daily resend invoice, update track time
		g.Log().Infof(ctx, "trackForSubscriptionLatestProcessInvoice start track invoiceId:%s", one.InvoiceId)
		_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().LastTrackTime: timeNow,
			dao.Invoice.Columns().GmtModify:     gtime.Now(),
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			fmt.Printf("trackForSubscriptionLatestProcessInvoice update err:%s", err.Error())
		}
		dayLeft := int((sub.CurrentPeriodEnd - timeNow + 7200) / 86400)
		subscription3.SendMerchantSubscriptionWebhookBackground(sub, dayLeft, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_INVOICE_TRACK)
	}
}

// trackForSubscription dunning system for subscription
func trackForSubscription(ctx context.Context, one *entity.Subscription, timeNow int64) {
	g.Log().Infof(ctx, "trackForSubscription sub:%s", one.SubscriptionId)
	if one.LastTrackTime+86400 < timeNow {
		// dunning: daily resend invoice, update track time
		g.Log().Infof(ctx, "trackForSubscription start track SubscriptionId:%s", one.SubscriptionId)
		_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().LastTrackTime: timeNow,
			dao.Subscription.Columns().GmtModify:     gtime.Now(),
		}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			fmt.Printf("trackForSubscription update err:%s", err.Error())
		}
		dayLeft := int((one.LastUpdateTime - timeNow + 7200) / 86400)
		subscription3.SendMerchantSubscriptionWebhookBackground(one, dayLeft, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_TRACK)
		if one.CancelAtPeriodEnd == 1 && (one.Status != consts.SubStatusCancelled && one.Status != consts.SubStatusExpired) {
			subscription3.SendMerchantSubscriptionWebhookBackground(one, dayLeft, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_TRACK_WILLCANCEL)
		} else if one.Status == consts.SubStatusExpired || one.Status == consts.SubStatusCancelled {
			if query.GetLatestActiveOrIncompleteOrCreateSubscriptionByUserId(ctx, one.UserId, one.MerchantId) == nil {
				key := fmt.Sprintf("UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_TRACK_USER_OUTOFSUBSCRIBE-%d-%d", one.MerchantId, one.UserId)
				if config2.GetConfigInstance().IsProd() {
					if utility.TryLock(ctx, key, 1800) {
						g.Log().Print(ctx, "UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_TRACK_USER_OUTOFSUBSCRIBE-%s-%s", "GetLock 1800s", key)
						subscription3.SendMerchantSubscriptionWebhookBackground(one, dayLeft, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_TRACK_USER_OUTOFSUBSCRIBE)
					}
				} else {
					g.Log().Print(ctx, "UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_TRACK_USER_OUTOFSUBSCRIBE-%s-%s", "GetLock 1800s", key)
					subscription3.SendMerchantSubscriptionWebhookBackground(one, dayLeft, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_TRACK_USER_OUTOFSUBSCRIBE)
				}
			}
		}
	}
}

func isSubscriptionExpireExcludePending(ctx context.Context, sub *entity.Subscription, timeNow int64) bool {
	if timeNow > utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)+config.GetMerchantSubscriptionConfig(ctx, sub.MerchantId).IncompleteExpireTime {
		// expire after periodEnd or trialEnd, depends on incompleteExpireTime config
		return true
	} else if sub.Status == consts.SubStatusIncomplete && sub.CurrentPeriodPaid != 1 && timeNow > sub.CurrentPeriodPaid {
		// manual set sub status to incomplete for several days
		return true
	}
	return false
}
