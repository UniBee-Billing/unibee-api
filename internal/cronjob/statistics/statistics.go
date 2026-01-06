package statistics

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/analysis/statistics"
)

func TaskForUpdateAllMerchantStatistics(ctx context.Context) {
	g.Log().Infof(ctx, "TaskForUpdateAllMerchantStatistics start")
	statistics.UpdateMerchantStatsCron(ctx)
	g.Log().Infof(ctx, "TaskForUpdateAllMerchantStatistics end")
}
