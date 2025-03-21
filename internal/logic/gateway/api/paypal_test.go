package api

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"unibee/internal/cmd/config"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/gateway/api/paypal"
	"unibee/internal/logic/gateway/util"
	entity "unibee/internal/model/entity/default"
	_ "unibee/test"
	"unibee/utility"
)

func GetPaypalHost() string {
	var apiHost = "https://api-m.paypal.com"
	if !config.GetConfigInstance().IsProd() {
		apiHost = "https://api-m.sandbox.paypal.com"
	}
	return apiHost
}

func GetPaymentByGatewayPaymentId(ctx context.Context, gatewayPaymentId string) (one *entity.Payment) {
	if len(gatewayPaymentId) == 0 {
		return nil
	}
	err := dao.Payment.Ctx(ctx).Where(dao.Payment.Columns().GatewayPaymentId, gatewayPaymentId).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func TestPaypal_Gateway(t *testing.T) {
	ctx := context.Background()
	gateway := util.GetGatewayById(ctx, 47)
	c, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, GetPaypalHost())
	_, err := c.GetAccessToken(context.Background())
	utility.AssertError(err, "Test Paypal Error")

	t.Run("Test Paypal Automatic payment", func(t *testing.T) {
		tokens, err := c.GetPaymentMethodTokens(ctx, "BEEB8ANDETATED")
		if err != nil {
			t.Errorf("Not expected error for GetPaymentMethodTokens(), got %s", err.Error())
		}
		fmt.Println(utility.MarshalToJsonString(tokens))
		refund, err := c.GetRefund(ctx, "36J9867447391970W")
		if err != nil {
			t.Errorf("Not expected error for GetRefund(), got %s", err.Error())
		}
		fmt.Println(utility.MarshalToJsonString(refund))
		order, err := c.GetOrder(ctx, "2F888037H01542134")
		if err != nil {
			t.Errorf("Not expected error for GetOrder(), got %s", err.Error())
		}
		payment := GetPaymentByGatewayPaymentId(ctx, "2F888037H01542134")
		detail, _ := Paypal{}.parsePaypalPayment(ctx, gateway, order, payment)
		fmt.Println(detail)
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
				Paypal: &paypal.PaymentSourcePaypal{
					VaultId: "5a848461yc8729645",
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
		amountValue := "1.10"
		var items = make([]paypal.Item, 0)
		items = append(items, paypal.Item{
			Name:        "default product",
			Description: "default product",
			UnitAmount: &paypal.Money{
				Value:    amountValue,
				Currency: "EUR",
			},
			Quantity: "1",
		})

		orderResponse, err := c.CreateOrder(
			ctx,
			paypal.OrderIntentCapture,
			[]paypal.PurchaseUnitRequest{
				{
					Amount: &paypal.PurchaseUnitAmount{
						Value:    amountValue,
						Currency: "EUR",
					},
					SoftDescriptor: "Default Product",
				},
			},
			&paypal.CreateOrderPayer{},
			&paypal.PaymentSource{
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
		if order.PurchaseUnits[0].Amount.Value != amountValue {
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
		if order.PurchaseUnits[0].Amount.Value != amountValue {
			t.Errorf("CreateOrder amount incorrect")
		}

		var gatewayPaymentMethod string
		if order.PaymentSource != nil &&
			order.PaymentSource.Paypal != nil &&
			order.PaymentSource.Paypal.Attributes != nil &&
			order.PaymentSource.Paypal.Attributes.Vault != nil &&
			len(order.PaymentSource.Paypal.Attributes.Vault.Id) > 0 && strings.Compare(order.PaymentSource.Paypal.Attributes.Vault.Status, "VAULTED") == 0 {
			gatewayPaymentMethod = order.PaymentSource.Paypal.Attributes.Vault.Id
		}
		if gatewayPaymentMethod != "" {
			orderResponse, err = c.CreateOrder(
				ctx,
				paypal.OrderIntentCapture,
				[]paypal.PurchaseUnitRequest{
					{
						Amount: &paypal.PurchaseUnitAmount{
							Value:    amountValue,
							Currency: "USD",
						},
					},
				},
				&paypal.CreateOrderPayer{},
				&paypal.PaymentSource{
					Paypal: &paypal.PaymentSourcePaypal{
						VaultId: gatewayPaymentMethod,
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
			if order.PurchaseUnits[0].Amount.Value != amountValue {
				t.Errorf("CreateOrder amount incorrect")
			}
		}
	})
}
