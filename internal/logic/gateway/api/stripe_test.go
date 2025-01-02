package api

import (
	"context"
	"fmt"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/paymentmethod"
	"github.com/stripe/stripe-go/v78/refund"
	"strings"
	"unibee/utility"

	"github.com/gogf/gf/v2/test/gtest"
	"testing"
	"unibee/internal/query"
	_test "unibee/test"
)

func init() {

}

func setUniBeeAppInfo() {
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "unibee.server",
		Version: "0.0.1",
		URL:     "https://unibee.dev",
	})
}

func TestCheckout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		gateway := query.GetGatewayById(ctx, 25)
		_test.AssertNotNil(gateway)
		stripe.Key = gateway.GatewaySecret
		setUniBeeAppInfo()
		{
			var items []*stripe.CheckoutSessionLineItemParams
			items = append(items, &stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(strings.ToLower("EUR")),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(fmt.Sprintf("%s", "Test Checkout")),
					},
					UnitAmount: stripe.Int64(200),
				},
				Quantity: stripe.Int64(1),
			})

			checkoutParams := &stripe.CheckoutSessionParams{
				Customer:  stripe.String("cus_Q53EmPEk3hxJF9"),
				Currency:  stripe.String(strings.ToLower("EUR")),
				LineItems: items,
				PaymentMethodTypes: stripe.StringSlice([]string{
					//"card",
					//"link",
					"au_becs_debit",
				}),
				PaymentMethodData: &stripe.CheckoutSessionPaymentMethodDataParams{AllowRedisplay: stripe.String(string(stripe.PaymentMethodAllowRedisplayAlways))},
				SuccessURL:        stripe.String("http://merchant.unibee.top"),
				CancelURL:         stripe.String("http://merchant.unibee.top"),
				PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
					SetupFutureUsage: stripe.String(string(stripe.PaymentIntentSetupFutureUsageOffSession)),
				},
			}
			checkoutParams.Mode = stripe.String(string(stripe.CheckoutSessionModePayment))
			detail, err := session.New(checkoutParams)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(utility.MarshalToJsonString(detail))
			}
		}
	})
}

func TestStripe(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		gateway := query.GetGatewayById(ctx, 25)
		_test.AssertNotNil(gateway)
		stripe.Key = gateway.GatewaySecret
		setUniBeeAppInfo()
		gatewayUser := QueryAndCreateGatewayUser(ctx, gateway, 2235427988)

		{
			params := &stripe.CustomerListPaymentMethodsParams{
				Customer: stripe.String(gatewayUser.GatewayUserId),
			}
			params.Limit = stripe.Int64(10)
			result := customer.ListPaymentMethods(params)
			fmt.Println(utility.MarshalToJsonString(result))
		}
		{
			params := &stripe.CustomerRetrievePaymentMethodParams{
				Customer: stripe.String(gatewayUser.GatewayUserId),
			}
			result, err := customer.RetrievePaymentMethod("pm_1OmpHYHhgikz9ijMWYrNNhs5", params)
			if err != nil {
				fmt.Println(utility.MarshalToJsonString(err))
			}
			fmt.Println(utility.MarshalToJsonString(result))
		}

	})

	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		gateway := query.GetGatewayById(ctx, 25)
		_test.AssertNotNil(gateway)
		stripe.Key = gateway.GatewaySecret
		setUniBeeAppInfo()
		gatewayUser := QueryAndCreateGatewayUser(ctx, gateway, 2235427988)
		{
			params := &stripe.PaymentMethodListParams{
				Type:     stripe.String(string(stripe.PaymentMethodTypeCard)),
				Customer: stripe.String(gatewayUser.GatewayUserId),
			}
			params.Limit = stripe.Int64(3)
			result := paymentmethod.List(params)
			fmt.Println(utility.MarshalToJsonString(result))
		}
		{
			params := &stripe.PaymentMethodParams{}
			result, err := paymentmethod.Get("pm_1OmpHYHhgikz9ijMWYrNNhs5", params)
			if err != nil {
				fmt.Println(utility.MarshalToJsonString(err))
			}
			fmt.Println(utility.MarshalToJsonString(result))
		}
		{
			params := &stripe.PaymentMethodAttachParams{
				Customer: stripe.String(gatewayUser.GatewayUserId),
			}
			result, err := paymentmethod.Attach(gatewayUser.GatewayDefaultPaymentMethod, params)
			if err != nil {
				fmt.Println(utility.MarshalToJsonString(err))
			}
			fmt.Println(utility.MarshalToJsonString(result))
		}
	})
}

func TestStripeQueryAllRefunds(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		stripe.Key = ""
		setUniBeeAppInfo()
		response := refund.List(&stripe.RefundListParams{
			PaymentIntent: stripe.String("pi_3QWhDnDaLWZKMs9N2ecEfXU8"),
		})
		fmt.Println(utility.MarshalToJsonString(response))
	})
}
