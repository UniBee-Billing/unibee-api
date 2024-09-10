package util

import "github.com/gogf/gf/v2/errors/gcode"

var (
	GatewayError = gcode.New(70, "Gateway Failed", nil)
)
