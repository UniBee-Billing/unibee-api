package consts

const (
	UsePendingUpdate                              = true
	DownGradeUsePendingUpdate                     = true
	PendingSubUpdateEffectImmediateWithOutChannel = false
)

type InvoiceStatusEnum int

const (
	InvoiceStatusInit       = 0 //初始化
	InvoiceStatusPending    = 1 // 草稿-渠道状态
	InvoiceStatusProcessing = 2 //打开可用于支付-渠道状态
	InvoiceStatusPaid       = 3 //已支付-渠道状态
	InvoiceStatusFailed     = 4 //支付失败-渠道状态
	InvoiceStatusCancelled  = 5 //取消-渠道状态
)

type SubscriptionStatusEnum int

const (
	SubStatusInit             = 0 //初始化
	SubStatusCreate           = 1 //创建-渠道状态
	SubStatusActive           = 2 //有效-渠道状态
	SubStatusPendingInActive  = 3 //PendingInActive
	SubStatusCancelled        = 4 //取消-渠道状态
	SubStatusExpired          = 5 //过期-渠道状态
	SubStatusSuspended        = 6 //暂停-渠道状态
	SubStatusIncomplete       = 7 //付款支付失败-渠道状态
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

type SubscriptionPlanChannelStatusEnum int

const (
	PlanChannelStatusInit     = 0 //初始化
	PlanChannelStatusCreate   = 1 //创建-渠道状态
	PlanChannelStatusActive   = 2 //有效-渠道状态
	PlanChannelStatusInActive = 3 //无效-渠道状态
)

func (status SubscriptionPlanChannelStatusEnum) Description() string {
	switch status {
	case PlanChannelStatusInit:
		return "STATUES_CREATE"
	case PlanChannelStatusCreate:
		return "STATUES_CREATE"
	case PlanChannelStatusActive:
		return "PlanChannelStatusActive"
	case PlanChannelStatusInActive:
		return "PlanChannelStatusInActive"
	default:
		return "PlanChannelStatusInit"
	}
}
