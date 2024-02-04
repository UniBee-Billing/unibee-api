package query

import (
	"context"
	"fmt"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"strconv"
	"strings"
)

func GetPlanById(ctx context.Context, id int64) (one *entity.SubscriptionPlan) {
	if id <= 0 {
		return nil
	}
	err := dao.SubscriptionPlan.Ctx(ctx).Where(entity.SubscriptionPlan{Id: uint64(id)}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetPlanBindingAddonsByPlanId(ctx context.Context, id int64) (list []*entity.SubscriptionPlan) {
	if id <= 0 {
		return nil
	}
	var one *entity.SubscriptionPlan
	err := dao.SubscriptionPlan.Ctx(ctx).Where(entity.SubscriptionPlan{Id: uint64(id)}).OmitEmpty().Scan(&one)
	if err != nil && one == nil {
		return nil
	}
	if len(one.BindingAddonIds) == 0 {
		return nil
	}
	var addonIdsList []uint64
	if len(one.BindingAddonIds) > 0 {
		//初始化
		strList := strings.Split(one.BindingAddonIds, ",")

		for _, s := range strList {
			num, err := strconv.ParseInt(s, 10, 64) // 将字符串转换为整数
			if err != nil {
				fmt.Println("Internal Error converting string to int:", err)
				return nil
			}
			addonIdsList = append(addonIdsList, uint64(num)) // 添加到整数列表中
		}
	}
	err = dao.SubscriptionPlan.Ctx(ctx).WhereIn(dao.SubscriptionPlan.Columns().Id, addonIdsList).Scan(&list)
	if err != nil {
		return nil
	}
	return list
}
