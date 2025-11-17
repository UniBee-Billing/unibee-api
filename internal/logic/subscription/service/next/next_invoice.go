package next

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	"unibee/internal/query"
	"unibee/utility"
)

func GetSubscriptionNextInvoiceData(ctx context.Context, subId string) *bean.SubscriptionNextInvoiceData {
	if subId == "" {
		return nil
	}
	sub := query.GetSubscriptionBySubscriptionId(ctx, subId)
	if sub == nil {
		return nil
	}
	var data *bean.SubscriptionNextInvoiceData
	if len(sub.NextInvoiceData) > 0 {
		_ = utility.UnmarshalFromJsonString(sub.NextInvoiceData, &data)
	}
	return data
}

func SaveSubscriptionNextInvoiceData(ctx context.Context, subId string, data *bean.SubscriptionNextInvoiceData) {
	if len(subId) == 0 || data == nil {
		return
	}
	_, _ = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().NextInvoiceData: utility.MarshalToJsonString(data),
		dao.Subscription.Columns().GmtModify:       gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subId).Update()
}

func ApplySubscriptionNextInvoiceData(ctx context.Context, subId string, invoiceId string) {
	if len(subId) == 0 || len(invoiceId) == 0 {
		return
	}
	data := GetSubscriptionNextInvoiceData(ctx, subId)
	if data != nil {
		data.ApplyInvoiceId = invoiceId
		_, _ = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().NextInvoiceData: utility.MarshalToJsonString(data),
			dao.Subscription.Columns().GmtModify:       gtime.Now(),
		}).Where(dao.Subscription.Columns().SubscriptionId, subId).Update()
	}
}

func ClearSubscriptionNextInvoiceData(ctx context.Context, subId string, invoiceId string) {
	if len(subId) == 0 || len(invoiceId) == 0 {
		return
	}
	data := GetSubscriptionNextInvoiceData(ctx, subId)
	if data != nil && data.ApplyInvoiceId == invoiceId {
		_, _ = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().NextInvoiceData: "",
			dao.Subscription.Columns().GmtModify:       gtime.Now(),
		}).Where(dao.Subscription.Columns().SubscriptionId, subId).Update()
	}
}
