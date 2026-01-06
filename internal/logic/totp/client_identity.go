package totp

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean/detail"
	"unibee/utility"
)

func IsClientIdentityExist(ctx context.Context, email string, clientIdentity string) bool {
	var clientIdentityMap = map[string]int64{}
	data, err := g.Redis().Get(ctx, detail.GetClientIdentityListCacheKey(email))
	if err == nil && data != nil && data.String() != "" {
		_ = utility.UnmarshalFromJsonString(data.String(), &clientIdentityMap)
	}
	if _, ok := clientIdentityMap[clientIdentity]; ok {
		return true
	}
	return false
}

func IsClientIdentityValid(ctx context.Context, email string, clientIdentity string) bool {
	var clientIdentityMap = map[string]int64{}
	data, err := g.Redis().Get(ctx, detail.GetClientIdentityListCacheKey(email))
	if err == nil && data != nil && data.String() != "" {
		_ = utility.UnmarshalFromJsonString(data.String(), &clientIdentityMap)
	}
	if lastTime, ok := clientIdentityMap[clientIdentity]; ok {
		if gtime.Now().Timestamp()-lastTime < detail.MemberDeviceExpireDays*24*60*60 {
			return true
		}
	}
	return false
}

func SaveClientIdentity(ctx context.Context, email string, clientIdentity string) {
	var clientIdentityMap = map[string]int64{}
	data, err := g.Redis().Get(ctx, detail.GetClientIdentityListCacheKey(email))
	if err == nil && data != nil && data.String() != "" {
		_ = utility.UnmarshalFromJsonString(data.String(), &clientIdentityMap)
	}
	clientIdentityMap[clientIdentity] = gtime.Now().Timestamp()
	_ = g.Redis().SetEX(ctx, detail.GetClientIdentityListCacheKey(email), clientIdentityMap, 60*24*60*60)
}

func DeleteClientIdentity(ctx context.Context, email string, clientIdentity string) {
	var clientIdentityMap = map[string]int64{}
	data, err := g.Redis().Get(ctx, detail.GetClientIdentityListCacheKey(email))
	if err == nil && data != nil && data.String() != "" {
		_ = utility.UnmarshalFromJsonString(data.String(), &clientIdentityMap)
	}
	delete(clientIdentityMap, clientIdentity)
	_ = g.Redis().SetEX(ctx, detail.GetClientIdentityListCacheKey(email), clientIdentityMap, 60*24*60*60)
}
