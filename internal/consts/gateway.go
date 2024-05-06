package consts

type GatewayTypeEnum int

const (
	GatewayTypeDefault      = 1 //
	GatewayTypeCrypto       = 2 //
	GatewayTypeWireTransfer = 3 //
)

func (status GatewayTypeEnum) Description() string {
	switch status {
	case GatewayTypeDefault:
		return "GatewayTypeDefault"
	case GatewayTypeCrypto:
		return "GatewayTypeCrypto"
	case GatewayTypeWireTransfer:
		return "GatewayTypeWireTransfer"
	default:
		return "GatewayTypeDefault"
	}
}
