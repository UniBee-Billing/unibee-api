package cronjob

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"unibee/internal/cmd/config"
	"unibee/internal/cronjob/batch"
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
	// every 10 second
	var name = "JobSubscriptionCycle"
	g.Log().Infof(ctx, "CronJob Start......")
	_, err = gcron.AddSingleton(ctx, "@every 10s", func(ctx context.Context) {
		sub.TaskForSubscriptionBillingCycleDunningInvoice(ctx, name)
	}, name)
	if err != nil {
		g.Log().Errorf(ctx, "StartCronJobs Name:%s Err:%s\n", name, err.Error())
	}
	// every 10 second
	var backName = "JobSubscriptionCycleBackup"
	_, err = gcron.AddSingleton(ctx, "@every 10s", func(ctx context.Context) {
		sub.TaskForSubscriptionBillingCycleDunningInvoice(ctx, backName)
	}, backName)
	if err != nil {
		g.Log().Errorf(ctx, "StartCronJobs Name:%s Err:%s\n", backName, err.Error())
	}
	// every 1 min
	var other1MinTask = "Job1MinTask"
	_, err = gcron.Add(ctx, "@every 1m", func(ctx context.Context) {
		discount.TaskForExpireDiscounts(ctx)
		invoice.TaskForExpireInvoices(ctx)
		batch.TaskForExpireBatchTasks(ctx)
	}, other1MinTask)
	if err != nil {
		g.Log().Errorf(ctx, "StartCronJobs Name:%s Err:%s\n", other1MinTask, err.Error())
	}

	// every 10 min
	var other10MinTask = "Job10MinTask"
	_, err = gcron.Add(ctx, "@every 10m", func(ctx context.Context) {
		sub.TaskForSubscriptionTrackAfterCancelledOrExpired(ctx, other10MinTask)
		sub.TaskForSubscriptionInitFailed(ctx, other10MinTask)
	}, other10MinTask)
	if err != nil {
		g.Log().Errorf(ctx, "StartCronJobs Name:%s Err:%s\n", other10MinTask, err.Error())
	}

	// every Hour
	var hourTask = "JobHourlyTask"
	_, err = gcron.Add(ctx, "@hourly", func(ctx context.Context) {
		gateway_log.TaskForDeleteChannelLogs(ctx)
		gateway_log.TaskForDeleteWebhookMessage(ctx)
		sub.TaskForUserSubCompensate(ctx, hourTask)
		if !config.GetConfigInstance().IsProd() {
			invoice.TaskForCompensateSubUpDownInvoices(ctx)
		}
	}, hourTask)
	if err != nil {
		g.Log().Errorf(ctx, "StartCronJobs Name:%s Err:%s\n", hourTask, err.Error())
	}

	// every day
	var dailyTask = "JobDailyTask"
	_, err = gcron.Add(ctx, "@daily", func(ctx context.Context) {
		invoice.TaskForCompensateSubUpDownInvoices(ctx)
	}, dailyTask)
	if err != nil {
		g.Log().Errorf(ctx, "StartCronJobs Name:%s Err:%s\n", dailyTask, err.Error())
	}

	return
}
