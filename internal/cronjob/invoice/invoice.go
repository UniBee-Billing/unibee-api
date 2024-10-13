package invoice

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"time"
	"unibee/api/bean"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/invoice/detail"
	"unibee/internal/logic/invoice/service"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func TaskForExpireInvoices(ctx context.Context) {
	var list []*entity.Invoice
	err := dao.Invoice.Ctx(ctx).
		Where(dao.Invoice.Columns().Status, consts.InvoiceStatusProcessing).
		Where(dao.Invoice.Columns().IsDeleted, 0).
		OrderAsc(dao.Invoice.Columns().FinishTime).
		Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "TaskForExpireInvoices error:%s", err.Error())
		return
	}
	for _, one := range list {
		key := fmt.Sprintf("TaskForExpireInvoices-%v", one.Id)
		if utility.TryLock(ctx, key, 60) {
			g.Log().Debugf(ctx, "TaskForExpireInvoices GetLock 60s, key:%s", key)
			defer func() {
				utility.ReleaseLock(ctx, key)
				g.Log().Debugf(ctx, "TaskForExpireInvoices ReleaseLock, key:%s", key)
			}()
			if one.FinishTime == 0 {
				_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
					dao.Invoice.Columns().FinishTime: one.GmtModify.Timestamp(),
					dao.Invoice.Columns().GmtModify:  gtime.Now(),
				}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
				if err != nil {
					g.Log().Errorf(ctx, "TaskForExpireInvoices Update FinishTime error:%s", err.Error())
				}
			} else if one.FinishTime+(one.DayUtilDue*86400)+1200 < gtime.Now().Timestamp() { // task delay 20 minutes to expire
				//Invoice Expire
				err = service.ProcessingInvoiceFailure(ctx, one.InvoiceId, "TaskForExpireInvoices")
				if err != nil {
					g.Log().Errorf(ctx, "TaskForExpireInvoices Failure invoice error:%s", err.Error())
				}
			}
		} else {
			g.Log().Errorf(ctx, "TaskForExpireInvoices GetLock Failure, key:%s", key)
			return
		}
	}
}

func ExpireUserSubInvoices(ctx context.Context, sub *entity.Subscription, timeNow int64) {
	if sub == nil {
		return
	}
	var list []*entity.Invoice
	err := dao.Invoice.Ctx(ctx).
		Where(dao.Invoice.Columns().UserId, sub.UserId).
		Where(dao.Invoice.Columns().SubscriptionId, sub.SubscriptionId).
		Where(dao.Invoice.Columns().BizType, consts.BizTypeSubscription).
		Where(dao.Invoice.Columns().Status, consts.InvoiceStatusProcessing).
		Where(dao.Invoice.Columns().IsDeleted, 0).
		OrderAsc(dao.Invoice.Columns().FinishTime).
		Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "ExpireUserSubInvoices error:%s", err.Error())
		return
	}
	for _, one := range list {
		key := fmt.Sprintf("TaskForExpireInvoices-%v", one.Id)
		if utility.TryLock(ctx, key, 60) {
			g.Log().Debugf(ctx, "ExpireUserSubInvoices GetLock 60s, key:%s", key)
			defer func() {
				utility.ReleaseLock(ctx, key)
				g.Log().Debugf(ctx, "ExpireUserSubInvoices ReleaseLock, key:%s", key)
			}()
			if one.FinishTime == 0 {
				_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
					dao.Invoice.Columns().FinishTime: one.GmtModify.Timestamp(),
					dao.Invoice.Columns().GmtModify:  gtime.Now(),
				}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
				if err != nil {
					g.Log().Errorf(ctx, "ExpireUserSubInvoices Update FinishTime error:%s", err.Error())
				}
			} else if one.FinishTime+(one.DayUtilDue*86400)+600 < timeNow {
				//Invoice Expire
				err = service.ProcessingInvoiceFailure(ctx, one.InvoiceId, "ExpireUserSubInvoices")
				if err != nil {
					g.Log().Errorf(ctx, "ExpireUserSubInvoices Failure invoice error:%s", err.Error())
				}
			}
		} else {
			g.Log().Errorf(ctx, "ExpireUserSubInvoices GetLock Failure, key:%s", key)
			return
		}
	}
}

