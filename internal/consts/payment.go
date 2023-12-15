package consts

const (
	PAYMENT_BIZ_TYPE_ORDER = 1 //支付交易上游订单类型-

	WAITING_AUTHORIZED = 0 //待授权
	AUTHORIZED         = 1 // 已授权
	CAPTURE_REQUEST    = 2 //已发起捕获货已捕获

)

type PayStatusEnum int

const (
	TO_BE_PAID  = 10 //待支付
	PAY_SUCCESS = 20 // 支付成功
	PAY_FAILED  = 30 //支付失败
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
	REFUND_ING     = 10 //退款中
	REFUND_SUCCESS = 20 //退款成功
	REFUND_FAILED  = 30 //退款失败
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
