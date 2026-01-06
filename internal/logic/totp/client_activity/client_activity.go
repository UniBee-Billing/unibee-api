package client_activity

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/utility"
)

type ClientIdentityActivity struct {
	LastLoginTime    int64 `json:"lastLoginTime"`
	LastActivityTime int64 `json:"lastActivityTime"`
}

func GetClientIdentityActivityCacheKey(clientIdentity string) string {
	return fmt.Sprintf("Merchant#Totp#client#identity#activity#%s", clientIdentity)
}

func GetClientIdentityActivity(ctx context.Context, clientIdentity string) *ClientIdentityActivity {
	var clientIdentityActivity *ClientIdentityActivity
	data, err := g.Redis().Get(ctx, GetClientIdentityActivityCacheKey(clientIdentity))
	if err == nil && data != nil && data.String() != "" {
		_ = utility.UnmarshalFromJsonString(data.String(), &clientIdentityActivity)
	}
	return clientIdentityActivity
}

func UpdateClientIdentityLoginTime(ctx context.Context, clientIdentity string) {
	var clientIdentityActivity = &ClientIdentityActivity{}
	data, err := g.Redis().Get(ctx, GetClientIdentityActivityCacheKey(clientIdentity))
	if err == nil && data != nil && data.String() != "" {
		_ = utility.UnmarshalFromJsonString(data.String(), &clientIdentityActivity)
	}
	clientIdentityActivity.LastLoginTime = gtime.Now().Timestamp()
	clientIdentityActivity.LastActivityTime = gtime.Now().Timestamp()
	_ = g.Redis().SetEX(ctx, GetClientIdentityActivityCacheKey(clientIdentity), utility.MarshalToJsonString(clientIdentityActivity), 90*24*60*60)
}

func UpdateClientIdentityActivityTime(ctx context.Context, clientIdentity string) {
	var clientIdentityActivity = &ClientIdentityActivity{}
	data, err := g.Redis().Get(ctx, GetClientIdentityActivityCacheKey(clientIdentity))
	if err == nil && data != nil && data.String() != "" {
		_ = utility.UnmarshalFromJsonString(data.String(), &clientIdentityActivity)
	}
	clientIdentityActivity.LastActivityTime = gtime.Now().Timestamp()
	_ = g.Redis().SetEX(ctx, GetClientIdentityActivityCacheKey(clientIdentity), utility.MarshalToJsonString(clientIdentityActivity), 90*24*60*60)
}
