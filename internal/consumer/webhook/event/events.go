package event

import (
	"sort"
)

type MerchantWebhookEvent string

const (
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CREATED   = "subscription.created"
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_UPDATED   = "subscription.updated"
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CANCELLED = "subscription.cancelled"
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_EXPIRED   = "subscription.expired"
)

var ListeningEventList = []string{
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CREATED,
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_UPDATED,
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CANCELLED,
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_EXPIRED,
}

// var ListeningEventList []string
func initListeningEventList() {
	//if ListeningEventList == nil {
	//	ListeningEventList = make([]string, 0, 100)
	//}
}

//
//func RegisterListenerEvent(i redismq.IMessageListener) {
//	if i == nil {
//		return
//	}
//	initListeningEvents()
//	if len(ListeningEventList) > 100 {
//		fmt.Println("Project Register Events Too Much ，Merge Please")
//		return
//	}
//
//	e := i.GetTag()
//	utility.Assert(!EventInListeningEvents(e), fmt.Sprintf("duplicated listener, event:%s already listened\n", e))
//
//	ListeningEventList = append(ListeningEventList, e)
//	sort.Strings(ListeningEventList)
//	fmt.Printf("Merchant_Webhook_Subscription Register Event:%s\n", e)
//}

func EventInListeningEvents(target MerchantWebhookEvent) bool {
	if len(target) <= 0 {
		return false
	}
	initListeningEventList()
	index := sort.SearchStrings(ListeningEventList, string(target))
	//index should in：[0,len(str_array)]
	if index < len(ListeningEventList) && ListeningEventList[index] == string(target) {
		return true
	}
	return false
}
