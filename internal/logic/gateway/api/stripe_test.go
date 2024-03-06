package api

import (
	"context"
	"fmt"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/paymentmethod"
	"unibee/utility"

	"github.com/gogf/gf/v2/test/gtest"
	"testing"
	"unibee/internal/query"
	_test "unibee/test"
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
		//{
		//	paymentIntentDetail, err := GetGatewayServiceProvider(ctx, gateway.Id).GatewayPaymentDetail(ctx, gateway, "pi_3OmpHZHhgikz9ijM0a87ACwq")
		//	if err != nil {
		//		fmt.Println(utility.MarshalToJsonString(err))
		//	}
		//	fmt.Println(utility.MarshalToJsonString(paymentIntentDetail))
		//}
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
		//{
		//	params := &stripe.PaymentIntentParams{
		//		Customer: stripe.String(gatewayUser.GatewayUserId),
		//		Confirm:  stripe.Bool(true),
		//		Amount:   stripe.Int64(101),
		//		Currency: stripe.String(strings.ToLower("USD")),
		//		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
		//			Enabled: stripe.Bool(true),
		//		},
		//		//Metadata:  createPayContext.MetaData,
		//		ReturnURL:        stripe.String("http://user.unibee.top"),
		//		SetupFutureUsage: stripe.String(string(stripe.PaymentIntentSetupFutureUsageOffSession)),
		//	}
		//	params.PaymentMethod = stripe.String(gatewayUser.GatewayDefaultPaymentMethod)
		//	targetIntent, err := paymentintent.New(params)
		//	if err != nil {
		//		fmt.Println(utility.MarshalToJsonString(err))
		//	}
		//	fmt.Println(utility.MarshalToJsonString(targetIntent))
		//}
	})
}
