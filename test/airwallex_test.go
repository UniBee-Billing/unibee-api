package test

import (
	"context"
	"testing"

	"unibee/internal/logic/gateway/api"
	gateway_bean "unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"

	"github.com/stretchr/testify/assert"
)

func TestAirwallexGatewayInfo(t *testing.T) {
	ctx := context.Background()
	airwallex := &api.Airwallex{}

	info := airwallex.GatewayInfo(ctx)

	assert.NotNil(t, info)
	assert.Equal(t, "airwallex", info.Name)
	assert.Equal(t, "Airwallex", info.DisplayName)
	assert.Equal(t, "Airwallex global payment gateway", info.Description)
	assert.Equal(t, "https://www.airwallex.com/", info.GatewayWebsiteLink)
	assert.Equal(t, "https://api.airwallex.com", info.Host)
	assert.Equal(t, "api_key", info.PublicKeyName)
	assert.Equal(t, "client_secret", info.PrivateSecretName)
	assert.Equal(t, "https://api.unibee.top/oss/file/airwallex-logo.png", info.GatewayLogo)
	assert.Equal(t, int64(230), info.Sort)
	assert.True(t, info.AutoChargeEnabled)
	assert.False(t, info.IsStaging)
}

func TestAirwallexGatewayTest(t *testing.T) {
	// This test is a placeholder for future implementation
	// In actual environment, real API keys need to be provided for testing
	t.Skip("Skipping Airwallex gateway test - requires real API keys")
}

func TestAirwallexUserCreate(t *testing.T) {
	ctx := context.Background()
	airwallex := &api.Airwallex{}

	// Mock gateway and user
	gateway := &entity.MerchantGateway{
		GatewayKey:    "test_key",
		GatewaySecret: "test_secret",
	}

	user := &entity.UserAccount{
		Id: 123,
	}

	resp, err := airwallex.GatewayUserCreate(ctx, gateway, user)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "airwallex_123", resp.GatewayUserId)
}

func TestAirwallexMerchantBalancesQuery(t *testing.T) {
	ctx := context.Background()
	airwallex := &api.Airwallex{}

	gateway := &entity.MerchantGateway{
		GatewayKey:    "test_key",
		GatewaySecret: "test_secret",
	}

	resp, err := airwallex.GatewayMerchantBalancesQuery(ctx, gateway)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.AvailableBalance)
	assert.NotNil(t, resp.PendingBalance)
}

func TestAirwallexPaymentMethods(t *testing.T) {
	ctx := context.Background()
	airwallex := &api.Airwallex{}

	gateway := &entity.MerchantGateway{
		GatewayKey:    "test_key",
		GatewaySecret: "test_secret",
	}

	// Test payment method list query
	resp, err := airwallex.GatewayUserPaymentMethodListQuery(ctx, gateway, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Test payment method attachment
	resp2, err := airwallex.GatewayUserAttachPaymentMethodQuery(ctx, gateway, 123, "method_id")
	assert.NoError(t, err)
	assert.NotNil(t, resp2)

	// Test payment method detachment
	resp3, err := airwallex.GatewayUserDeAttachPaymentMethodQuery(ctx, gateway, 123, "method_id")
	assert.NoError(t, err)
	assert.NotNil(t, resp3)

	// Test create and bind payment method
	resp4, err := airwallex.GatewayUserCreateAndBindPaymentMethod(ctx, gateway, 123, "USD", nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp4)
	assert.Equal(t, "airwallex_pm_123", resp4.PaymentMethod.Id)
}

func TestAirwallexPaymentOperations(t *testing.T) {
	ctx := context.Background()
	airwallex := &api.Airwallex{}

	gateway := &entity.MerchantGateway{
		GatewayKey:    "test_key",
		GatewaySecret: "test_secret",
	}

	// Test new payment
	createPayContext := &gateway_bean.GatewayNewPaymentReq{
		Pay: &entity.Payment{
			PaymentId: "test_payment_123",
			UserId:    123,
		},
	}

	resp, err := airwallex.GatewayNewPayment(ctx, gateway, createPayContext)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "airwallex_test_payment_123", resp.GatewayPaymentId)

	// Test payment capture
	payment := &entity.Payment{
		PaymentId: "test_payment_123",
		UserId:    123,
	}

	resp2, err := airwallex.GatewayCapture(ctx, gateway, payment)
	assert.NoError(t, err)
	assert.NotNil(t, resp2)

	// Test payment cancel
	resp3, err := airwallex.GatewayCancel(ctx, gateway, payment)
	assert.NoError(t, err)
	assert.NotNil(t, resp3)
}

func TestAirwallexRefundOperations(t *testing.T) {
	ctx := context.Background()
	airwallex := &api.Airwallex{}

	gateway := &entity.MerchantGateway{
		GatewayKey:    "test_key",
		GatewaySecret: "test_secret",
	}

	// Test refund
	createRefundContext := &gateway_bean.GatewayNewPaymentRefundReq{
		Refund: &entity.Refund{
			RefundId:      "test_refund_123",
			MerchantId:    456,
			RefundComment: "Test refund",
			RefundAmount:  1000,
		},
		Payment: &entity.Payment{
			PaymentId:        "test_payment_123",
			UserId:           123,
			MerchantId:       456,
			GatewayPaymentId: "airwallex_payment_123",
			Currency:         "USD",
		},
	}

	resp, err := airwallex.GatewayRefund(ctx, gateway, createRefundContext)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "airwallex_refund_test_refund_123", resp.GatewayRefundId)

	// Test refund cancel
	payment := &entity.Payment{
		PaymentId:  "test_payment_123",
		UserId:     123,
		MerchantId: 456,
	}
	refund := &entity.Refund{
		RefundId:     "test_refund_123",
		MerchantId:   456,
		RefundAmount: 1000,
		Currency:     "USD",
	}

	resp2, err := airwallex.GatewayRefundCancel(ctx, gateway, payment, refund)
	assert.NoError(t, err)
	assert.NotNil(t, resp2)
}

func TestAirwallexUnsupportedMethods(t *testing.T) {
	ctx := context.Background()
	airwallex := &api.Airwallex{}

	// Test unsupported crypto transaction
	_, err := airwallex.GatewayCryptoFiatTrans(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Not Supported")
}
