package jwt

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
	"unibee/internal/consts"
	"unibee/utility"
)

type TokenType string

const (
	TOKEN_PREFIX            = "UniBee.Portal."
	TOKENTYPEUSER           = "USER"
	TOKENTYPEMERCHANTMember = "MERCHANT_MEMBER"
)

type TokenClaims struct {
	TokenType  TokenType `json:"tokenType"`
	Id         uint64    `json:"id"`
	Email      string    `json:"email"`
	MerchantId uint64    `json:"merchantId"`
	jwt.RegisteredClaims
}

func IsPortalToken(token string) bool {
	return strings.HasPrefix(token, TOKEN_PREFIX)
}

func ParsePortalToken(accessToken string) *TokenClaims {
	utility.Assert(len(consts.GetConfigInstance().Server.TokenKey) > 0, "server error: tokenKey is nil")
	accessToken = strings.Replace(accessToken, TOKEN_PREFIX, "", 1)
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(consts.GetConfigInstance().Server.TokenKey), nil
	})
	return parsedAccessToken.Claims.(*TokenClaims)
}

func CreatePortalToken(tokenType TokenType, merchantId uint64, id uint64, email string) (string, error) {
	utility.Assert(len(consts.GetConfigInstance().Server.TokenKey) > 0, "server error: tokenKey is nil")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"tokenType":  tokenType,
			"merchantId": merchantId,
			"id":         id,
			"email":      email,
			"exp":        time.Now().Add(time.Hour * 1).Unix(),
		})

	tokenString, err := token.SignedString([]byte(consts.GetConfigInstance().Server.TokenKey))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", TOKEN_PREFIX, tokenString), nil
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
