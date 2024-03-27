package cycle

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	dao "unibee/internal/dao/oversea_pay"
	handler2 "unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/invoice/invoice_compute"
	"unibee/internal/logic/payment/service"
	subscription2 "unibee/internal/logic/subscription"
	"unibee/internal/logic/subscription/billingcycle/expire"
	"unibee/internal/logic/subscription/config"
	"unibee/internal/logic/subscription/handler"
	service2 "unibee/internal/logic/subscription/service"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type BillingCycleWalkRes struct {
	WalkHasDeal bool
	Message     string
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
		var needInvoiceGenerate = true
		var needInvoiceFirstTryPayment = false
		if latestInvoice != nil && (latestInvoice.Status == consts.InvoiceStatusProcessing) {
			needInvoiceGenerate = false
			if len(latestInvoice.PaymentId) == 0 && timeNow > utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd) {
				needInvoiceFirstTryPayment = true
			}
		} else if latestInvoice != nil && latestInvoice.Status == consts.InvoiceStatusPaid && timeNow < latestInvoice.PeriodStart {
			needInvoiceGenerate = false
		} else if timeNow < sub.DunningTime {
			needInvoiceGenerate = false
		}

		if sub.Status == consts.SubStatusExpired || sub.Status == consts.SubStatusCancelled {
			return &BillingCycleWalkRes{WalkHasDeal: false, Message: "Nothing Todo As Sub Cancelled Or Expired"}, nil
		} else if sub.Status == consts.SubStatusCreate {
			if utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)+(2*24*60*60) < timeNow {
				// first time create sub expired
				err := expire.SubscriptionExpire(ctx, sub, "NotPayAfter48Hours")
				if err != nil {
					g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice SubscriptionExpire SubStatus:Created", err.Error())
					return nil, err
				} else {
					return &BillingCycleWalkRes{WalkHasDeal: true, Message: "SubscriptionExpire From Create Status As Payment Out Of 2 Days"}, nil
				}
			} else {
				return &BillingCycleWalkRes{WalkHasDeal: false, Message: "Nothing Todo As Sub At Create Status NotPayBefore48Hours"}, nil
			}
		} else if timeNow > utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd) && sub.CancelAtPeriodEnd == 1 && sub.Status != consts.SubStatusCancelled {
			// sub set cancelAtPeriodEnd, need cancel by system
			needInvoiceGenerate = false
			err = service2.SubscriptionCancel(ctx, sub.SubscriptionId, false, false, "CancelAtPeriodEndBySystem")
			if err != nil {
				g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice SubscriptionCancel err:", err.Error())
				return nil, err
			} else {
				return &BillingCycleWalkRes{WalkHasDeal: true, Message: "SubscriptionCancel At Billing Cycle End By CurrentPeriodEnd Set"}, nil
			}
		} else if !needInvoiceGenerate && !needInvoiceFirstTryPayment && timeNow > utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)+config.GetMerchantSubscriptionConfig(ctx, sub.MerchantId).IncompleteExpireTime {
			// invoice not generate and sub out of time, need expired by system
			err = expire.SubscriptionExpire(ctx, sub, "CycleExpireWithoutPay")
			if err != nil {
				g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice SubscriptionExpire", err.Error())
				return nil, err
			} else {
				return &BillingCycleWalkRes{WalkHasDeal: true, Message: "SubscriptionExpire From Billing Cycle As Payment Out Of Permission Days"}, nil
			}
		} else {
			if sub.CancelAtPeriodEnd == 1 {
				return &BillingCycleWalkRes{WalkHasDeal: false, Message: "Nothing Todo As CancelPeriodEnd Set"}, nil
			}
			// Unpaid after period end or trial end
			if utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd) < timeNow && sub.Status != consts.SubStatusIncomplete {
				err = handler.HandleSubscriptionIncomplete(ctx, sub.SubscriptionId, timeNow)
				if err != nil {
					g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice HandleSubscriptionIncomplete err:", err.Error())
					return nil, err
				} else {
					return &BillingCycleWalkRes{WalkHasDeal: true, Message: "HandleSubscriptionIncomplete As Not Paid After CurrentPeriodEnd Or TrialEnd"}, nil
				}
			}

			if needInvoiceGenerate {
				var invoice *bean.InvoiceSimplify
				pendingUpdate := query.GetUnfinishedSubscriptionPendingUpdateByPendingUpdateId(ctx, sub.PendingUpdateId)
				if pendingUpdate != nil {
					//generate PendingUpdate cycle invoice
					plan := query.GetPlanById(ctx, pendingUpdate.UpdatePlanId)
					var nextPeriodStart = sub.CurrentPeriodEnd
					if sub.TrialEnd > sub.CurrentPeriodEnd {
						nextPeriodStart = sub.TrialEnd
					}
					var nextPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, nextPeriodStart, plan.Id)
					invoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
						Currency:      pendingUpdate.UpdateCurrency,
						PlanId:        pendingUpdate.UpdatePlanId,
						Quantity:      pendingUpdate.UpdateQuantity,
						AddonJsonData: pendingUpdate.UpdateAddonData,
						TaxScale:      sub.TaxScale,
						PeriodStart:   nextPeriodStart,
						PeriodEnd:     nextPeriodEnd,
						InvoiceName:   "SubscriptionDowngrade",
					})
				} else {
					//generate cycle invoice from sub
					plan := query.GetPlanById(ctx, sub.PlanId)

					var nextPeriodStart = sub.CurrentPeriodEnd
					if sub.TrialEnd > sub.CurrentPeriodEnd {
						nextPeriodStart = sub.TrialEnd
					}
					var nextPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, nextPeriodStart, plan.Id)

					invoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
						Currency:      sub.Currency,
						PlanId:        sub.PlanId,
						Quantity:      sub.Quantity,
						AddonJsonData: sub.AddonData,
						TaxScale:      sub.TaxScale,
						PeriodStart:   nextPeriodStart,
						PeriodEnd:     nextPeriodEnd,
						InvoiceName:   "SubscriptionCycle",
					})
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
					}).Where(dao.SubscriptionPendingUpdate.Columns().UpdateSubscriptionId, pendingUpdate.UpdateSubscriptionId).OmitNil().Update()
					if err != nil {
						return nil, err
					}
				}
				g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice CreateProcessingInvoiceForSub:", utility.MarshalToJsonString(one))
				return &BillingCycleWalkRes{WalkHasDeal: true, Message: fmt.Sprintf("Subscription Generate Invoice Result:%s", utility.MarshalToJsonString(one))}, nil
			} else {
				if latestInvoice != nil && latestInvoice.Status == consts.InvoiceStatusProcessing {
					trackForSubscriptionLatest(ctx, sub, timeNow)
				}
				if latestInvoice != nil && len(latestInvoice.PaymentId) == 0 && latestInvoice.Status == consts.InvoiceStatusProcessing && sub.CurrentPeriodEnd < timeNow {
					// finish the payment
					createRes, err := service.CreateSubInvoiceAutomaticPayment(ctx, sub, latestInvoice)
					if err != nil {
						g.Log().Print(ctx, "EndTrialManual CreateSubInvoiceAutomaticPayment err:", err.Error())
						return nil, err
					}
					payment := query.GetPaymentByPaymentId(ctx, createRes.PaymentId)
					if payment != nil && createRes.Status == consts.PaymentSuccess {
						_ = handler.HandleSubscriptionNextBillingCyclePaymentSuccess(ctx, sub, payment)
					}
					return &BillingCycleWalkRes{WalkHasDeal: true, Message: fmt.Sprintf("Subscription Finish Invoice Payment Result:%s", utility.MarshalToJsonString(createRes))}, nil
				} else {
					return &BillingCycleWalkRes{WalkHasDeal: false, Message: "Nothing Todo, Seems Invoice Does not Need Generate"}, nil
				}
			}
		}
	} else {
		g.Log().Print(ctx, source, "GetLock Failure", key)
		return &BillingCycleWalkRes{WalkHasDeal: false, Message: "Sub Get Lock Failure"}, nil
	}
}

