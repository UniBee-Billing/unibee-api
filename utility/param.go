package utility

import "github.com/gogf/gf/v2/os/genv"

func GetEnvParam(name string) string {
	v := genv.GetWithCmd(name)
	if v != nil {
		return v.String()
	}
	return ""
}
