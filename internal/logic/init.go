package logic

import (
	"context"
	"unibee/internal/cmd/config"
	"unibee/internal/logic/email"
	"unibee/internal/logic/merchant"
)

func StandaloneInit(ctx context.Context) {
	if config.GetConfigInstance().Mode != "cloud" {
		merchant.StandAloneInit(ctx)
		email.StandAloneInit(ctx)
	}
}
