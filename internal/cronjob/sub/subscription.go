package sub

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	v1 "go-oversea-pay/api/open/payment"
	redismq2 "go-oversea-pay/internal/cmd/redismq"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/channel/ro"
	"go-oversea-pay/internal/logic/invoice/invoice_compute"
	handler2 "go-oversea-pay/internal/logic/payment/handler"
	"go-oversea-pay/internal/logic/payment/service"
	subscription2 "go-oversea-pay/internal/logic/subscription"
	service2 "go-oversea-pay/internal/logic/subscription/service"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/redismq"
	"go-oversea-pay/utility"
	"strconv"
	"time"
)

func mainTask(ctx context.Context) {
	//3 分钟 Invoice 未支付邮件提醒
	//subscription cycle 支付邮件
	//invoice 保留 3 天时间，每天到点提醒
}

var (
	SubscriptionDelayPaymentPermissionTime int64 = 24 * 60 * 60 // 24h expire after
)

func SubscriptionBillingCycleDunningInvoice(ctx context.Context, taskName string) {
	g.Log().Print(ctx, taskName, "Start......")
	var timeNow = gtime.Now().Timestamp()
	var subs []*entity.Subscription
	var sortKey = "task_time asc"
	var status = []int{consts.SubStatusActive}
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
			// todo mark need consider trial end
			if sub.CurrentPeriodEnd+SubscriptionDelayPaymentPermissionTime < timeNow {
				// sub out of time, need expired by system
				err := SubscriptionExpire(ctx, sub, "CycleExpireWithoutPay")
				if err != nil {
					g.Log().Print(ctx, taskName, "SubscriptionExpire", err.Error())
				}
			} else {
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
						g.Log().Print(ctx, taskName, "SubscriptionBillingCycleDunningInvoice pay channel not found")
						continue
					}
					merchantInfo := query.GetMerchantInfoById(ctx, sub.MerchantId)
					if merchantInfo == nil {
						g.Log().Print(ctx, taskName, "SubscriptionBillingCycleDunningInvoice merchantInfo not found")
						continue
					}
					createRes, err := service.DoChannelPay(ctx, &ro.CreatePayContext{
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
					if err != nil {
						g.Log().Print(ctx, taskName, "SubscriptionBillingCycleDunningInvoice err:", err.Error())
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

	g.Log().Print(ctx, taskName, "End......")
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
