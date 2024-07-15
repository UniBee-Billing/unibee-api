package consts

const (
	ProrationUsingUniBeeCompute         = true
	SubscriptionCycleUnderUniBeeControl = true
	SubPendingTimeout                   = 3 * 24 * 60 * 60
)

type SubscriptionStatusEnum int

const (
	SubTypeDefault            = 0
	SubTypeUniBeeControl      = 1
	SubStatusInit             = 0
	SubStatusPending          = 1
	SubStatusActive           = 2
	SubStatusPendingInActive  = 3 // deprecated
	SubStatusCancelled        = 4
	SubStatusExpired          = 5
	SubStatusSuspended        = 6
	SubStatusIncomplete       = 7
	SubStatusProcessing       = 8
	SubStatusFailed           = 9
	PendingSubStatusInit      = 0
	PendingSubStatusCreate    = 1
	PendingSubStatusFinished  = 2
	PendingSubStatusCancelled = 3
)

func (status SubscriptionStatusEnum) Description() string {
	switch status {
	case SubStatusInit:
		return "Init"
	case SubStatusPending:
		return "Pending"
	case SubStatusActive:
		return "Active"
	case SubStatusCancelled:
		return "Cancelled"
	case SubStatusExpired:
		return "Expired"
	case SubStatusSuspended:
		return "Suspended"
	case SubStatusIncomplete:
		return "Incomplete"
	case SubStatusProcessing:
		return "Processing"
	default:
		return "Init"
	}
}

func SubStatusToEnum(status int) SubscriptionStatusEnum {
	switch status {
	case SubStatusInit:
		return SubStatusInit
	case SubStatusPending:
		return SubStatusPending
	case SubStatusActive:
		return SubStatusActive
	case SubStatusCancelled:
		return SubStatusCancelled
	case SubStatusExpired:
		return SubStatusExpired
	case SubStatusSuspended:
		return SubStatusSuspended
	case SubStatusIncomplete:
		return SubStatusIncomplete
	case SubStatusProcessing:
		return SubStatusProcessing
	default:
		return SubStatusInit
	}
}

const (
	SubTimeLineStatusPending    = 0
	SubTimeLineStatusProcessing = 1
	SubTimeLineStatusFinished   = 2
	SubTimeLineStatusCancelled  = 3
	SubTimeLineStatusExpired    = 4
)
