package cronjob

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"unibee/internal/cronjob/discount"
	"unibee/internal/cronjob/gateway_log"
	"unibee/internal/cronjob/invoice"
	"unibee/internal/cronjob/sub"
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
	var name = "SubscriptionCycle"
	g.Log().Print(ctx, "CronJob Start......")
	_, err = gcron.AddSingleton(ctx, "@every 10s", func(ctx context.Context) {
		sub.TaskForSubscriptionBillingCycleDunningInvoice(ctx, name)
	}, name)
	if err != nil {
		g.Log().Printf(ctx, "StartCronJobs Name:%s Err:%s\n", name, err.Error())
	}
	// every 10 second
	var backName = "SubscriptionCycleBackup"
	_, err = gcron.AddSingleton(ctx, "@every 10s", func(ctx context.Context) {
		sub.TaskForSubscriptionBillingCycleDunningInvoice(ctx, backName)
	}, backName)
	if err != nil {
		g.Log().Printf(ctx, "StartCronJobs Name:%s Err:%s\n", backName, err.Error())
	}
	// every 1 min
	var other1MinTask = "Other1MinTask"
	_, err = gcron.AddSingleton(ctx, "@every 1m", func(ctx context.Context) {
		discount.TaskForExpireDiscounts(ctx)
		invoice.TaskForExpireInvoices(ctx)
	}, other1MinTask)
	if err != nil {
		g.Log().Printf(ctx, "StartCronJobs Name:%s Err:%s\n", other1MinTask, err.Error())
	}

	// every 10 min
	var other10MinTask = "Other10MinTask"
	_, err = gcron.AddSingleton(ctx, "@every 10m", func(ctx context.Context) {
		sub.TaskForSubscriptionTrackAfterCancelledOrExpired(ctx, other10MinTask)
	}, other10MinTask)
	if err != nil {
		g.Log().Printf(ctx, "StartCronJobs Name:%s Err:%s\n", other10MinTask, err.Error())
	}

	// every Hour
	var hourTask = "OtherTask"
	_, err = gcron.AddSingleton(ctx, "@hourly", func(ctx context.Context) {
		gateway_log.TaskForDeleteChannelLogs(ctx)
	}, hourTask)
	if err != nil {
		g.Log().Printf(ctx, "StartCronJobs Name:%s Err:%s\n", hourTask, err.Error())
	}
	return
}
