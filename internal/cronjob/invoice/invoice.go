package invoice

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

func TaskForExpireInvoices(ctx context.Context) {
	var list []*entity.Invoice
	err := dao.Invoice.Ctx(ctx).
		Where(dao.Invoice.Columns().Status, consts.InvoiceStatusProcessing).
		WhereLT(dao.Invoice.Columns().FinishTime, gtime.Now().Timestamp()).
		Where(dao.Invoice.Columns().IsDeleted, 0).
		Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "TaskForExpireInvoices error:%s", err.Error())
		return
	}
	for _, one := range list {
		key := fmt.Sprintf("TaskForExpireInvoices-%s", one.Id)
		if utility.TryLock(ctx, key, 60) {
			g.Log().Print(ctx, "TaskForExpireInvoices GetLock 60s", key)
			defer func() {
				utility.ReleaseLock(ctx, key)
				g.Log().Print(ctx, "TaskForExpireInvoices ReleaseLock", key)
			}()
			// todo mark expire invoice
		} else {
			g.Log().Print(ctx, "TaskForExpireInvoices GetLock Failure", key)
			return
		}
	}
}