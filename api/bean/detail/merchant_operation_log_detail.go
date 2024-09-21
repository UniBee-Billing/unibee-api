package detail

type MerchantOperationLogDetail struct {
	Id             uint64                `json:"id"                 description:"id"`                                                   // id
	MerchantId     uint64                `json:"merchantId"         description:"merchant Id"`                                          // merchant Id
	MemberId       uint64                `json:"memberId"           description:"member_id"`                                            // member_id
	OptAccount     string                `json:"optAccount"          description:"operation Account"`                                   // operation target
	OptAccountId   string                `json:"optAccountId"       description:"operation account id"`                                 // operation account id
	OptAccountType int                   `json:"optAccountType"     description:"opt_account_type, 0-Member|1-User|2-OpenApi|3-System"` // opt_account_type, 0-Member|1-User|2-OpenApi|3-System
	OptTarget      string                `json:"optTarget"          description:"operation target"`                                     // operation target
	OptContent     string                `json:"optContent"         description:"operation content"`                                    // operation content
	CreateTime     int64                 `json:"createTime"         description:"operation create utc time"`                            // operation create utc time
	SubscriptionId string                `json:"subscriptionId"     description:"subscription_id"`                                      // subscription_id
	UserId         uint64                `json:"userId"             description:"user_id"`                                              // user_id
	InvoiceId      string                `json:"invoiceId"          description:"invoice id"`                                           // invoice id
	PlanId         uint64                `json:"planId"             description:"plan id"`                                              // plan id
	DiscountCode   string                `json:"discountCode"       description:"discount_code"`                                        // discount_code
	Member         *MerchantMemberDetail `json:"member"             description:"Member"`
	//User           *bean.UserAccountSimplify  `json:"user"               description:"User"`
	//Subscription   *bean.SubscriptionSimplify `json:"subscription"       description:"Subscription"`
}
