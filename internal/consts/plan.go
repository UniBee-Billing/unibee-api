package consts

type PlanType int

const (
	PlanTypeMain           = 1
	PlanTypeRecurringAddon = 2
	PlanTypeOnetimeAddon   = 3
)

type PlanStatusEnum int

const (
	PlanStatusEditable           = 1
	PlanStatusActive             = 2
	PlanStatusExpired            = 3
	PlanPublishStatusPublished   = 2
	PlanPublishStatusUnPublished = 1
)

func (status PlanStatusEnum) Description() string {
	switch status {
	case PlanStatusEditable:
		return "PlanStatusEditable"
	case PlanStatusActive:
		return "PlanStatusActive"
	case PlanStatusExpired:
		return "PlanStatusExpired"
	default:
		return "PlanStatusEditable"
	}
}

type GatewayPlanStatusEnum int

const (
	GatewayPlanStatusInit     = 0
	GatewayPlanStatusCreate   = 1
	GatewayPlanStatusActive   = 2
	GatewayPlanStatusInActive = 3
)

func (status GatewayPlanStatusEnum) Description() string {
	switch status {
	case GatewayPlanStatusInit:
		return "GatewayPlanStatusInit"
	case GatewayPlanStatusCreate:
		return "GatewayPlanStatusCreate"
	case GatewayPlanStatusActive:
		return "GatewayPlanStatusActive"
	case GatewayPlanStatusInActive:
		return "GatewayPlanStatusInActive"
	default:
		return "GatewayPlanStatusInit"
	}
}