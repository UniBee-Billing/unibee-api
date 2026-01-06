package test

import (
	"context"
	"testing"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/gateway_bean"
	entity "unibee/internal/model/entity/default"
)

func TestBlockonomicsGateway(t *testing.T) {
	ctx := context.Background()

	// Test gateway info
	gateway := &api.Blockonomics{}
	info := gateway.GatewayInfo(ctx)

	if info == nil {
		t.Fatal("GatewayInfo should not be nil")
	}

	if info.Name != "Blockonomics" {
		t.Errorf("Expected gateway name to be 'Blockonomics', got '%s'", info.Name)
	}

	if info.DisplayName != "Bitcoin" {
		t.Errorf("Expected display name to be 'Bitcoin', got '%s'", info.DisplayName)
	}

	if info.GatewayType != 2 { // 2 = Crypto
		t.Errorf("Expected gateway type to be 2 (Crypto), got %d", info.GatewayType)
	}

	if info.Sort != 85 {
		t.Errorf("Expected sort to be 85, got %d", info.Sort)
	}

	if info.AutoChargeEnabled {
		t.Error("Expected AutoChargeEnabled to be false")
	}

	t.Logf("Gateway info: %+v", info)
}

func TestBlockonomicsGatewayTest(t *testing.T) {
	ctx := context.Background()

	gateway := &api.Blockonomics{}

	// Test with API key (Blockonomics currencies endpoint may be public)
	req := &_interface.GatewayTestReq{
		Key: "",
	}

	icon, gatewayType, err := gateway.GatewayTest(ctx, req)
	if err != nil {
		t.Fatalf("Gateway test failed: %v", err)
	}

	if icon == "" {
		t.Error("Expected icon URL to be returned")
	}

	if gatewayType != 2 { // 2 = Crypto
		t.Errorf("Expected gateway type to be 2 (Crypto), got %d", gatewayType)
	}

	t.Logf("Gateway test successful: icon=%s, type=%d", icon, gatewayType)
}

func TestBlockonomicsUserCreate(t *testing.T) {
	ctx := context.Background()

	gateway := &api.Blockonomics{}

	// Mock user account
	user := &entity.UserAccount{
		Id: 12345,
	}

	resp, err := gateway.GatewayUserCreate(ctx, nil, user)
	if err != nil {
		t.Fatalf("GatewayUserCreate failed: %v", err)
	}

	if resp.GatewayUserId != "12345" {
		t.Errorf("Expected GatewayUserId to be '12345', got '%s'", resp.GatewayUserId)
	}

	t.Logf("User create response: %+v", resp)
}

func TestBlockonomicsPaymentMethodList(t *testing.T) {
	ctx := context.Background()

	gateway := &api.Blockonomics{}

	req := &gateway_bean.GatewayUserPaymentMethodReq{
		UserId: 12345,
	}

	resp, err := gateway.GatewayUserPaymentMethodListQuery(ctx, nil, req)
	if err != nil {
		t.Fatalf("GatewayUserPaymentMethodListQuery failed: %v", err)
	}

	if len(resp.PaymentMethods) != 1 {
		t.Errorf("Expected 1 payment method, got %d", len(resp.PaymentMethods))
	}

	if resp.PaymentMethods[0].Id != "BTC" {
		t.Errorf("Expected payment method ID to be 'BTC', got '%s'", resp.PaymentMethods[0].Id)
	}

	t.Logf("Payment method list response: %+v", resp)
}
