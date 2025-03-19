package utility

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
)

func TryLock(ctx context.Context, redisKey string, second int64) bool {
	result, err := g.Redis().Do(ctx, "SET", redisKey, "1", "NX", "EX", second)
	if err != nil {
		fmt.Printf("UtilityLock_TryLock err: %s\n", err.Error())
	}
	return result != nil && !result.IsNil() && err == nil
}
func ReleaseLock(ctx context.Context, redisKey string) bool {
	_, err := g.Redis().Del(ctx, redisKey)
	if err != nil {
		fmt.Printf("UtilityLock_ReleaseLock err: %s\n", err.Error())
	}
	return err == nil
}
