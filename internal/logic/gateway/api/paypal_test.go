package api

import (
	"context"
	"fmt"
	"testing"
	"unibee/internal/logic/gateway/api/paypal"
	"unibee/internal/query"
	_ "unibee/test"
	"unibee/utility"
)

func TestPaypal_Gateway(t *testing.T) {
	ctx := context.Background()
	gateway := query.GetGatewayById(ctx, 45)
	c, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, gateway.Host)
	_, err := c.GetAccessToken(context.Background())
	utility.AssertError(err, "Test Paypal Error")

	t.Run("Test Paypal Automatic payment", func(t *testing.T) {
		order, err := c.GetOrder(ctx, "9CH16636GU041544S")
		if err != nil {
			t.Errorf("Not expected error for GetOrder(), got %s", err.Error())
		}
		orderResponse, err := c.CreateOrder(
			ctx,
			paypal.OrderIntentCapture,
			[]paypal.PurchaseUnitRequest{
				{
					Amount: &paypal.PurchaseUnitAmount{
						Value:    "1.00",
						Currency: "USD",
					},
				},
			},
			&paypal.CreateOrderPayer{},
			&paypal.PaymentSource{
				//Card: &paypal.PaymentSourceCard{Attributes: &paypal.PaymentSourceAttributes{
				//	Vault: paypal.PaymentSourceAttributesVault{
				//		StoreInVault: "ON_SUCCESS",
				//	},
				//	Verification: paypal.PaymentSourceAttributesVerification{Method: "SCA_WHEN_REQUIRED"},
				//}},
				Paypal: &paypal.PaymentSourcePaypal{
					VaultId: "5a848461yc8729645",
					//Attributes: &paypal.PaymentSourceAttributes{
					//	Vault: &paypal.PaymentSourceAttributesVault{
					//		StoreInVault: "ON_SUCCESS",
					//		UsageType:    "MERCHANT",
					//	},
					//},
				},
			},
			&paypal.ApplicationContext{
				BrandName:          "",
				Locale:             "",
				ShippingPreference: "",
				UserAction:         "",
				PaymentMethod:      paypal.PaymentMethod{},
				ReturnURL:          "https://merchant.unibee.top",
				CancelURL:          "https://user.unibee.top",
			},
			utility.CreatePaymentId(),
		)
		if err != nil {
			t.Errorf("Not expected error for CreateOrder(), got %s", err.Error())
		}
		order, err = c.GetOrder(ctx, orderResponse.ID)
		if err != nil {
			t.Errorf("Not expected error for GetOrder(), got %s", err.Error())
		}
		if order.PurchaseUnits[0].Amount.Value != "1.00" {
			t.Errorf("CreateOrder amount incorrect")
		}
	})

	t.Run("Test Paypal Checkout payment", func(t *testing.T) {
		orderResponse, err := c.CreateOrder(
			ctx,
			paypal.OrderIntentCapture,
			[]paypal.PurchaseUnitRequest{
				{
					Amount: &paypal.PurchaseUnitAmount{
						Value:    "1.00",
						Currency: "USD",
					},
				},
			},
			&paypal.CreateOrderPayer{},
			&paypal.PaymentSource{
				//Card: &paypal.PaymentSourceCard{Attributes: &paypal.PaymentSourceAttributes{
				//	Vault: paypal.PaymentSourceAttributesVault{
				//		StoreInVault: "ON_SUCCESS",
				//	},
				//	Verification: paypal.PaymentSourceAttributesVerification{Method: "SCA_WHEN_REQUIRED"},
				//}},
				Paypal: &paypal.PaymentSourcePaypal{
					Attributes: &paypal.PaymentSourceAttributes{
						Vault: &paypal.PaymentSourceAttributesVault{
							StoreInVault: "ON_SUCCESS",
							UsageType:    "MERCHANT",
						},
					},
				},
			},
			&paypal.ApplicationContext{
				BrandName:          "",
				Locale:             "",
				ShippingPreference: "",
				UserAction:         "",
				PaymentMethod:      paypal.PaymentMethod{},
				ReturnURL:          "https://merchant.unibee.top",
				CancelURL:          "https://user.unibee.top",
			},
			utility.CreatePaymentId(),
		)
		if err != nil {
			t.Errorf("Not expected error for CreateOrder(), got %s", err.Error())
		}
		order, err := c.GetOrder(ctx, orderResponse.ID)
		if err != nil {
			t.Errorf("Not expected error for GetOrder(), got %s", err.Error())
		}
		if order.PurchaseUnits[0].Amount.Value != "1.00" {
			t.Errorf("CreateOrder amount incorrect")
		}

		captureOrder, err := c.CaptureOrder(ctx, order.ID, paypal.CaptureOrderRequest{})
		if err != nil {
			t.Errorf("Not expected error for CaptureOrder(), got %s", err.Error())
		}
		fmt.Println(utility.MarshalToJsonString(captureOrder))

		order, err = c.GetOrder(ctx, orderResponse.ID)
		if err != nil {
			t.Errorf("Not expected error for GetOrder(), got %s", err.Error())
		}
		if order.PurchaseUnits[0].Amount.Value != "1.00" {
			t.Errorf("CreateOrder amount incorrect")
		}

	})
}
