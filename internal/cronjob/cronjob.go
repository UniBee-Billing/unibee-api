package cronjob

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"go-oversea-pay/internal/cronjob/sub"
)

func StartCronJobs() {
	var (
		err error
		ctx = gctx.New()
	)
	//if consts.GetConfigInstance().IsServerDev() {
	//	return
	//}
	// every 10 second
	var name = "SubscriptionBillingCycleDunningInvoice"
	g.Log().Print(ctx, "CronJob Start......")
	_, err = gcron.AddSingleton(ctx, "*/10 * * * * *", func(ctx context.Context) {
		sub.SubscriptionBillingCycleDunningInvoice(ctx, name)
	}, name)
	if err != nil {
		g.Log().Print(ctx, "StartCronJobs Name:%s Err:%s", name, err.Error())
	}
	// every 10 second
	var backName = "SubscriptionBillingCycleDunningInvoiceBackup"
	_, err = gcron.AddSingleton(ctx, "*/10 * * * * *", func(ctx context.Context) {
		sub.SubscriptionBillingCycleDunningInvoice(ctx, backName)
	}, backName)
	if err != nil {
		g.Log().Print(ctx, "StartCronJobs Name:%s Err:%s", backName, err.Error())
	}
	// every hour
	var httpLogDeleteTaskName = "httpLogDeleteTaskName"
	_, err = gcron.AddSingleton(ctx, "* * */1 * * *", func(ctx context.Context) {

	}, backName)
	if err != nil {
		g.Log().Print(ctx, "StartCronJobs Name:%s Err:%s", httpLogDeleteTaskName, err.Error())
	}
	return
}
