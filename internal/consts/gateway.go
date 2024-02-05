package consts

type GatewayTypeEnum int

const (
	GatewayTypeOneTimePayment = 0 //
	GatewayTypeSubscription   = 1 //
)

func (status GatewayTypeEnum) Description() string {
	switch status {
	case GatewayTypeOneTimePayment:
		return "GatewayTypeOneTimePayment"
	case GatewayTypeSubscription:
		return "GatewayTypeSubscription"
	default:
		return "GatewayTypeOneTimePayment"
	}
}
