package consts

type SubscriptionPlanStatusEnum int

const (
	PlanStatusCreate   = 0 //创建
	PlanStatusActive   = 1 //有效
	PlanStatusInActive = 2 //无效
)

func (status SubscriptionPlanStatusEnum) Description() string {
	switch status {
	case PlanStatusCreate:
		return "STATUES_CREATE"
	case PlanStatusActive:
		return "PlanStatusActive"
	case PlanStatusInActive:
		return "PlanStatusInActive"
	default:
		return "PlanStatusActive"
	}
}
