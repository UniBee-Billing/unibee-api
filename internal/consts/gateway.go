package consts

type GatewayTypeEnum int

const (
	GatewayTypeCard         = 1 //
	GatewayTypeCrypto       = 2 //
	GatewayTypeWireTransfer = 3 //
)

func (status GatewayTypeEnum) Description() string {
	switch status {
	case GatewayTypeCard:
		return "GatewayTypeCard"
	case GatewayTypeCrypto:
		return "GatewayTypeCrypto"
	case GatewayTypeWireTransfer:
		return "GatewayTypeWireTransfer"
	default:
		return "GatewayTypeCard"
	}
}
