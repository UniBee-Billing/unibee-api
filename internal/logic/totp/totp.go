package totp

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"strings"
	"time"
	"unibee/internal/cmd/config"
	"unibee/internal/logic/merchant_config"
	"unibee/internal/logic/merchant_config/update"
	"unibee/utility"
)

const (
	KeyMerchantTotpGlobalStatus = "KEY_MERCHANT_TOTP_GLOBAL_STATUS"
	TotpTypeGeneral             = 1 //Google Authenticator | Microsoft Authenticator | Authy | 1Password / LastPass | FreeOTP | Other TOTP
)

func GetMerchantTotpGlobalConfig(ctx context.Context, merchantId uint64) bool {
	keyConfig := merchant_config.GetMerchantConfig(ctx, merchantId, KeyMerchantTotpGlobalStatus)
	if keyConfig != nil {
		value := keyConfig.ConfigValue
		if value == "true" {
			return true
		}
	}
	return false
}

func UpdateMerchantTotpGlobalConfig(ctx context.Context, merchantId uint64, activate bool) {
	if activate {
		_ = update.SetMerchantConfig(ctx, merchantId, KeyMerchantTotpGlobalStatus, "true")
	} else {
		_ = update.SetMerchantConfig(ctx, merchantId, KeyMerchantTotpGlobalStatus, "false")
	}
}

func GetMerchantMemberTotpSecret(ctx context.Context, totpType int, email string) (secret string, url string, err error) {
	issuer := fmt.Sprintf("UniBee (%s)", email)
	if strings.Contains(config.GetConfigInstance().Server.GetServerPath(), "beta") {
		issuer = fmt.Sprintf("UniBee-beta (%s)", email)
	} else if !config.GetConfigInstance().IsProd() {
		issuer = fmt.Sprintf("UniBee-%s (%s)", config.GetConfigInstance().Env, email)
	}
	if totpType > 0 {
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      issuer, // Application Name
			AccountName: email,  // Account
		})
		if err != nil {
			return "", "", err
		}
		secret = key.Secret()
		url = key.URL()
	} else {
		return "", "", errors.New("2FA type not supported")
	}
	g.Log().Infof(ctx, "GetMerchantMemberTotpSecret totpType:%d email:%s secret:%s url:%s", totpType, email, secret, url)
	return secret, url, nil
}

func ValidateMerchantMemberTotp(ctx context.Context, totpType int, email string, secret string, code string, clientIdentity string) (result bool) {
	if len(secret) < 0 || len(code) < 0 {
		return false
	}
	if !config.GetConfigInstance().IsProd() && code == "666666" {
		g.Log().Infof(ctx, "ValidateMerchantMemberTotp totpType:%d email:%s secret:%s with magic code:%s", totpType, email, secret, code)
		SaveClientIdentity(ctx, email, clientIdentity)
		return true
	}
	if totpType > 0 {
		validateResult, err := totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
			Period:    30,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		})
		utility.AssertError(err, "Validate2FA")
		result = validateResult
	} else {
		utility.Assert(false, "2FA type not supported")
	}
	g.Log().Infof(ctx, "ValidateMerchantMemberTotp totpType:%d email:%s secret:%s code:%s", totpType, email, secret, code)
	if result {
		SaveClientIdentity(ctx, email, clientIdentity)
	}

	return result
}
