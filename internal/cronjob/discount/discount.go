package discount

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/discount"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

func TaskForExpireDiscounts(ctx context.Context) {
	var list []*entity.MerchantDiscountCode
	err := dao.MerchantDiscountCode.Ctx(ctx).
		Where(dao.MerchantDiscountCode.Columns().Status, discount.DiscountStatusDeActive).
		WhereLT(dao.MerchantDiscountCode.Columns().EndTime, gtime.Now().Timestamp()).
		Where(dao.MerchantDiscountCode.Columns().IsDeleted, 0).
		Scan(&list)
	if err != nil {
		g.Log().Errorf(ctx, "TaskForExpireDiscounts error:%s", err.Error())
		return
	}
	for _, one := range list {
		key := fmt.Sprintf("TaskForExpireDiscounts-%s", one.Id)
		if utility.TryLock(ctx, key, 60) {
			g.Log().Print(ctx, "TaskForExpireDiscounts GetLock 60s", key)
			defer func() {
				utility.ReleaseLock(ctx, key)
				g.Log().Print(ctx, "TaskForExpireDiscounts ReleaseLock", key)
			}()
			_, err := dao.MerchantDiscountCode.Ctx(ctx).Data(g.Map{
				dao.MerchantDiscountCode.Columns().Status:    discount.DiscountStatusExpired,
				dao.MerchantDiscountCode.Columns().GmtModify: gtime.Now(),
			}).Where(dao.MerchantDiscountCode.Columns().Id, one.Id).OmitNil().Update()
			if err != nil {
				g.Log().Errorf(ctx, "TaskForExpireDiscounts Expire id:%d error:%s", one.Id, err.Error())
			}
		} else {
			g.Log().Print(ctx, "TaskForExpireDiscounts GetLock Failure", key)
			return
		}
	}
}