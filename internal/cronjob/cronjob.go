package cronjob

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/cronjob/sub"
)

func StartCronJobs() {
	var (
		err error
		ctx = gctx.New()
	)
	if !consts.GetConfigInstance().IsServerDev() {
		return
	}
	var name = "SubscriptionBillingCycleDunningInvoice"
	g.Log().Print(ctx, "CronJob Start......")
	_, err = gcron.AddSingleton(ctx, "*/10 * * * * *", func(ctx context.Context) {
		sub.SubscriptionBillingCycleDunningInvoice(ctx, name)
	}, name)
	if err != nil {
		g.Log().Print(ctx, "StartCronJobs Name:%s Err:%s", name, err.Error())
	}
	var backName = "SubscriptionBillingCycleDunningInvoiceBackup"
	_, err = gcron.AddSingleton(ctx, "*/10 * * * * *", func(ctx context.Context) {
		sub.SubscriptionBillingCycleDunningInvoice(ctx, backName)
	}, backName)
	if err != nil {
		g.Log().Print(ctx, "StartCronJobs Name:%s Err:%s", backName, err.Error())
	}
	return
}
