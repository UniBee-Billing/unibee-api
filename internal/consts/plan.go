package consts

type PlanType int

const (
	PlanTypeMain           = 1
	PlanTypeRecurringAddon = 2
	PlanTypeOnetime        = 3
)

type PlanStatusEnum int

const (
	PlanStatusEditable           = 1
	PlanStatusActive             = 2
	PlanStatusInActive           = 3
	PlanStatusSoftArchive        = 4
	PlanStatusHardArchive        = 5
	PlanPublishStatusPublished   = 2
	PlanPublishStatusUnPublished = 1
)

func (status PlanStatusEnum) Description() string {
	switch status {
	case PlanStatusEditable:
		return "Editable"
	case PlanStatusActive:
		return "Active"
	case PlanStatusSoftArchive:
		return "SoftArchive"
	case PlanStatusHardArchive:
		return "HardArchive"
	default:
		return "Editable"
	}
}

func PLanStatusToEnum(status int) PlanStatusEnum {
	switch status {
	case PlanStatusEditable:
		return PlanStatusEditable
	case PlanStatusActive:
		return PlanStatusActive
	case PlanStatusSoftArchive:
		return PlanStatusSoftArchive
	case PlanStatusHardArchive:
		return PlanStatusHardArchive
	default:
		return PlanStatusEditable
	}
}

func (status PlanType) Description() string {
	switch status {
	case PlanTypeMain:
		return "MainPlan"
	case PlanTypeRecurringAddon:
		return "RecurringAddon"
	case PlanTypeOnetime:
		return "OnetimePlan"
	default:
		return "MainPlan"
	}
}

func PlanTypeToEnum(one int) PlanType {
	switch one {
	case PlanTypeMain:
		return PlanTypeMain
	case PlanTypeRecurringAddon:
		return PlanTypeRecurringAddon
	case PlanTypeOnetime:
		return PlanTypeOnetime
	default:
		return PlanTypeMain
	}
}

func PlanTypeDescriptionToEnum(one string) PlanType {
	switch one {
	case "MainPlan":
		return PlanTypeMain
	case "RecurringAddon":
		return PlanTypeRecurringAddon
	case "OnetimePlan":
		return PlanTypeOnetime
	default:
		return PlanTypeMain
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