func TaskForCompensateSubUpDownInvoices(ctx context.Context) {
	var count = 100
	var page = 0
	for {
		var list []*entity.SubscriptionPendingUpdate
		err := dao.SubscriptionPendingUpdate.Ctx(ctx).
			WhereGT(dao.SubscriptionPendingUpdate.Columns().Status, 0).
			Where(dao.SubscriptionPendingUpdate.Columns().IsDeleted, 0).
			WhereNotNull(dao.SubscriptionPendingUpdate.Columns().InvoiceId).
			Limit(page*count, count).
			Scan(&list)
		if err != nil {
			g.Log().Errorf(ctx, "TaskForCompensateSubUpDownInvoices Get List error:%s", err.Error())
			return
		}
		for _, one := range list {
			time.Sleep(1 * time.Second)
			invoiceDetail := detail.InvoiceDetail(ctx, one.InvoiceId)
			if invoiceDetail == nil {
				g.Log().Infof(ctx, "TaskForCompensateSubUpDownInvoices invoice not found, ignore")
				continue
			}
			if invoiceDetail.Metadata != nil {
				if _, ok := invoiceDetail.Metadata["IsUpgrade"]; ok {
					if _, ok2 := invoiceDetail.Metadata["SubscriptionUpdate"]; ok2 {
						g.Log().Infof(ctx, "TaskForCompensateSubUpDownInvoices not old version, ignore invoiceId:%s ", invoiceDetail.InvoiceId)
						continue
					}
				}
			}
			metaData := invoiceDetail.Metadata
			if metaData == nil {
				metaData = make(map[string]interface{})
			}
			var isUpgrade bool
			var changed bool
			//check isUpgrade or not
			oldPlan := query.GetPlanById(ctx, one.PlanId)
			plan := query.GetPlanById(ctx, one.UpdatePlanId)
			if plan.IntervalUnit != oldPlan.IntervalUnit || plan.IntervalCount != oldPlan.IntervalCount {
				isUpgrade = true
				changed = true
			} else if plan.Amount > oldPlan.Amount || plan.Amount*one.UpdateQuantity > oldPlan.Amount*one.Quantity {
				isUpgrade = true
				changed = true
			} else if plan.Amount < oldPlan.Amount || plan.Amount*one.UpdateQuantity < oldPlan.Amount*one.Quantity {
				isUpgrade = false
				changed = true
			} else {
				var oldAddonParams []*bean.PlanAddonParam
				err := utility.UnmarshalFromJsonString(one.AddonData, &oldAddonParams)
				utility.Assert(err == nil, fmt.Sprintf("UnmarshalFromJsonString internal err:%v", err))
				var oldAddonMap = make(map[uint64]int64)
				for _, oldAddon := range oldAddonParams {
					if _, ok := oldAddonMap[oldAddon.AddonPlanId]; ok {
						oldAddonMap[oldAddon.AddonPlanId] = oldAddonMap[oldAddon.AddonPlanId] + oldAddon.Quantity
					} else {
						oldAddonMap[oldAddon.AddonPlanId] = oldAddon.Quantity
					}
				}
				var addonParams []*bean.PlanAddonParam
				err = utility.UnmarshalFromJsonString(one.UpdateAddonData, &addonParams)
				utility.Assert(err == nil, fmt.Sprintf("UnmarshalFromJsonString internal err:%v", err))
				var newAddonMap = make(map[uint64]int64)
				for _, newAddon := range addonParams {
					if _, ok := newAddonMap[newAddon.AddonPlanId]; ok {
						newAddonMap[newAddon.AddonPlanId] = newAddonMap[newAddon.AddonPlanId] + newAddon.Quantity
					} else {
						newAddonMap[newAddon.AddonPlanId] = newAddon.Quantity
					}
				}
				for newAddonPlanId, newAddonQuantity := range newAddonMap {
					if oldAddonQuantity, ok := oldAddonMap[newAddonPlanId]; ok {
						if oldAddonQuantity < newAddonQuantity {
							isUpgrade = true
							changed = true
							break
						}
					} else {
						isUpgrade = true
						changed = true
						break
					}
				}
				if len(oldAddonMap) != len(newAddonMap) {
					changed = true
				} else {
					for newAddonPlanId, newAddonQuantity := range newAddonMap {
						if oldAddonQuantity, ok := oldAddonMap[newAddonPlanId]; ok {
							if oldAddonQuantity != newAddonQuantity {
								changed = true
								break
							}
						} else {
							changed = true
							break
						}
					}
				}
			}

			if !changed {
				g.Log().Infof(ctx, "TaskForCompensateSubUpDownInvoices nochange ignore invoiceId:%s ", invoiceDetail.InvoiceId)
				continue
			}
			g.Log().Infof(ctx, "TaskForCompensateSubUpDownInvoices invoiceId:%s oldMetadata:%s", invoiceDetail.InvoiceId, utility.MarshalToJsonString(metaData))
			metaData["IsUpgrade"] = isUpgrade
			metaData["SubscriptionUpdate"] = true
			metaData["FromCompensate"] = true
			_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
				dao.Invoice.Columns().MetaData:  utility.MarshalToJsonString(metaData),
				dao.Invoice.Columns().GmtModify: gtime.Now(),
			}).Where(dao.Invoice.Columns().InvoiceId, invoiceDetail.InvoiceId).OmitNil().Update()
			if err != nil {
				g.Log().Errorf(ctx, "TaskForCompensateSubUpDownInvoices Update Invoice newMetadata error:%s", err.Error())
			} else {
				g.Log().Infof(ctx, "TaskForCompensateSubUpDownInvoices invoiceId:%s Update newMetadata:%s", invoiceDetail.InvoiceId, utility.MarshalToJsonString(metaData))
			}
		}
		// next page
		page = page + 1
		if list == nil || len(list) == 0 {
			break
		}
	}
}
