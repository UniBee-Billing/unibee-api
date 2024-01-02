package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/subscription/service"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strconv"
	"strings"

	"go-oversea-pay/api/subscription/v1"
)

func (c *ControllerV1) SubscriptionPlanChannelTransferAndActivate(ctx context.Context, req *v1.SubscriptionPlanChannelTransferAndActivateReq) (res *v1.SubscriptionPlanChannelTransferAndActivateRes, err error) {
	utility.Assert(req.PlanId > 0, "plan should > 0")
	//utility.Assert(req.ChannelId > 0, "ConfirmChannelId should > 0")
	plan := query.GetSubscriptionPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "plan not found")
	//多个渠道Plan 创建并激活
	list := query.GetListSubscriptionTypePayChannels(ctx) // todo mark 需改造成获取 merchantId 相关的 Channel
	utility.Assert(len(list) > 0, "no channel found, need at least one")
	for _, channel := range list {
		err = service.SubscriptionPlanChannelTransferAndActivate(ctx, req.PlanId, int64(channel.Id))
		if err != nil {
			utility.FailureJsonExit(g.RequestFromCtx(ctx), fmt.Sprintf("%s", err))
			return
		}
	}

	//发布 Plan
	err = service.SubscriptionPlanActivate(ctx, req.PlanId)
	if err != nil {
		utility.FailureJsonExit(g.RequestFromCtx(ctx), fmt.Sprintf("%s", err))
		return
	}

	if len(plan.BindingAddonIds) > 0 {
		//addon 检查
		var addonIds []int64
		var addonIdsList []int64
		if len(plan.BindingAddonIds) > 0 {
			//初始化
			strList := strings.Split(plan.BindingAddonIds, ",")
			for _, s := range strList {
				num, err := strconv.ParseInt(s, 10, 64) // 将字符串转换为整数
				if err != nil {
					fmt.Println("Internal Error converting string to int:", err)
					return nil, err
				}
				addonIdsList = append(addonIdsList, num) // 添加到整数列表中
				addonIds = append(addonIds, num)         // 添加到整数列表中
			}
		}
		//检查 addonIds 类型
		var allAddonList []*entity.SubscriptionPlan
		err = dao.SubscriptionPlan.Ctx(ctx).WhereIn(dao.SubscriptionPlan.Columns().Id, addonIds).Scan(&allAddonList)
		for _, addonPlan := range allAddonList {
			if addonPlan.Status != consts.PlanStatusPublished {
				//发布 addonPlan
				for _, channel := range list {
					err = service.SubscriptionPlanChannelTransferAndActivate(ctx, int64(addonPlan.Id), int64(channel.Id))
					if err != nil {
						utility.FailureJsonExit(g.RequestFromCtx(ctx), fmt.Sprintf("%s", err))
						return
					}
				}

				//发布 Plan
				err = service.SubscriptionPlanActivate(ctx, int64(addonPlan.Id))
				if err != nil {
					utility.FailureJsonExit(g.RequestFromCtx(ctx), fmt.Sprintf("%s", err))
					return
				}
			}
		}
	}

	utility.SuccessJsonExit(g.RequestFromCtx(ctx), nil)
	return
}