// trackForSubscriptionLatest dunning system for subscription invoice
func trackForSubscriptionLatest(ctx context.Context, sub *entity.Subscription, timeNow int64) {
	g.Log().Infof(ctx, "trackForSubscriptionLatest sub:%s", sub.SubscriptionId)
	one := query.GetInvoiceByInvoiceId(ctx, sub.LatestInvoiceId)
	g.Log().Infof(ctx, "trackForSubscriptionLatest invoiceId:%s", one.InvoiceId)
	if one != nil && one.Status == consts.InvoiceStatusProcessing && one.FinishTime+(one.DayUtilDue*86400) < timeNow {
		// todo mark expire invoice
	} else if one != nil && one.Status == consts.InvoiceStatusProcessing && one.LastTrackTime+86400 < timeNow && one.FinishTime+(one.DayUtilDue*86400) > timeNow {
		// dunning: daily resend invoice, update track time
		g.Log().Infof(ctx, "trackForSubscriptionLatest start track invoiceId:%s", one.InvoiceId)
		_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().LastTrackTime: gtime.Now().Timestamp(),
			dao.Invoice.Columns().GmtModify:     gtime.Now(),
		}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			fmt.Printf("SendInvoiceEmailToUser update err:%s", err.Error())
		}
		//err = handler2.SendInvoiceEmailToUser(ctx, one.InvoiceId)
		//if err != nil {
		//	glog.Errorf(ctx, "trackForSubscriptionLatest error:%s", err.Error())
		//}
		subscription3.SendMerchantSubscriptionWebhookBackground(sub, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_INVOICE_TRACK)
	}
}
