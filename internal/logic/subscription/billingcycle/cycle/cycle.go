package cycle

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"unibee/api/bean"
	config2 "unibee/internal/cmd/config"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	dao "unibee/internal/dao/default"
	config3 "unibee/internal/logic/credit/config"
	"unibee/internal/logic/discount"
	handler3 "unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/invoice/invoice_compute"
	handler2 "unibee/internal/logic/invoice/service"
	"unibee/internal/logic/payment/service"
	subscription2 "unibee/internal/logic/subscription"
	"unibee/internal/logic/subscription/billingcycle/expire"
	"unibee/internal/logic/subscription/config"
	"unibee/internal/logic/subscription/handler"
	"unibee/internal/logic/subscription/pending_update_cancel"
	service2 "unibee/internal/logic/subscription/service"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/logic/user/vat"
	entity "unibee/internal/model/entity/default"
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
	plan := query.GetPlanById(ctx, sub.PlanId)
	if plan == nil {
		return &BillingCycleWalkRes{Message: "Plan Not Found"}, nil
	}
	key := fmt.Sprintf("SubscriptionCycleWalk-%s", sub.SubscriptionId)
	if utility.TryLock(ctx, key, 60) {
		g.Log().Debugf(ctx, source, "GetLock 60s", key)
		defer func() {
			utility.ReleaseLock(ctx, key)
			g.Log().Debugf(ctx, source, "ReleaseLock", key)
		}()
		var err error
		_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().TaskTime: gtime.Now(),
		}).Where(dao.Subscription.Columns().Id, sub.Id).OmitNil().Update()
		if err != nil {
			g.Log().Errorf(ctx, source, "SubscriptionBillingCycleDunningInvoice Update TaskTime err:", err.Error())
		}

		if len(sub.PendingUpdateId) > 0 {
			pendingUpdate := query.GetSubscriptionPendingUpdateByPendingUpdateId(ctx, sub.PendingUpdateId)
			if pendingUpdate.EffectImmediate == 1 && pendingUpdate.Status < consts.PendingSubStatusFinished && pendingUpdate.CreateTime+3600 < timeNow { // one hour
				err = pending_update_cancel.SubscriptionPendingUpdateCancel(ctx, pendingUpdate.PendingUpdateId, "EffectTimeout")
				if err != nil {
					g.Log().Errorf(ctx, source, "SubPipeBillingCycleWalk SubscriptionPendingUpdateCancel pendingUpdateId:%s err:", pendingUpdate.PendingUpdateId, err.Error())
				}
			}
		}

		// generate invoice and payment ahead
		latestInvoice := query.GetInvoiceByInvoiceId(ctx, sub.LatestInvoiceId)

		trackForSubscription(ctx, sub, timeNow)

		var lastAutomaticTryTime int64 = 0
		var lastPayment *entity.Payment
		if latestInvoice != nil && len(latestInvoice.PaymentId) > 0 {
			lastPayment = query.GetPaymentByPaymentId(ctx, latestInvoice.PaymentId)
			if lastPayment != nil {
				lastAutomaticTryTime = lastPayment.CreateTime
			}
		}
		var needInvoiceGenerate = true
		var needTryInvoiceAutomaticPayment = false
		if latestInvoice != nil && (latestInvoice.Status == consts.InvoiceStatusProcessing) {
			needInvoiceGenerate = false
			if timeNow > utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)-config.GetMerchantSubscriptionConfig(ctx, sub.MerchantId).TryAutomaticPaymentBeforePeriodEnd &&
				(timeNow-lastAutomaticTryTime) > 43200 &&
				latestInvoice.GatewayId > 0 {
				needTryInvoiceAutomaticPayment = true
			}
		} else if latestInvoice != nil && latestInvoice.Status == consts.InvoiceStatusPaid && timeNow < latestInvoice.PeriodStart {
			needInvoiceGenerate = false
		} else if timeNow < sub.DunningTime {
			needInvoiceGenerate = false
		}
		// lock the invoice and payment creation half an hour before period end
		if timeNow > utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)-27*60 && timeNow < utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd) {
			needInvoiceGenerate = false
			needTryInvoiceAutomaticPayment = false
		}
		if plan.DisableAutoCharge > 0 {
			needInvoiceGenerate = false
			needTryInvoiceAutomaticPayment = false
		}

		if sub.Status == consts.SubStatusExpired || sub.Status == consts.SubStatusFailed || sub.Status == consts.SubStatusCancelled {
			return &BillingCycleWalkRes{WalkUnfinished: false, Message: "Nothing Todo As Sub Cancelled Or Expired Or Failed"}, nil
		} else if sub.Status == consts.SubStatusPending || sub.Status == consts.SubStatusProcessing {
			if sub.GmtCreate.Timestamp()+consts.SubPendingTimeout < timeNow {
				// first time create sub expired
				err = expire.SubscriptionExpire(ctx, sub, "NotPayAfter36Hours")
				if err != nil {
					g.Log().Errorf(ctx, source, "SubscriptionBillingCycleDunningInvoice SubscriptionExpire SubStatus:Created", err.Error())
					return nil, err
				} else {
					return &BillingCycleWalkRes{WalkUnfinished: true, Message: "SubscriptionExpire From Create Status As Payment Out Of 2 Days"}, nil
				}
			} else {
				return &BillingCycleWalkRes{WalkUnfinished: false, Message: "Nothing Todo As Sub At Create Status NotPayBefore48Hours"}, nil
			}
		} else if timeNow > utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd) && sub.CancelAtPeriodEnd == 1 && sub.Status != consts.SubStatusCancelled {
			//if latestInvoice != nil && latestInvoice.Status == consts.InvoiceStatusProcessing &&
			//	(latestInvoice.InvoiceName == "SubscriptionUpdate" || latestInvoice.InvoiceName == "SubscriptionRenew") &&
			//	needTryInvoiceAutomaticPayment {
			//	createRes, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, latestInvoice, false, "", "", "SubscriptionBillingCycle", timeNow)
			//	if err != nil {
			//		g.Log().Errorf(ctx, "AutomaticPaymentByCycle CreateSubInvoicePaymentDefaultAutomatic err:", err.Error())
			//		return nil, err
			//	}
			//
			//	if createRes.Payment != nil && createRes.Status == consts.PaymentSuccess {
			//		_ = handler.HandleSubscriptionNextBillingCyclePaymentSuccess(ctx, sub, latestInvoice)
			//		return &BillingCycleWalkRes{WalkUnfinished: true, Message: fmt.Sprintf("Subscription Update Invoice Payment Result:%s", utility.MarshalToJsonString(createRes))}, nil
			//	} else {
			//		_ = handler2.ProcessingInvoiceFailure(ctx, latestInvoice.InvoiceId, "TaskForExpireInvoices")
			//		return &BillingCycleWalkRes{WalkUnfinished: true, Message: fmt.Sprintf("Subscription Update Invoice Payment Result:%s", utility.MarshalToJsonString(createRes))}, nil
			//	}
			//}
			// sub set cancelAtPeriodEnd, need cancel by system
			needInvoiceGenerate = false
			err = service2.SubscriptionCancel(ctx, sub.SubscriptionId, false, false, "CancelAtPeriodEndBySystem")
			if err != nil {
				g.Log().Errorf(ctx, source, "SubscriptionBillingCycleDunningInvoice SubscriptionCancel err:", err.Error())
				return nil, err
			} else {
				return &BillingCycleWalkRes{WalkUnfinished: true, Message: "SubscriptionCancel At Billing Cycle End By CurrentPeriodEnd Set"}, nil
			}
		} else if !needInvoiceGenerate && !needTryInvoiceAutomaticPayment && isSubscriptionExpireExcludePending(ctx, sub, timeNow) {
			// invoice not generate and sub out of time, need expired by system
			err = expire.SubscriptionExpire(ctx, sub, "AutoRenewFailure")
			if err != nil {
				g.Log().Errorf(ctx, source, "SubscriptionBillingCycleDunningInvoice SubscriptionExpire", err.Error())
				return nil, err
			} else {
				_, _ = redismq.Send(&redismq.Message{
					Topic:      redismq2.TopicSubscriptionAutoRenewFailure.Topic,
					Tag:        redismq2.TopicSubscriptionAutoRenewFailure.Tag,
					Body:       sub.SubscriptionId,
					CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
				})
				return &BillingCycleWalkRes{WalkUnfinished: true, Message: "SubscriptionExpire For AutoRenew Failed"}, nil
			}
		} else {
			if sub.CancelAtPeriodEnd == 1 {
				//if latestInvoice != nil && latestInvoice.Status == consts.InvoiceStatusProcessing &&
				//	(latestInvoice.InvoiceName == "SubscriptionUpdate" || latestInvoice.InvoiceName == "SubscriptionRenew") &&
				//	needTryInvoiceAutomaticPayment {
				//	createRes, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, latestInvoice, false, "", "", "SubscriptionBillingCycle", timeNow)
				//	if err != nil {
				//		g.Log().Errorf(ctx, "AutomaticPaymentByCycle CreateSubInvoicePaymentDefaultAutomatic err:", err.Error())
				//		return nil, err
				//	}
				//
				//	if createRes.Payment != nil && createRes.Status == consts.PaymentSuccess {
				//		_ = handler.HandleSubscriptionNextBillingCyclePaymentSuccess(ctx, sub, latestInvoice)
				//		return &BillingCycleWalkRes{WalkUnfinished: true, Message: fmt.Sprintf("Subscription Update Invoice Payment Result:%s", utility.MarshalToJsonString(createRes))}, nil
				//	} else {
				//		_ = handler2.ProcessingInvoiceFailure(ctx, latestInvoice.InvoiceId, "TaskForExpireInvoices")
				//		return &BillingCycleWalkRes{WalkUnfinished: true, Message: fmt.Sprintf("Subscription Update Invoice Payment Result:%s", utility.MarshalToJsonString(createRes))}, nil
				//	}
				//}
				return &BillingCycleWalkRes{WalkUnfinished: false, Message: "Nothing Todo As CancelPeriodEnd Set"}, nil
			}
			// Unpaid after period end or trial end
			if utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd) < timeNow && sub.Status != consts.SubStatusIncomplete {
				err = handler.HandleSubscriptionIncomplete(ctx, sub.SubscriptionId, timeNow)
				if err != nil {
					g.Log().Errorf(ctx, source, "SubscriptionBillingCycleDunningInvoice HandleSubscriptionIncomplete err:", err.Error())
					return nil, err
				} else {
					return &BillingCycleWalkRes{WalkUnfinished: true, Message: "HandleSubscriptionIncomplete As Not Paid After CurrentPeriodEnd Or TrialEnd"}, nil
				}
			}
			var taxPercentage = sub.TaxPercentage
			percentage, countryCode, vatNumber, err := vat.GetUserTaxPercentage(ctx, sub.UserId)
			if err == nil {
				taxPercentage = percentage
			}

			if needInvoiceGenerate {
				var invoice *bean.Invoice
				var discountCode = ""
				canApply, isRecurring, _ := discount.UserDiscountApplyPreview(ctx, &discount.UserDiscountApplyReq{
					MerchantId:         sub.MerchantId,
					UserId:             sub.UserId,
					DiscountCode:       sub.DiscountCode,
					Currency:           sub.Currency,
					SubscriptionId:     sub.SubscriptionId,
					PLanId:             sub.PlanId,
					TimeNow:            timeNow,
					IsRecurringApply:   true,
					IsUpgrade:          false,
					IsChangeToLongPlan: false,
					IsRenew:            false,
					IsNewUser:          false,
				})
				if canApply && isRecurring {
					discountCode = sub.DiscountCode
				}
				pendingUpdate := query.GetUnfinishedSubscriptionPendingUpdateByPendingUpdateId(ctx, sub.PendingUpdateId)
				if pendingUpdate != nil && pendingUpdate.EffectImmediate == 1 {
					pendingUpdate = nil
				}
				applyPromoCredit := config3.CheckCreditConfigRecurring(ctx, sub.MerchantId, consts.CreditAccountTypePromo, sub.Currency)
				if config3.CheckCreditConfigDiscountCodeExclusive(ctx, sub.MerchantId, consts.CreditAccountTypePromo, sub.Currency) && len(discountCode) > 0 {
					applyPromoCredit = false
				}
				if pendingUpdate != nil {
					//generate PendingUpdate cycle invoice
					updatePlan := query.GetPlanById(ctx, pendingUpdate.UpdatePlanId)
					var nextPeriodStart = utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)
					var nextPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, nextPeriodStart, sub.BillingCycleAnchor, updatePlan.Id)

					invoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
						UserId:           sub.UserId,
						Currency:         pendingUpdate.UpdateCurrency,
						DiscountCode:     discountCode,
						TimeNow:          timeNow,
						PlanId:           pendingUpdate.UpdatePlanId,
						Quantity:         pendingUpdate.UpdateQuantity,
						AddonJsonData:    pendingUpdate.UpdateAddonData,
						VatNumber:        vatNumber,
						CountryCode:      countryCode,
						TaxPercentage:    taxPercentage,
						PeriodStart:      nextPeriodStart,
						PeriodEnd:        nextPeriodEnd,
						InvoiceName:      "SubscriptionDowngrade",
						FinishTime:       timeNow,
						CreateFrom:       consts.InvoiceAutoChargeFlag,
						Metadata:         map[string]interface{}{"SubscriptionUpdate": true, "IsUpgrade": false},
						ApplyPromoCredit: applyPromoCredit,
					})
				} else {
					//generate cycle invoice from sub
					plan = query.GetPlanById(ctx, sub.PlanId)

					var nextPeriodStart = utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)
					var nextPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, nextPeriodStart, sub.BillingCycleAnchor, plan.Id)

					invoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
						UserId:           sub.UserId,
						Currency:         sub.Currency,
						DiscountCode:     discountCode,
						TimeNow:          timeNow,
						PlanId:           sub.PlanId,
						Quantity:         sub.Quantity,
						AddonJsonData:    sub.AddonData,
						VatNumber:        vatNumber,
						CountryCode:      countryCode,
						TaxPercentage:    taxPercentage,
						PeriodStart:      nextPeriodStart,
						PeriodEnd:        nextPeriodEnd,
						InvoiceName:      "SubscriptionCycle",
						FinishTime:       timeNow,
						CreateFrom:       consts.InvoiceAutoChargeFlag,
						ApplyPromoCredit: applyPromoCredit,
					})
				}
				if sub.TrialEnd > 0 && sub.TrialEnd == sub.CurrentPeriodEnd {
					invoice.TrialEnd = -2 // mark this invoice is the first invoice after trial
				}
				gatewayId, paymentMethodId := sub_update.VerifyPaymentGatewayMethod(ctx, sub.UserId, nil, "", sub.SubscriptionId)
				if gatewayId > 0 && (gatewayId != sub.GatewayId || paymentMethodId != sub.GatewayDefaultPaymentMethod) {
					_, _ = dao.Subscription.Ctx(ctx).Data(g.Map{
						dao.Subscription.Columns().GmtModify:                   gtime.Now(),
						dao.Subscription.Columns().GatewayId:                   gatewayId,
						dao.Subscription.Columns().GatewayDefaultPaymentMethod: paymentMethodId,
					}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()
					sub.GatewayId = gatewayId
					sub.GatewayDefaultPaymentMethod = paymentMethodId
				}
				one, err := handler2.CreateProcessingInvoiceForSub(ctx, sub.PlanId, invoice, sub, sub.GatewayId, sub.GatewayDefaultPaymentMethod, true, timeNow)
				if err != nil {
					g.Log().Errorf(ctx, source, "SubscriptionBillingCycleDunningInvoice CreateProcessingInvoiceForSub err:%s", err.Error())
					return nil, err
				}
				if pendingUpdate != nil {
					_, err = dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
						dao.SubscriptionPendingUpdate.Columns().GmtModify: gtime.Now(),
						dao.SubscriptionPendingUpdate.Columns().InvoiceId: one.InvoiceId,
					}).Where(dao.SubscriptionPendingUpdate.Columns().PendingUpdateId, pendingUpdate.PendingUpdateId).OmitNil().Update()
					if err != nil {
						g.Log().Errorf(ctx, source, "SubscriptionBillingCycleDunningInvoice update pendingUpdate err:", err.Error())
						return nil, err
					}
				}
				g.Log().Debugf(ctx, source, "SubscriptionBillingCycleDunningInvoice CreateProcessingInvoiceForSub:", utility.MarshalToJsonString(one))
				return &BillingCycleWalkRes{WalkUnfinished: true, Message: fmt.Sprintf("Subscription Generate Invoice Result:%s", utility.MarshalToJsonString(one))}, nil
			} else {
				if sub.CancelAtPeriodEnd == 1 {
					return &BillingCycleWalkRes{WalkUnfinished: false, Message: "Nothing Todo As CancelPeriodEnd Set"}, nil
				}
				if latestInvoice != nil && latestInvoice.Status == consts.InvoiceStatusProcessing {
					trackForSubscriptionLatestProcessInvoice(ctx, sub, timeNow)
				}

				if needTryInvoiceAutomaticPayment {
					// finish the payment
					if latestInvoice.TotalAmount == 0 {
						paidInvoice, err := handler3.MarkInvoiceAsPaidForZeroPayment(ctx, latestInvoice.InvoiceId)
						if err != nil || paidInvoice.Status != consts.InvoiceStatusPaid {
							if err != nil {
								g.Log().Errorf(ctx, "AutomaticPaymentByCycle MarkInvoiceAsPaidForZeroPayment invoiceId:%s err:%s", latestInvoice.InvoiceId, err.Error())
							} else {
								g.Log().Errorf(ctx, "AutomaticPaymentByCycle MarkInvoiceAsPaidForZeroPayment failed invoiceId:%s ", latestInvoice.InvoiceId)
							}
							return nil, err
						}
						_ = handler.HandleSubscriptionNextBillingCyclePaymentSuccess(ctx, sub, latestInvoice)
						return &BillingCycleWalkRes{WalkUnfinished: true, Message: fmt.Sprintf("Subscription Finish Zero Invoice Payment Result:%s", utility.MarshalToJsonString(paidInvoice))}, nil
					} else {
						// gatewayId, paymentMethodId := user.VerifyPaymentGatewayMethod(ctx, sub.UserId, nil, "", sub.SubscriptionId)
						createRes, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, latestInvoice, false, "", "", "SubscriptionBillingCycle", timeNow)
						if err != nil {
							g.Log().Errorf(ctx, "AutomaticPaymentByCycle CreateSubInvoicePaymentDefaultAutomatic err:%s", err.Error())
							return nil, err
						}

						if createRes.Payment != nil && createRes.Status == consts.PaymentSuccess {
							_ = handler.HandleSubscriptionNextBillingCyclePaymentSuccess(ctx, sub, latestInvoice)
						}
						return &BillingCycleWalkRes{WalkUnfinished: true, Message: fmt.Sprintf("Subscription Finish Invoice Payment Result:%s", utility.MarshalToJsonString(createRes))}, nil
					}
				} else if latestInvoice != nil && latestInvoice.GatewayId <= 0 {
					return &BillingCycleWalkRes{WalkUnfinished: false, Message: "Nothing Todo, Seems Invoice Gateway Need Specified"}, nil
				} else {
					return &BillingCycleWalkRes{WalkUnfinished: false, Message: "Nothing Todo, Seems Invoice Does not Need Generate"}, nil
				}
			}
		}
	} else {
		g.Log().Debugf(ctx, source, "GetLock Failure", key)
		return &BillingCycleWalkRes{WalkUnfinished: false, Message: "Sub Get Lock Failure"}, nil
	}
}

