package batch

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func TaskForExpireBatchTasks(ctx context.Context) {
	var list []*entity.MerchantBatchTask
	err := dao.MerchantBatchTask.Ctx(ctx).
		WhereLT(dao.MerchantBatchTask.Columns().CreateTime, gtime.Now().Timestamp()-6*3500).
		WhereLT(dao.MerchantBatchTask.Columns().Status, 2).
		Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "TaskForExpireBatchTasks error:%s", err.Error())
		return
	}
	for _, one := range list {
		if len(one.FailReason) == 0 {
			_, _ = dao.MerchantBatchTask.Ctx(ctx).Data(g.Map{
				dao.MerchantBatchTask.Columns().Status:         3,
				dao.MerchantBatchTask.Columns().FailReason:     "Expired After 6 hours",
				dao.MerchantBatchTask.Columns().LastUpdateTime: gtime.Now().Timestamp(),
				dao.MerchantBatchTask.Columns().GmtModify:      gtime.Now(),
			}).Where(dao.MerchantBatchTask.Columns().Id, one.Id).OmitNil().Update()
		}
	}
}
