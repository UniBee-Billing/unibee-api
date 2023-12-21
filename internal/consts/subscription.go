package consts

type SubscriptionStatusEnum int

const (
	SubStatusInit     = 0 //初始化
	SubStatusCreate   = 1 //创建-渠道状态
	SubStatusActive   = 2 //有效-渠道状态
	SubStatusInActive = 3 //无效-渠道状态
)

func (status SubscriptionStatusEnum) Description() string {
	switch status {
	case SubStatusInit:
		return "SubStatusInit"
	case SubStatusCreate:
		return "SubStatusCreate"
	case SubStatusActive:
		return "SubStatusActive"
	case SubStatusInActive:
		return "SubStatusInActive"
	default:
		return "SubStatusInit"
	}
}

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
