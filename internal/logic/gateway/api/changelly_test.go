package api

import (
	"context"
	"testing"
)

func TestForApiKey(t *testing.T) {
	ctx := context.Background()
	changelly := Changelly{}
	err := changelly.GatewayTest(ctx, "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")
	if err != nil {
		return
	}
}