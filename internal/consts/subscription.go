package consts

type InvoiceStatusEnum int

const (
	InvoiceStatusInit        = 0 //初始化
	InvoiceStatusDraft       = 1 // 草稿-渠道状态
	InvoiceStatusOpen        = 2 //打开-渠道状态
	InvoiceStatusPaid        = 3 //已支付-渠道状态
	InvoiceStatusUnCollected = 4 //未收集-渠道状态
	InvoiceStatusVoid        = 5 // Void-渠道状态
)

type SubscriptionStatusEnum int

const (
	SubStatusInit      = 0 //初始化
	SubStatusCreate    = 1 //创建-渠道状态
	SubStatusActive    = 2 //有效-渠道状态
	SubStatusSuspended = 3 //暂停-渠道状态
	SubStatusCancelled = 4 //取消-渠道状态
	SubStatusExpired   = 5 //过期-渠道状态
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
	PlanStatusEditable  = 1
	PlanStatusPublished = 2
	PlanStatusExpired   = 3
)

func (status SubscriptionPlanStatusEnum) Description() string {
	switch status {
	case PlanStatusEditable:
		return "PlanStatusEditable"
	case PlanStatusPublished:
		return "PlanStatusPublished"
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
