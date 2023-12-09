package utility

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

func TryLock(ctx context.Context, redisKey string, second int64) bool {
	result, err := g.Redis().Do(ctx, "SET", redisKey, "1", "NX", "EX", second)
	return result != nil && err == nil
}
func ReleaseLock(ctx context.Context, redisKey string) bool {
	_, err := g.Redis().Del(ctx, redisKey)
	return err == nil
}
