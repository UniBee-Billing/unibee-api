package consts

const (
	BizTypeOneTime      = 1
	BizTypeInvoice      = 2
	BizTypeSubscription = 3

	WaitingAuthorized = 0 //
	Authorized        = 1 //
	CaptureRequest    = 2 //

)

type PayStatusEnum int

const (
	PaymentCreated   = 10 //
	PaymentSuccess   = 20 //
	PaymentFailed    = 30 //
	PaymentCancelled = 40 //
)

func (action PayStatusEnum) Description() string {
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
	RefundIng     = 10 //
	RefundSuccess = 20 //
	RefundFailed  = 30 //
	RefundReverse = 40 //
)

func (action RefundStatusEnum) Description() string {
	switch action {
	case RefundIng:
		return "REFUND_ING"
	case RefundSuccess:
		return "REFUND_SUCCESS"
	case RefundFailed:
		return "REFUND_FAILED"
	default:
		return "REFUND_ING"
	}
}
