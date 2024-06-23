package consts

type InvoiceStatusEnum int

const (
	DEFAULT_DAY_UTIL_DUE    = 3
	InvoiceStatusInit       = 0
	InvoiceStatusPending    = 1
	InvoiceStatusProcessing = 2
	InvoiceStatusPaid       = 3
	InvoiceStatusFailed     = 4
	InvoiceStatusCancelled  = 5
	InvoiceStatusReversed   = 6
)

type InvoiceSendStatusEnum int

const (
	InvoiceSendStatusUnSend      = 0
	InvoiceSendStatusSend        = 1
	InvoiceSendStatusUnnecessary = 2
)

func (status InvoiceStatusEnum) Description() string {
	switch status {
	case InvoiceStatusInit:
		return "Init"
	case InvoiceStatusPending:
		return "Pending"
	case InvoiceStatusPaid:
		return "Active"
	case InvoiceStatusCancelled:
		return "Cancelled"
	case InvoiceStatusFailed:
		return "Failed"
	case InvoiceStatusProcessing:
		return "Processing"
	case InvoiceStatusReversed:
		return "Processing"
	default:
		return "Init"
	}
}

func InvoiceStatusToEnum(status int) InvoiceStatusEnum {
	switch status {
	case InvoiceStatusInit:
		return InvoiceStatusInit
	case InvoiceStatusPending:
		return InvoiceStatusPending
	case InvoiceStatusPaid:
		return InvoiceStatusPaid
	case InvoiceStatusCancelled:
		return InvoiceStatusCancelled
	case InvoiceStatusFailed:
		return InvoiceStatusFailed
	case InvoiceStatusProcessing:
		return InvoiceStatusProcessing
	case InvoiceStatusReversed:
		return InvoiceStatusReversed
	default:
		return InvoiceStatusInit
	}
}
