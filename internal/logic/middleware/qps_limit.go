package middleware

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

const luaScript = `
local key = KEYS[1]
local current = redis.call("INCR", key)

if current == 1 then
    redis.call("PEXPIRE", key, ARGV[1])
end

return current
`

func CheckQPSLimit(ctx context.Context, key string, maxQPS int, expireMs int) (bool, int64) {
	// Run Lua Script
	result, err := g.Redis().Do(ctx, "EVAL", luaScript, "1", key, expireMs)
	if err != nil {
		g.Log().Errorf(ctx, "CheckQPSLimit Error:%s\n", err.Error())
		return true, 0
	}
	if result == nil || result.IsNil() || !result.IsInt() {
		g.Log().Errorf(ctx, "CheckQPSLimit Error result is nil or not int\n")
		return true, 0
	}
	if result.Int() > maxQPS {
		g.Log().Infof(ctx, "CheckQPSLimit reachMax key:%s currenct:%d maxQPS:%d\n", key, result.Int(), maxQPS)
	}
	return result.Int() <= maxQPS, result.Int64()
}
