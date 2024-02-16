package cronjob

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"unibee-api/internal/cronjob/sub"
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
		g.Log().Printf(ctx, "StartCronJobs Name:%s Err:%s\n", name, err.Error())
	}
	// every 10 second
	var backName = "SubscriptionBillingCycleDunningInvoiceBackup"
	_, err = gcron.AddSingleton(ctx, "*/10 * * * * *", func(ctx context.Context) {
		sub.SubscriptionBillingCycleDunningInvoice(ctx, backName)
	}, backName)
	if err != nil {
		g.Log().Printf(ctx, "StartCronJobs Name:%s Err:%s\n", backName, err.Error())
	}
	// every hour
	var httpLogDeleteTaskName = "httpLogDeleteTaskName"
	_, err = gcron.AddSingleton(ctx, "* * */1 * * *", func(ctx context.Context) {

	}, httpLogDeleteTaskName)
	if err != nil {
		g.Log().Printf(ctx, "StartCronJobs Name:%s Err:%s\n", httpLogDeleteTaskName, err.Error())
	}
	return
}
