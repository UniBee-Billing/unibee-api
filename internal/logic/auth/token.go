package auth

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"unibee-api/internal/consts"
)

var secretKey = []byte("3^&secret-key-for-UniBee*1!8*")

func CreateToken(email string, userId uint64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"id":    userId,
			"exp":   time.Now().Add(time.Hour * 1).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getAuthTokenRedisKey(token string) string {
	return fmt.Sprintf("auth#%s#%s", consts.GetConfigInstance().Env, token)
}

func PutAuthTokenToCache(ctx context.Context, token string, value string) bool {
	err := g.Redis().SetEX(ctx, getAuthTokenRedisKey(token), value, consts.GetConfigInstance().Auth.Login.Expire)
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

func ResetAuthTokenTTL(ctx context.Context, token string) bool {
	expire, err := g.Redis().Expire(ctx, getAuthTokenRedisKey(token), consts.GetConfigInstance().Auth.Login.Expire)
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
