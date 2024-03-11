package consts

type GatewayTypeEnum int

const (
	GatewayTypeDefault = 1 //
	GatewayTypeCrypto  = 2 //
)

func (status GatewayTypeEnum) Description() string {
	switch status {
	case GatewayTypeDefault:
		return "GatewayTypeDefault"
	case GatewayTypeCrypto:
		return "GatewayTypeCrypto"
	default:
		return "GatewayTypeDefault"
	}
}
