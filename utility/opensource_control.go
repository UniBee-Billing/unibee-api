package utility

import (
	"unibee/internal/cmd/config"
)

func IsOpenSourceVersion() bool {
	if config.GetConfigInstance().Mode == "standalone" || config.GetConfigInstance().Mode == "cloud" {
		return false
	}
	return true
}
