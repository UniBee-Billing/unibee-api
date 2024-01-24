package auth

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
)

func getAuthTokenRedisKey(token string) string {
	return fmt.Sprintf("auth#%s", token)
}

func PutAuthTokenToCache(ctx context.Context, token string, value string, ttlSecond int64) bool {
	err := g.Redis().SetEX(ctx, getAuthTokenRedisKey(token), value, ttlSecond)
	if err != nil {
		return false
	}
	return true
}

func IsAuthTokenExpired(ctx context.Context, token string) bool {
	get, err := g.Redis().Get(ctx, getAuthTokenRedisKey(token))
	if err != nil {
		return false
	}
	if get != nil && len(get.String()) > 0 {
		return true
	}
	return false
}

func SetAuthTokenNewTTL(ctx context.Context, token string, newTTLSecond int64) bool {
	expire, err := g.Redis().Expire(ctx, getAuthTokenRedisKey(token), newTTLSecond)
	if err != nil {
		return false
	}
	return expire == 1
}

func DelAuthToken(ctx context.Context, token string) bool {
	_, err := g.Redis().Del(ctx, getAuthTokenRedisKey(token))
	if err != nil {
		return false
	}
	return true
}
