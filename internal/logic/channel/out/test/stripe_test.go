package test

import (
	"fmt"
	"github.com/stripe/stripe-go/v76"
	sub "github.com/stripe/stripe-go/v76/subscription"
	"testing"
)

func TestChangeBillingCycleAnchor(t *testing.T) {
	//go func() {
	//ctx := context.Background()
	//channelEntity := util.GetOverseaPayChannel(ctx, 25)
	//utility.Assert(channelEntity != nil, "支付渠道异常 channel not found")
	stripe.Key = "sk_test_51ONBbPHhgikz9ijMYqMBiSjihBxHApWQmt9s9mtd8vgA4O9PupcHHNVQWc4Dd7JYZ2xRrg9quSrp0g1XfO64Xkbq00zOkBRIoz"
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "unibee.server",
		Version: "0.0.1",
		URL:     "https://unibee.dev",
	})

	detailResponse, err := sub.Get("sub_1OV191Hhgikz9ijMPTz8X9Wh", &stripe.SubscriptionParams{})
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	fmt.Printf("detail current cycle:%d-%d\n", detailResponse.CurrentPeriodStart, detailResponse.CurrentPeriodEnd)
	fmt.Printf("detail TrialEnd:%d\n", detailResponse.TrialEnd)

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
