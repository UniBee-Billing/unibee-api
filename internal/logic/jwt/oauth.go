package jwt

import (
	"context"
	"fmt"
	"time"
	"unibee/internal/cmd/config"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
)

// OAuthJsClaims represents the claims in Auth.js JWT token
type OAuthJsClaims struct {
	Email         string `json:"email"`
	Provider      string `json:"provider"`
	ProviderId    string `json:"providerId"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	EmailVerified bool   `json:"emailVerified"`
	Sign          string `json:"sign"`
	Iat           int64  `json:"iat"`
	Exp           int64  `json:"exp"`
	jwt.RegisteredClaims
}

func GenerateOAuthSign(email, name string) string {
	secret := config.GetConfigInstance().Server.JwtKey
	payload := fmt.Sprintf("%s|%s|%s|%s", secret, email, name, secret)
	return utility.MD5(payload)
}

// ValidateOAuthJsJWT validates Auth.js JWT token and returns user information
func ValidateOAuthJsJWT(ctx context.Context, tokenString string) (*OAuthJsClaims, error) {
	// Get OAuth token secret from config
	secret := config.GetConfigInstance().OAuth.TokenSecret
	utility.Assert(len(secret) > 0, "OAuth token secret not configured")

	// Parse and validate JWT token without time validation
	token, err := jwt.ParseWithClaims(tokenString, &OAuthJsClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}, jwt.WithoutClaimsValidation()) // Skip all claims validation, including time validation

	if err != nil {
		g.Log().Errorf(ctx, "Failed to parse Auth.js JWT token: %v", err)
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// Manually validate token validity (excluding time)
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*OAuthJsClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Check if token exceeds one month validity period
	if claims.Exp > 0 {
		oneMonthAgo := time.Now().AddDate(0, -1, 0).Unix() // One month ago in Unix timestamp
		if claims.Exp < oneMonthAgo {
			return nil, fmt.Errorf("oauth token is expired")
		}
	}

	//// Validate required fields
	//if claims.Email == "" {
	//	return nil, fmt.Errorf("email is required in token")
	//}
	if claims.Provider == "" {
		return nil, fmt.Errorf("provider is required in token")
	}
	if claims.ProviderId == "" {
		return nil, fmt.Errorf("providerId is required in token")
	}

	sign := GenerateOAuthSign(claims.Email, claims.Name)

	if claims.Sign != sign {
		g.Log().Errorf(ctx, "ValidateOAuthJsJWT Invalid token sign:%s verifySign:%s email:%s name:%s", claims.Sign, sign, claims.Email, claims.Name)
		return nil, fmt.Errorf("invalid oauth token sign")
	}

	return claims, nil
}

// GetOAuthJsTokenFromHeader extracts Auth.js JWT token from request header
func GetOAuthJsTokenFromHeader(ctx context.Context) string {
	// Try different header keys for Auth.js token
	authJsToken := g.RequestFromCtx(ctx).Header.Get("X-Auth-JS-Token")
	if authJsToken == "" {
		authJsToken = g.RequestFromCtx(ctx).Header.Get("X-Auth-Token")
	}
	if authJsToken == "" {
		authJsToken = g.RequestFromCtx(ctx).Header.Get("X-OAuth-Token")
	}
	return authJsToken
}
