package consts

type SubscriptionPlanStatusEnum int

const (
	PlanStatusInit     = 0 //初始化
	PlanStatusCreate   = 1 //创建-渠道状态
	PlanStatusActive   = 2 //有效-渠道状态
	PlanStatusInActive = 3 //无效-渠道状态
)

func (status SubscriptionPlanStatusEnum) Description() string {
	switch status {
	case PlanStatusInit:
		return "STATUES_CREATE"
	case PlanStatusCreate:
		return "STATUES_CREATE"
	case PlanStatusActive:
		return "PlanStatusActive"
	case PlanStatusInActive:
		return "PlanStatusInActive"
	default:
		return "PlanStatusInit"
	}
}
