package merchant

import (
	"fmt"
	"unibee/internal/cmd/config"
	"unibee/utility"
)

func GenerateMerchantAPIKey() string {
	if config.GetConfigInstance().IsProd() {
		return fmt.Sprintf("ub_%s", utility.GenerateRandomAlphanumeric(32))
	} else {
		return fmt.Sprintf("ub_test_%s", utility.GenerateRandomAlphanumeric(32))
	}
}
