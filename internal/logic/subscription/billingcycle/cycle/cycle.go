package cycle

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/gateway/ro"
	"unibee/internal/logic/invoice/invoice_compute"
	"unibee/internal/logic/payment/service"
	subscription2 "unibee/internal/logic/subscription"
	"unibee/internal/logic/subscription/billingcycle/expire"
	"unibee/internal/logic/subscription/handler"
	service2 "unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"
)

var (
	SubscriptionCycleDelayPaymentPermissionTime int64 = 24 * 60 * 60 // 24h expire after
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
		if sub.Status == consts.SubStatusCreate {
			if utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)+(2*24*60*60) < timeNow {
				err := expire.SubscriptionExpire(ctx, sub, "NotPayAfter48Hours")
				if err != nil {
					g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice SubscriptionExpire SubStatus:Created", err.Error())
					return nil, err
				} else {
					return &BillingCycleWalkRes{WalkHasDeal: true, Message: "SubscriptionExpire From Create Status As Payment Out Of 2 Days"}, nil
				}
			} else {
				return &BillingCycleWalkRes{WalkHasDeal: false, Message: "Nothing Todo As Sub Get Lock Failure"}, nil
			}
		}
		if utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)+SubscriptionCycleDelayPaymentPermissionTime < timeNow {
			// sub out of time, need expired by system
			err := expire.SubscriptionExpire(ctx, sub, "CycleExpireWithoutPay")
			if err != nil {
				g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice SubscriptionExpire", err.Error())
				return nil, err
			} else {
				return &BillingCycleWalkRes{WalkHasDeal: true, Message: "SubscriptionExpire From Billing Cycle As Payment Out Of Permission Days"}, nil
			}
		} else if utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd) < timeNow {
			if sub.CancelAtPeriodEnd == 1 {
				// Cancel At Period End
				err = service2.SubscriptionCancel(ctx, sub.SubscriptionId, false, false, "CancelAtPeriodEndBySystem")
				if err != nil {
					g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice SubscriptionCancel err:", err.Error())
					return nil, err
				} else {
					return &BillingCycleWalkRes{WalkHasDeal: true, Message: "SubscriptionCancel At Billing Cycle End By CurrentPeriodEnd Set"}, nil
				}
			} else {
				// Not Paid To Incomplete
				err = handler.SubscriptionIncomplete(ctx, sub.SubscriptionId, timeNow)
				if err != nil {
					g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice SubscriptionIncomplete err:", err.Error())
					return nil, err
				} else {
					return &BillingCycleWalkRes{WalkHasDeal: true, Message: "SubscriptionIncomplete As Not Paid After CurrentPeriodEnd Or TrialEnd"}, nil
				}
			}
		} else {
			if sub.CancelAtPeriodEnd == 1 {
				return &BillingCycleWalkRes{WalkHasDeal: false, Message: "Nothing Todo As CurrentPeriodEnd Set"}, nil
			}
			// generate invoice and payment ahead
			latestInvoice := query.GetInvoiceByInvoiceId(ctx, sub.LatestInvoiceId)
			var needGenerate = true
			if latestInvoice != nil && (latestInvoice.Status == consts.InvoiceStatusProcessing || latestInvoice.Status == consts.InvoiceStatusPending) {
				needGenerate = false
			} else if latestInvoice != nil && latestInvoice.Status == consts.InvoiceStatusPaid && latestInvoice.PeriodEnd > sub.CurrentPeriodEnd {
				needGenerate = false
			}
			if needGenerate {
				var invoice *ro.InvoiceDetailSimplify
				var billingReason = ""
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
					})
					billingReason = "SubscriptionDowngrade"
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
					})
					billingReason = "SubscriptionCycle"
				}
				createRes, err := service.CreateSubInvoicePayment(ctx, sub, invoice, billingReason, true)
				if err != nil {
					g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice CreateSubInvoicePayment err:", err.Error())
					return nil, err
				}
				g.Log().Print(ctx, source, "SubscriptionBillingCycleDunningInvoice GatewayPaymentCreate:", utility.MarshalToJsonString(createRes))
				return &BillingCycleWalkRes{WalkHasDeal: true, Message: fmt.Sprintf("Subscription Generate Invoice And Try Payment Result:%s", utility.MarshalToJsonString(createRes))}, nil
			} else {
				return &BillingCycleWalkRes{WalkHasDeal: true, Message: "Nothing Todo, Seems Invoice Does not Need Generate"}, nil
			}
		}
	} else {
		g.Log().Print(ctx, source, "GetLock Failure", key)
		return &BillingCycleWalkRes{WalkHasDeal: false, Message: "Sub Get Lock Failure"}, nil
	}
}
