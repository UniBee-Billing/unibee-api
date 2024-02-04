package sub

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq2 "go-oversea-pay/internal/cmd/redismq"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/channel/ro"
	"go-oversea-pay/internal/logic/invoice/invoice_compute"
	handler2 "go-oversea-pay/internal/logic/payment/handler"
	"go-oversea-pay/internal/logic/payment/service"
	subscription2 "go-oversea-pay/internal/logic/subscription"
	"go-oversea-pay/internal/logic/subscription/handler"
	service2 "go-oversea-pay/internal/logic/subscription/service"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/redismq"
	"go-oversea-pay/utility"
	"time"
)

func mainTask(ctx context.Context) {
	//3 分钟 Invoice 未支付邮件提醒
	//subscription cycle 支付邮件
	//invoice 保留 3 天时间，每天到点提醒
}

var (
	SubscriptionCycleDelayPaymentPermissionTime int64 = 24 * 60 * 60 // 24h expire after
)

func SubscriptionBillingCycleDunningInvoice(ctx context.Context, taskName string) {
	g.Log().Debug(ctx, taskName, "Start......")
	var timeNow = gtime.Now().Timestamp()
	var subs []*entity.Subscription
	var sortKey = "task_time asc"
	var status = []int{consts.SubStatusCreate, consts.SubStatusActive, consts.SubStatusIncomplete}
	// query sub which dunningTime expired
	err := dao.Subscription.Ctx(ctx).
		Where(dao.Subscription.Columns().IsDeleted, 0).
		WhereLT(dao.Subscription.Columns().DunningTime, timeNow). //  dunning < now
		Where(dao.Subscription.Columns().Type, consts.SubTypeUniBeeControl).
		WhereIn(dao.Subscription.Columns().Status, status).
		Limit(0, 10).
		Order(sortKey).
		OmitEmpty().Scan(&subs)
	if err != nil {
		g.Log().Errorf(ctx, "%s Error:%s", taskName, err.Error())
		return
	}

	for _, sub := range subs {
		key := fmt.Sprintf("SubscriptionCycle-%s", sub.SubscriptionId)
		if utility.TryLock(ctx, key, 60) {
			g.Log().Print(ctx, taskName, "GetLock 60s", key)
			if sub.Status == consts.SubStatusCreate {
				if utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)+(2*24*60*60) < timeNow {
					err := SubscriptionExpire(ctx, sub, "NotPayAfter48Hours")
					if err != nil {
						g.Log().Print(ctx, taskName, "SubscriptionBillingCycleDunningInvoice SubscriptionExpire SubStatus:Created", err.Error())
					}
				}
				continue
			}
			if utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)+SubscriptionCycleDelayPaymentPermissionTime < timeNow {
				// sub out of time, need expired by system
				err := SubscriptionExpire(ctx, sub, "CycleExpireWithoutPay")
				if err != nil {
					g.Log().Print(ctx, taskName, "SubscriptionBillingCycleDunningInvoice SubscriptionExpire", err.Error())
				}
			} else if utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd) < timeNow {
				if sub.CancelAtPeriodEnd == 1 {
					// Cancel At Period End
					err = service2.SubscriptionCancel(ctx, sub.SubscriptionId, false, false, "CancelAtPeriodEndBySystem")
					if err != nil {
						g.Log().Print(ctx, taskName, "SubscriptionBillingCycleDunningInvoice SubscriptionCancel err:", err.Error())
					}
				} else {
					// Not Paid To Incomplete
					err = handler.SubscriptionIncomplete(ctx, sub.SubscriptionId)
					if err != nil {
						g.Log().Print(ctx, taskName, "SubscriptionBillingCycleDunningInvoice SubscriptionIncomplete err:", err.Error())
					}
				}
			} else {
				if sub.CancelAtPeriodEnd == 1 {
					continue
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
					createRes, err := service.CreateSubInvoicePayment(ctx, sub, invoice, billingReason)
					if err != nil {
						g.Log().Print(ctx, taskName, "SubscriptionBillingCycleDunningInvoice CreateSubInvoicePayment err:", err.Error())
						continue
					}
					g.Log().Print(ctx, taskName, "SubscriptionBillingCycleDunningInvoice DoChannelPay:", utility.MarshalToJsonString(createRes))
					_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
						dao.Subscription.Columns().TaskTime: gtime.Now(),
					}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
					if err != nil {
						g.Log().Print(ctx, taskName, "SubscriptionBillingCycleDunningInvoice UpdateTaskTime err:", err.Error())
					}
				}
			}
			// compute cycle
			time.Sleep(10 * time.Second)
			utility.ReleaseLock(ctx, key)
			g.Log().Print(ctx, taskName, "ReleaseLock", key)
		} else {
			g.Log().Print(ctx, taskName, "GetLock Failure", key)
		}
	}

	g.Log().Debug(ctx, taskName, "End......")
}

func SubscriptionExpire(ctx context.Context, sub *entity.Subscription, reason string) error {
	//Expire SubscriptionPendingUpdate
	var pendingUpdates []*entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().SubscriptionId, sub.SubscriptionId).
		WhereLT(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusFinished).
		Limit(0, 100).
		OmitEmpty().Scan(&pendingUpdates)
	if err != nil {
		return err
	}
	for _, p := range pendingUpdates {
		err = service2.SubscriptionPendingUpdateCancel(ctx, p.UpdateSubscriptionId, reason)
		if err != nil {
			fmt.Printf("MakeSubscriptionExpired SubscriptionPendingUpdateCancel error:%s", err.Error())
		}
	}
	//Cancel Subscription Remaining Payment
	var paymentList []*entity.Payment
	err = dao.Payment.Ctx(ctx).
		Where(dao.Payment.Columns().SubscriptionId, sub.SubscriptionId).
		Where(dao.Payment.Columns().Status, consts.TO_BE_PAID).
		Limit(0, 100).
		OmitEmpty().Scan(&paymentList)
	if err != nil {
		fmt.Printf("SubscriptionExpire GetPaymentList error:%s", err.Error())
	}
	for _, p := range paymentList {
		// todo mark should use DoChannelCancel
		err := handler2.HandlePayCancel(ctx, &handler2.HandlePayReq{
			PaymentId:     p.PaymentId,
			PayStatusEnum: consts.PAY_CANCEL,
			Reason:        reason,
		})
		if err != nil {
			fmt.Printf("SubscriptionExpire HandlePayCancel error:%s", err.Error())
		}
	}
	//Expire Subscription UnFinished Invoice, May No Need
	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:       consts.SubStatusExpired,
		dao.Subscription.Columns().CancelReason: reason,
		dao.Subscription.Columns().TrialEnd:     sub.CurrentPeriodStart - 1,
		dao.Subscription.Columns().GmtModify:    gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()
	if err != nil {
		fmt.Printf("SubscriptionExpire error:%s", err.Error())
		return err
	}

	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicSubscriptionExpire.Topic,
		Tag:   redismq2.TopicSubscriptionExpire.Tag,
		Body:  sub.SubscriptionId,
	})

	return nil
}
