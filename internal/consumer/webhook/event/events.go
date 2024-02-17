package event

import (
	"sort"
)

const (
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CREATED   = "subscription.created"
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_UPDATED   = "subscription.updated"
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CANCELLED = "subscription.cancelled"
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_EXPIRED   = "subscription.expired"
)

var ListeningEvents = []string{
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CREATED,
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_UPDATED,
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CANCELLED,
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_EXPIRED,
}

// var ListeningEvents []string
func initListeningEvents() {
	//if ListeningEvents == nil {
	//	ListeningEvents = make([]string, 0, 100)
	//}
}

//
//func RegisterListenerEvent(i redismq.IMessageListener) {
//	if i == nil {
//		return
//	}
//	initListeningEvents()
//	if len(ListeningEvents) > 100 {
//		fmt.Println("Project Register Events Too Much ，Merge Please")
//		return
//	}
//
//	e := i.GetTag()
//	utility.Assert(!EventInListeningEvents(e), fmt.Sprintf("duplicated listener, event:%s already listened\n", e))
//
//	ListeningEvents = append(ListeningEvents, e)
//	sort.Strings(ListeningEvents)
//	fmt.Printf("Merchant_Webhook_Subscription Register Event:%s\n", e)
//}

func EventInListeningEvents(target string) bool {
	if len(target) <= 0 {
		return false
	}
	initListeningEvents()
	index := sort.SearchStrings(ListeningEvents, target)
	//index should in：[0,len(str_array)]
	if index < len(ListeningEvents) && ListeningEvents[index] == target {
		return true
	}
	return false
}
