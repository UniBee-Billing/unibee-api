package api

import (
	"context"
	"fmt"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/paymentmethod"
	"unibee-api/utility"

	"github.com/gogf/gf/v2/test/gtest"
	"testing"
	"unibee-api/internal/query"
	_test "unibee-api/test"
)

func setUnibeeAppInfo() {
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "unibee.server",
		Version: "0.0.1",
		URL:     "https://unibee.dev",
	})
}

func TestStrip(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		gateway := query.GetGatewayById(ctx, 25)
		_test.AssertNotNil(gateway)
		stripe.Key = gateway.GatewaySecret
		setUnibeeAppInfo()
		gatewayUser := queryAndCreateChannelUser(ctx, gateway, 2235427988)

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
		setUnibeeAppInfo()
		gatewayUser := queryAndCreateChannelUser(ctx, gateway, 2235427988)
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
			paymentIntentDetail, err := GetGatewayServiceProvider(ctx, int64(gateway.Id)).GatewayPaymentDetail(ctx, gateway, "pi_3OmpHZHhgikz9ijM0a87ACwq")
			if err != nil {
				fmt.Println(utility.MarshalToJsonString(err))
			}
			fmt.Println(utility.MarshalToJsonString(paymentIntentDetail))
		}
	})
}