func trackForSubscriptionLatestProcessInvoice(ctx context.Context, sub *entity.Subscription, timeNow int64) {
	//g.Log().Infof(ctx, "trackForSubscriptionLatestProcessInvoice sub:%s", sub.SubscriptionId)
	//one := query.GetInvoiceByInvoiceId(ctx, sub.LatestInvoiceId)
	//g.Log().Infof(ctx, "trackForSubscriptionLatestProcessInvoice invoiceId:%s", one.InvoiceId)
	//if one != nil && one.Status == consts.InvoiceStatusProcessing && one.LastTrackTime+86400 < timeNow {
	//	// dunning: daily resend invoice, update track time
	//	g.Log().Infof(ctx, "trackForSubscriptionLatestProcessInvoice start track invoiceId:%s", one.InvoiceId)
	//	_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
	//		dao.Invoice.Columns().LastTrackTime: timeNow,
	//		dao.Invoice.Columns().GmtModify:     gtime.Now(),
	//	}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
	//	if err != nil {
	//		fmt.Printf("trackForSubscriptionLatestProcessInvoice update err:%s", err.Error())
	//	}
	//	dayLeft := int((sub.CurrentPeriodEnd - timeNow + 7200) / 86400)
	//	subscription3.SendMerchantSubscriptionWebhookBackground(sub, dayLeft, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_INVOICE_TRACK)
	//}
}

