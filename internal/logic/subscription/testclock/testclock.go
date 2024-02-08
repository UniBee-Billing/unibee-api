package testclock

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee-api/internal/consts"
)

type TestClockWalkRes struct {
}

func WalkSubscriptionToTestClock(ctx context.Context, newTestClock int64) (*TestClockWalkRes, error) {
	if consts.GetConfigInstance().IsProd() {
		return nil, gerror.New("Test Does Not Work For Prod Env")
	}
	//TestClock Walk
	return &TestClockWalkRes{}, nil
}
