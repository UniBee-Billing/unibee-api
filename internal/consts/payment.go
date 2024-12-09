package consts

const (
	BizTypeOneTime        = 1
	BizTypeInvoice        = 2
	BizTypeSubscription   = 3
	BizTypeCreditRecharge = 4

	WaitingAuthorized = 0
	Authorized        = 1
	CaptureRequest    = 2
)

type PaymentStatusEnum int

const (
	PaymentCreated   = 10
	PaymentSuccess   = 20
	PaymentFailed    = 30
	PaymentCancelled = 40
)

func (action PaymentStatusEnum) Description() string {
	switch action {
	case PaymentCreated:
		return "PAYMENT_CREATED"
	case PaymentSuccess:
		return "PAYMENT_SUCCESS"
	case PaymentFailed:
		return "PAYMENT_FAILED"
	default:
		return "PAYMENT_CREATED"
	}
}

type RefundStatusEnum int

const (
	RefundCreated       = 10
	RefundSuccess       = 20
	RefundFailed        = 30
	RefundCancelled     = 40
	RefundReverse       = 50
	RefundTypeGateway   = 1
	RefundTypeMarked    = 2
	TimelineTypePayment = 0
	TimelineTypeRefund  = 1
)

func (action RefundStatusEnum) Description() string {
	switch action {
	case RefundCreated:
		return "REFUND_CREATED"
	case RefundSuccess:
		return "REFUND_SUCCESS"
	case RefundFailed:
		return "REFUND_FAILED"
	case RefundCancelled:
		return "REFUND_CANCELLED"
	case RefundReverse:
		return "REFUND_REVERSE"
	default:
		return "REFUND_CREATED"
	}
}
