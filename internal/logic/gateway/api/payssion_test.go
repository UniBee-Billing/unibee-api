package api

import (
	"context"
	"testing"
)

func TestForPayssion(t *testing.T) {
	pay := &Payssion{}
	_, _, _ = pay.GatewayTest(context.Background(), "sandbox_6340c0569ae5339c", "hdvh5MkJMCQ5ZhtgatLzukbJXwbRMra4")

}
