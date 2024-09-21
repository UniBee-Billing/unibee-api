package jwt

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
	"unibee/api/bean/detail"
	"unibee/internal/cmd/config"
	"unibee/internal/model"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

const (
	TOKEN_PREFIX            = "UniBee.Portal."
	TOKENTYPEUSER           = "USER"
	TOKENTYPEMERCHANTMember = "MERCHANT_MEMBER"
)

func IsPortalToken(token string) bool {
	return strings.HasPrefix(token, TOKEN_PREFIX)
}

func ParsePortalToken(accessToken string) *model.TokenClaims {
	utility.Assert(len(config.GetConfigInstance().Server.JwtKey) > 0, "server error: tokenKey is nil")
	accessToken = strings.Replace(accessToken, TOKEN_PREFIX, "", 1)
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &model.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfigInstance().Server.JwtKey), nil
	})
	return parsedAccessToken.Claims.(*model.TokenClaims)
}

func CreatePortalToken(tokenType model.TokenType, merchantId uint64, id uint64, email string, lang string) (string, error) {
	utility.Assert(len(config.GetConfigInstance().Server.JwtKey) > 0, "server error: tokenKey is nil")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"tokenType":  tokenType,
			"merchantId": merchantId,
			"id":         id,
			"email":      email,
			"lang":       lang,
			"exp":        time.Now().Add(time.Hour * 1).Unix(),
		})

	tokenString, err := token.SignedString([]byte(config.GetConfigInstance().Server.JwtKey))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", TOKEN_PREFIX, tokenString), nil
}

func CreateMemberPortalToken(ctx context.Context, tokenType model.TokenType, merchantId uint64, memberId uint64, email string) (string, error) {
	utility.Assert(len(config.GetConfigInstance().Server.JwtKey) > 0, "server error: tokenKey is nil")
	one := query.GetMerchantMemberById(ctx, memberId)
	utility.Assert(one != nil, "member not found")
	permissionKey := GetMemberPermissionKey(ctx, one)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"tokenType":     tokenType,
			"merchantId":    merchantId,
			"id":            memberId,
			"email":         email,
			"exp":           time.Now().Add(time.Hour * 1).Unix(),
			"permissionKey": permissionKey,
		})

	tokenString, err := token.SignedString([]byte(config.GetConfigInstance().Server.JwtKey))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", TOKEN_PREFIX, tokenString), nil
}

func GetMemberPermissionKey(ctx context.Context, one *entity.MerchantMember) string {
	isOwner, permission := detail.ConvertMemberPermissions(ctx, one)
	permissionKey := fmt.Sprintf("%v_%s", isOwner, utility.MD5(utility.MarshalToJsonString(permission)))
	return permissionKey
}

func getAuthTokenRedisKey(token string) string {
	return fmt.Sprintf("auth#%s#%s", config.GetConfigInstance().Env, token)
}

func PutAuthTokenToCache(ctx context.Context, token string, value string) bool {
	err := g.Redis().SetEX(ctx, getAuthTokenRedisKey(token), value, config.GetConfigInstance().Auth.Login.Expire)
	if err != nil {
		return false
	}
	return true
}

func IsAuthTokenAvailable(ctx context.Context, token string) bool {
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
	expire, err := g.Redis().Expire(ctx, getAuthTokenRedisKey(token), config.GetConfigInstance().Auth.Login.Expire)
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
