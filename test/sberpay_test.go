package test

import (
	"context"
	"testing"

	"unibee/internal/logic/gateway/api"
	entity "unibee/internal/model/entity/default"

	"github.com/stretchr/testify/assert"
)

func TestSberPayGatewayInfo(t *testing.T) {
	ctx := context.Background()
	sberPay := &api.SberPay{}

	info := sberPay.GatewayInfo(ctx)

	assert.NotNil(t, info)
	assert.Equal(t, "SberPay", info.Name)
	assert.Equal(t, "SberPay", info.DisplayName)
	assert.Equal(t, "SberPay mobile payment gateway for Russian market", info.Description)
	assert.Equal(t, "https://www.sberbank.ru/", info.GatewayWebsiteLink)
	assert.Equal(t, "https://api.sber.ru", info.Host)
	assert.Equal(t, "gateway_key", info.PublicKeyName)
	assert.Equal(t, "gateway_secret", info.PrivateSecretName)
	assert.Equal(t, "https://api.unibee.dev/oss/file/dbyr5ts5ykzlbbn5ly.png", info.GatewayLogo)
	assert.Equal(t, int64(200), info.Sort)
	assert.False(t, info.AutoChargeEnabled)
	assert.True(t, info.IsStaging)
}

func TestSberPayGatewayTest(t *testing.T) {
	// This test is a placeholder for future implementation
	// In actual environment, real API keys need to be provided for testing
	t.Skip("Skipping SberPay gateway test - requires real API keys")
}

func TestSberPayUserCreate(t *testing.T) {
	ctx := context.Background()
	sberPay := &api.SberPay{}

	// Mock gateway and user
	gateway := &entity.MerchantGateway{
		GatewayKey:    "test_key",
		GatewaySecret: "test_secret",
	}

	user := &entity.UserAccount{
		Id: 123,
	}

	resp, err := sberPay.GatewayUserCreate(ctx, gateway, user)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "user_123", resp.GatewayUserId)
}

func TestSberPayMerchantBalancesQuery(t *testing.T) {
	ctx := context.Background()
	sberPay := &api.SberPay{}

	gateway := &entity.MerchantGateway{
		GatewayKey:    "test_key",
		GatewaySecret: "test_secret",
	}

	resp, err := sberPay.GatewayMerchantBalancesQuery(ctx, gateway)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.AvailableBalance, 1)
	assert.Len(t, resp.PendingBalance, 1)

	balance := resp.AvailableBalance[0]
	assert.Equal(t, int64(0), balance.Amount)
	assert.Equal(t, "RUB", balance.Currency)
}

func TestSberPayUnsupportedMethods(t *testing.T) {
	ctx := context.Background()
	sberPay := &api.SberPay{}

	gateway := &entity.MerchantGateway{
		GatewayKey:    "test_key",
		GatewaySecret: "test_secret",
	}

	// Test unsupported methods
	_, err := sberPay.GatewayCryptoFiatTrans(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SberPay does not support crypto transactions")

	_, err = sberPay.GatewayUserAttachPaymentMethodQuery(ctx, gateway, 123, "method_id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SberPay does not support payment method management")

	_, err = sberPay.GatewayUserDeAttachPaymentMethodQuery(ctx, gateway, 123, "method_id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SberPay does not support payment method management")

	_, err = sberPay.GatewayUserPaymentMethodListQuery(ctx, gateway, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SberPay does not support payment method management")

	_, err = sberPay.GatewayUserCreateAndBindPaymentMethod(ctx, gateway, 123, "RUB", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SberPay does not support payment method management")
}
