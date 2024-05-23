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
