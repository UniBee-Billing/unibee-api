// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantGateway is the golang structure for table merchant_gateway.
type MerchantGateway struct {
	Id               uint64      `json:"id"               description:"gateway_id"`                                                // gateway_id
	MerchantId       uint64      `json:"merchantId"       description:"merchant_id"`                                               // merchant_id
	EnumKey          int64       `json:"enumKey"          description:"enum key , match in gateway implementation"`                // enum key , match in gateway implementation
	GatewayType      int         `json:"gatewayType"      description:"gateway type，null or 0-Payment Type ｜ 1-Subscription Type"` // gateway type，null or 0-Payment Type ｜ 1-Subscription Type
	GatewayName      string      `json:"gatewayName"      description:"gateway name"`                                              // gateway name
	Name             string      `json:"name"             description:"name"`                                                      // name
	SubGateway       string      `json:"subGateway"       description:"sub_gateway_enum"`                                          // sub_gateway_enum
	BrandData        string      `json:"brandData"        description:""`                                                          //
	Logo             string      `json:"logo"             description:"gateway logo"`                                              // gateway logo
	Host             string      `json:"host"             description:"pay host"`                                                  // pay host
	GatewayAccountId string      `json:"gatewayAccountId" description:"gateway account id"`                                        // gateway account id
	GatewayKey       string      `json:"gatewayKey"       description:""`                                                          //
	GatewaySecret    string      `json:"gatewaySecret"    description:"secret"`                                                    // secret
	Custom           string      `json:"custom"           description:"custom"`                                                    // custom
	GmtCreate        *gtime.Time `json:"gmtCreate"        description:"create time"`                                               // create time
	GmtModify        *gtime.Time `json:"gmtModify"        description:"update time"`                                               // update time
	Description      string      `json:"description"      description:"description"`                                               // description
	WebhookKey       string      `json:"webhookKey"       description:"webhook_key"`                                               // webhook_key
	WebhookSecret    string      `json:"webhookSecret"    description:"webhook_secret"`                                            // webhook_secret
	UniqueProductId  string      `json:"uniqueProductId"  description:"unique  gateway productId, only stripe need"`               // unique  gateway productId, only stripe need
	CreateTime       int64       `json:"createTime"       description:"create utc time"`                                           // create utc time
}
