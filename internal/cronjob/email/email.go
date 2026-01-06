package email

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func TaskForCompensateEmailHistory(ctx context.Context) {
	var list []*entity.MerchantEmailHistory
	err := dao.MerchantEmailHistory.Ctx(ctx).
		Where(dao.MerchantEmailHistory.Columns().Status, 0).
		Limit(0, 1000).
		Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "TaskForCompensateEmailHistory error:%s", err.Error())
		return
	}
	for _, one := range list {
		status := one.Status
		if strings.Contains(one.Response, "202") {
			status = 1
		} else {
			status = 2
		}
		_, err = dao.MerchantEmailHistory.Ctx(ctx).Data(g.Map{
			dao.MerchantEmailHistory.Columns().Status: status,
		}).Where(dao.MerchantEmailHistory.Columns().Id, one.Id).OmitNil().Update()
	}
}
