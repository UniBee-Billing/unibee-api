package consts

type SubscriptionPlanStatusEnum int

const (
	CREATE   = 0 //创建
	ACTIVE   = 1 //有效
	INACTIVE = 2 //无效
)

func (status SubscriptionPlanStatusEnum) Description() string {
	switch status {
	case CREATE:
		return "CREATE"
	case ACTIVE:
		return "ACTIVE"
	case INACTIVE:
		return "INACTIVE"
	default:
		return "ACTIVE"
	}
}
