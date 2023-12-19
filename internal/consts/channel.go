package consts

type PayChannelTypeEnum int

const (
	PayChannelTypePayment      = 0 //
	PayChannelTypeSubscription = 1 //
)

func (status PayChannelTypeEnum) Description() string {
	switch status {
	case PayChannelTypePayment:
		return "PayChannelTypePayment"
	case PayChannelTypeSubscription:
		return "PayChannelTypeSubscription"
	default:
		return "PayChannelTypePayment"
	}
}
