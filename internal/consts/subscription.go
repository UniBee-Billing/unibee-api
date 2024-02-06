package consts

const (
	NonEffectImmediatelyUsePendingUpdate = false
	ProrationUsingUniBeeCompute          = true
	SubscriptionCycleUnderUniBeeControl  = true
)

type InvoiceStatusEnum int

const (
	InvoiceStatusInit       = 0 //
	InvoiceStatusPending    = 1 //
	InvoiceStatusProcessing = 2 //
	InvoiceStatusPaid       = 3 //
	InvoiceStatusFailed     = 4 //
	InvoiceStatusCancelled  = 5 //
)

type SubscriptionStatusEnum int

const (
	SubTypeDefault            = 0
	SubTypeUniBeeControl      = 1
	SubStatusInit             = 0 //
	SubStatusCreate           = 1 //
	SubStatusActive           = 2 //
	SubStatusPendingInActive  = 3 //
	SubStatusCancelled        = 4 //
	SubStatusExpired          = 5 //
	SubStatusSuspended        = 6 //
	SubStatusIncomplete       = 7 //
	PendingSubStatusInit      = 0
	PendingSubStatusCreate    = 1
	PendingSubStatusFinished  = 2
	PendingSubStatusCancelled = 3
)

func (status SubscriptionStatusEnum) Description() string {
	switch status {
	case SubStatusInit:
		return "SubStatusInit"
	case SubStatusCreate:
		return "SubStatusCreate"
	case SubStatusActive:
		return "SubStatusActive"
	case SubStatusSuspended:
		return "SubStatusSuspended"
	default:
		return "SubStatusInit"
	}
}

type SubscriptionPlanType int

const (
	PlanTypeMain  = 1
	PlanTypeAddon = 2
)

type SubscriptionPlanStatusEnum int

const (
	PlanStatusEditable           = 1
	PlanStatusActive             = 2
	PlanStatusExpired            = 3
	PlanPublishStatusPublished   = 2
	PlanPublishStatusUnPublished = 1
)

func (status SubscriptionPlanStatusEnum) Description() string {
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

type SubscriptionGatewayPlanStatusEnum int

const (
	GatewayPlanStatusInit     = 0
	GatewayPlanStatusCreate   = 1
	GatewayPlanStatusActive   = 2
	GatewayPlanStatusInActive = 3
)

func (status SubscriptionGatewayPlanStatusEnum) Description() string {
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
