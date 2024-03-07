package query

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

func GetPlanById(ctx context.Context, id uint64) (one *entity.Plan) {
	if id <= 0 {
		return nil
	}
	err := dao.Plan.Ctx(ctx).Where(dao.Plan.Columns().Id, id).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetPlanBindingAddonsByPlanId(ctx context.Context, id uint64) (list []*entity.Plan) {
	if id <= 0 {
		return nil
	}
	var one *entity.Plan
	err := dao.Plan.Ctx(ctx).Where(dao.Plan.Columns().Id, id).OmitEmpty().Scan(&one)
	if err != nil && one == nil {
		return nil
	}
	if len(one.BindingAddonIds) == 0 {
		return nil
	}
	var addonIdsList []uint64
	if len(one.BindingAddonIds) > 0 {
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
	err = dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, addonIdsList).Scan(&list)
	if err != nil {
		return nil
	}
	return list
}
