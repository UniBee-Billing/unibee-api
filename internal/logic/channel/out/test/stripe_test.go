package test

import (
	"fmt"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"go-oversea-pay/utility"
	"testing"
)

func TestChangeBillingCycleAnchor(t *testing.T) {
	//go func() {
	//ctx := context.Background()
	//channelEntity := util.GetOverseaPayChannel(ctx, 25)
	//utility.Assert(channelEntity != nil, "支付渠道异常 channel not found")
	stripe.Key = "***REMOVED***"
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "unibee.server",
		Version: "0.0.1",
		URL:     "https://unibee.dev",
	})

	//detailResponse, err := sub.Get("sub_1OV191Hhgikz9ijMPTz8X9Wh", &stripe.SubscriptionParams{})
	//if err != nil {
	//	fmt.Printf("err:%s\n", err.Error())
	//}
	//fmt.Printf("detail current cycle:%d-%d\n", detailResponse.CurrentPeriodStart, detailResponse.CurrentPeriodEnd)
	//fmt.Printf("detail TrialEnd:%d\n", detailResponse.TrialEnd)

	// Cancelled Without Proration
	//params := &stripe.SubscriptionCancelParams{}
	//params.InvoiceNow = stripe.Bool(false)
	//params.Prorate = stripe.Bool(false)
	//response, err := sub.Cancel("sub_1OV191Hhgikz9ijMPTz8X9Wh", params)
	//fmt.Printf("updateResponse:%s\n", utility.MarshalToJsonString(response))
	//fmt.Printf("detail current cycle:%d-%d\n", response.CurrentPeriodStart, response.CurrentPeriodEnd)
	//fmt.Printf("detail Status:%s\n", response.Status)

	params := &stripe.CustomerListPaymentMethodsParams{
		Customer: stripe.String("cus_PJmwrgrXuesjZv"),
	}
	params.Limit = stripe.Int64(3)
	_ = customer.ListPaymentMethods(params)
	//fmt.Printf("result:%s\n", utility.MarshalToJsonString(result.List()))
	customerResult, err := customer.Get("cus_PJmwrgrXuesjZv", &stripe.CustomerParams{})
	utility.AssertError(err, "queryAndCreateChannelUser")
	//if err != nil {
	//	fmt.Printf("queryAndCreateChannelUser:%s", err.Error())
	//}
	fmt.Printf("customerResult:%s\n", utility.MarshalToJsonString(customerResult))

	//params := &stripe.InvoicePayParams{}
	//params.ChannelDefaultPaymentMethod = stripe.String("pm_1OdQUNHhgikz9ijMs0UgkN6I")
	//response, err := invoice.Pay("in_1OdziFHhgikz9ijMM0zrMlTf", params)
	//fmt.Printf("response:%s\n", utility.MarshalToJsonString(response))
	//if err != nil {
	//	fmt.Printf("err:%s\n", err.Error())
	//}

	//response, err := sub.Get("sub_1OdQUOHhgikz9ijMWA7qzh3u", &stripe.SubscriptionParams{})
	//fmt.Printf("response:%s\n", utility.MarshalToJsonString(response))
	//if err != nil {
	//	fmt.Printf("err:%s\n", err.Error())
	//}

	//updateResponse, err := sub.Update("sub_1OV191Hhgikz9ijMPTz8X9Wh", &stripe.SubscriptionParams{
	//	//TrialEnd:          stripe.Int64(1706746815),
	//	TrialEndNow:       stripe.Bool(true),
	//	ProrationBehavior: stripe.String("none"),
	//})
	//if err != nil {
	//	fmt.Printf("err:%s\n", err.Error())
	//}
	//fmt.Printf("updateResponse:%s\n", utility.MarshalToJsonString(updateResponse))
	//fmt.Printf("detail current cycle:%d-%d\n", updateResponse.CurrentPeriodStart, updateResponse.CurrentPeriodEnd)
	//fmt.Printf("detail TrialEnd:%d\n", detailResponse.TrialEnd)

	//}()

}
