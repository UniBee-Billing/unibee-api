package system

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"unibee/api/system/auth"
	"unibee/internal/cmd/config"
	"unibee/internal/logic/jwt"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerAuth) TokenGenerator(ctx context.Context, req *auth.TokenGeneratorReq) (res *auth.TokenGeneratorRes, err error) {
	if len(req.Env) == 0 {
		req.Env = config.GetConfigInstance().Env
	}
	database := config.GetConfigInstance().RedisConfig.Default.DB
	if req.RedisDatabase != nil {
		database = *req.RedisDatabase
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.GetConfigInstance().RedisConfig.Default.Address,
		Password: config.GetConfigInstance().RedisConfig.Default.Pass,
		DB:       database,
	})
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("redis error:%s\n", err.Error())
		}
	}(client)

	if req.PortalType == 0 {
		one := query.GetMerchantMemberByEmail(ctx, req.Email)
		utility.Assert(one != nil, "email not found")
		token, err := jwt.CreateMemberPortalToken(ctx, jwt.TOKENTYPEMERCHANTMember, one.MerchantId, one.Id, req.Email)
		utility.AssertError(err, "Server Error")
		//utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("MerchantMember#%d", one.Id)), "Cache Error")
		_, err = client.Do(ctx, "SetEX", getAuthTokenRedisKey(req.Env, token), 86400*7, fmt.Sprintf("MerchantMember#%d", one.Id)).Result()
		utility.AssertError(err, "Server Error")
		return &auth.TokenGeneratorRes{Token: token}, nil
	} else {
		merchant := query.GetMerchantById(ctx, req.MerchantId)
		utility.Assert(merchant != nil, "Invalid merchantId")
		one := query.GetUserAccountByEmail(ctx, req.MerchantId, req.Email)
		utility.Assert(one != nil, "email not found")
		token, err := jwt.CreatePortalToken(jwt.TOKENTYPEUSER, one.MerchantId, one.Id, req.Email, one.Language)
		//utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("User#%d", one.Id), "Cache Error")
		utility.AssertError(err, "Server Error")
		_, err = client.Do(ctx, "SetEX", getAuthTokenRedisKey(req.Env, token), 86400*7, fmt.Sprintf("User#%d", one.Id)).Result()
		utility.AssertError(err, "Server Error")
		return &auth.TokenGeneratorRes{Token: token}, nil
	}
}

func getAuthTokenRedisKey(env string, token string) string {
	return fmt.Sprintf("auth#%s#%s", env, token)
}
