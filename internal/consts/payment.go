package consts

const (
	BIZ_TYPE_ONE_TIME     = 1
	BIZ_TYPE_INVOICE      = 2
	BIZ_TYPE_SUBSCRIPTION = 3

	WAITING_AUTHORIZED = 0 //
	AUTHORIZED         = 1 //
	CAPTURE_REQUEST    = 2 //

)

type PayStatusEnum int

const (
	TO_BE_PAID  = 10 //
	PAY_SUCCESS = 20 //
	PAY_FAILED  = 30 //
	PAY_CANCEL  = 40 //
)

func (action PayStatusEnum) Description() string {
	switch action {
	case TO_BE_PAID:
		return "TO_BE_PAID"
	case PAY_SUCCESS:
		return "PAY_SUCCESS"
	case PAY_FAILED:
		return "PAY_FAILED"
	default:
		return "TO_BE_PAID"
	}
}

type RefundStatusEnum int

const (
	REFUND_ING     = 10 //
	REFUND_SUCCESS = 20 //
	REFUND_FAILED  = 30 //
	REFUND_REVERSE = 40 //
)

func (action RefundStatusEnum) Description() string {
	switch action {
	case REFUND_ING:
		return "REFUND_ING"
	case REFUND_SUCCESS:
		return "REFUND_SUCCESS"
	case REFUND_FAILED:
		return "REFUND_FAILED"
	default:
		return "REFUND_ING"
	}
}
