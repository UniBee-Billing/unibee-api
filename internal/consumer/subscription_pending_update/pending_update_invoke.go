package subscription

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	"unibee/internal/query"
	"unibee/utility"
)

func init() {
	redismq.RegisterInvoke("GetPendingUpdateByInvoiceId", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		g.Log().Infof(ctx, "GetPendingUpdateByInvoiceId:%s", utility.MarshalToJsonString(request))
		if invoiceId, ok := request.(string); ok {
			if len(invoiceId) > 0 {
				pendingUpdate := query.GetSubscriptionPendingUpdateByInvoiceId(ctx, invoiceId)
				if pendingUpdate != nil {
					g.Log().Infof(ctx, "GetPendingUpdateByInvoiceId:%s get pendingUpdate:%s", utility.MarshalToJsonString(request), utility.MarshalToJsonString(pendingUpdate))
					return utility.MarshalToJsonString(pendingUpdate), nil
				} else {
					g.Log().Infof(ctx, "GetPendingUpdateByInvoiceId:%s get pendingUpdate not found", utility.MarshalToJsonString(request))
				}
			} else {
				g.Log().Infof(ctx, "GetPendingUpdateByInvoiceId:%s invoiceId length should greater than 0", utility.MarshalToJsonString(request))
			}
		} else {
			g.Log().Infof(ctx, "GetPendingUpdateByInvoiceId:%s request not string", utility.MarshalToJsonString(request))
		}
		return nil, nil
	})
}
