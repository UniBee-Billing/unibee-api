package time

import "github.com/gogf/gf/v2/os/gtime"

func init() {
	// Set Global TimeZone, Should Set Before Standard Time Package Init
	err := gtime.SetTimeZone("UTC")
	if err != nil {
		panic(err)
	}
}
