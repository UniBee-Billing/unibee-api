package utility

import (
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
)

type localCode struct {
	code    int         // Error code, usually an integer.
	message string      // Brief message for this error code.
	detail  interface{} // As type of interface, it is mainly designed as an extension field for error code.
}

// Code returns the integer number of current error code.
func (c localCode) Code() int {
	return c.code
}

// Message returns the brief message for current error code.
func (c localCode) Message() string {
	return c.message
}

// Detail returns the detailed information of current error code,
// which is mainly designed as an extension field for error code.
func (c localCode) Detail() interface{} {
	return c.detail
}

// String returns current error code as a string.
func (c localCode) String() string {
	if c.detail != nil {
		return fmt.Sprintf(`%d:%s %v`, c.code, c.message, c.detail)
	}
	if c.message != "" {
		return fmt.Sprintf(`%d:%s`, c.code, c.message)
	}
	return fmt.Sprintf(`%d`, c.code)
}

var (
	GatewayError = gcode.New(70, "Gateway Failed", nil)
)
