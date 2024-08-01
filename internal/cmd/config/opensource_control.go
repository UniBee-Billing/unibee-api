package config

func IsOpenSourceVersion() bool {
	if GetConfigInstance().Mode == "standalone" || GetConfigInstance().Mode == "cloud" {
		return false
	}
	return true
}