func trackForSubscription(ctx context.Context, one *entity.Subscription, timeNow int64) {
	g.Log().Debugf(ctx, "trackForSubscription sub:%s", one.SubscriptionId)
	plan := query.GetPlanById(ctx, one.PlanId)
	if plan == nil {
		return
	}
	if one.LastTrackTime+86400 <= timeNow && (timeNow-one.CurrentPeriodEnd) >= -15*86400 && (timeNow-one.CurrentPeriodEnd) <= 15*86400 {
		// dunning: daily resend invoice, update track time
		g.Log().Debugf(ctx, "trackForSubscription start track SubscriptionId:%s", one.SubscriptionId)
		_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().LastTrackTime: timeNow,
			dao.Subscription.Columns().GmtModify:     gtime.Now(),
		}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			g.Log().Errorf(ctx, "trackForSubscription update err:%s", err.Error())
		}
		dayLeft := int((one.CurrentPeriodEnd - timeNow + 7200) / 86400)
		subscription3.SendMerchantSubscriptionWebhookBackground(one, dayLeft, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_TRACK, nil)
		if one.CancelAtPeriodEnd == 1 && (one.Status != consts.SubStatusCancelled && one.Status != consts.SubStatusExpired) {
			subscription3.SendMerchantSubscriptionWebhookBackground(one, dayLeft, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_TRACK_WILLCANCEL, nil)
		} else if one.Status == consts.SubStatusExpired || one.Status == consts.SubStatusCancelled {
			if query.GetLatestActiveOrIncompleteOrCreateSubscriptionByUserId(ctx, one.UserId, one.MerchantId, plan.ProductId) == nil {
				key := fmt.Sprintf("UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_TRACK_USER_OUTOFSUBSCRIBE-%d-%d", one.MerchantId, one.UserId)
				if config2.GetConfigInstance().IsProd() {
					if utility.TryLock(ctx, key, 1800) {
						g.Log().Debugf(ctx, "UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_TRACK_USER_OUTOFSUBSCRIBE-%s-%s", "GetLock 1800s", key)
						subscription3.SendMerchantSubscriptionWebhookBackground(one, dayLeft, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_TRACK_USER_OUTOFSUBSCRIBE, nil)
					}
				} else {
					g.Log().Debugf(ctx, "UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_TRACK_USER_OUTOFSUBSCRIBE-%s-%s", "GetLock 1800s", key)
					subscription3.SendMerchantSubscriptionWebhookBackground(one, dayLeft, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_TRACK_USER_OUTOFSUBSCRIBE, nil)
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
