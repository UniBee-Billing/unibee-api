package invoice

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/invoice/service"
	entity "unibee/internal/model/entity/oversea_pay"
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
			g.Log().Print(ctx, "TaskForExpireInvoices GetLock 60s", key)
			defer func() {
				utility.ReleaseLock(ctx, key)
				g.Log().Print(ctx, "TaskForExpireInvoices ReleaseLock", key)
			}()
			if one.FinishTime == 0 {
				_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
					dao.Invoice.Columns().FinishTime: one.GmtModify.Timestamp(),
					dao.Invoice.Columns().GmtModify:  gtime.Now(),
				}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
				if err != nil {
					g.Log().Errorf(ctx, "TaskForExpireInvoices Update FinishTime error:", err.Error())
				}
			} else if one.FinishTime+(one.DayUtilDue*86400) < gtime.Now().Timestamp() {
				//Invoice Expire
				err = service.ProcessingInvoiceFailure(ctx, one.InvoiceId)
				if err != nil {
					g.Log().Errorf(ctx, "TaskForExpireInvoices Failure invoice error:", err.Error())
				}
			}
		} else {
			g.Log().Errorf(ctx, "TaskForExpireInvoices GetLock Failure", key)
			return
		}
	}
}

func ExpireUserOneTimeInvoices(ctx context.Context, sub *entity.Subscription) {
	if sub == nil {
		return
	}
	var timeNow = gtime.Now().Timestamp()
	var list []*entity.Invoice
	err := dao.Invoice.Ctx(ctx).
		Where(dao.Invoice.Columns().UserId, sub.UserId).
		Where(dao.Invoice.Columns().BizType, consts.BizTypeOneTime).
		Where(dao.Invoice.Columns().Status, consts.InvoiceStatusProcessing).
		Where(dao.Invoice.Columns().IsDeleted, 0).
		OrderAsc(dao.Invoice.Columns().FinishTime).
		Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "ExpireUserOneTimeInvoices error:%s", err.Error())
		return
	}
	for _, one := range list {
		key := fmt.Sprintf("TaskForExpireInvoices-%v", one.Id)
		if utility.TryLock(ctx, key, 60) {
			g.Log().Print(ctx, "ExpireUserOneTimeInvoices GetLock 60s", key)
			defer func() {
				utility.ReleaseLock(ctx, key)
				g.Log().Print(ctx, "ExpireUserOneTimeInvoices ReleaseLock", key)
			}()
			if one.FinishTime == 0 {
				_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
					dao.Invoice.Columns().FinishTime: one.GmtModify.Timestamp(),
					dao.Invoice.Columns().GmtModify:  gtime.Now(),
				}).Where(dao.Invoice.Columns().Id, one.Id).OmitNil().Update()
				if err != nil {
					g.Log().Errorf(ctx, "ExpireUserOneTimeInvoices Update FinishTime error:", err.Error())
				}
			} else if one.FinishTime+(one.DayUtilDue*86400) < timeNow {
				//Invoice Expire
				err = service.ProcessingInvoiceFailure(ctx, one.InvoiceId)
				if err != nil {
					g.Log().Errorf(ctx, "ExpireUserOneTimeInvoices Failure invoice error:", err.Error())
				}
			}
		} else {
			g.Log().Errorf(ctx, "ExpireUserOneTimeInvoices GetLock Failure", key)
			return
		}
	}
}
