// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantChannelConfig is the golang structure for table merchant_channel_config.
type MerchantChannelConfig struct {
	Id               uint64      `json:"id"               description:"channel_id"`                                                // channel_id
	MerchantId       int64       `json:"merchantId"       description:"merchant_id"`                                               // merchant_id
	EnumKey          int64       `json:"enumKey"          description:"enum key , match in channel implementation"`                // enum key , match in channel implementation
	ChannelType      int         `json:"channelType"      description:"channel type，null or 0-Payment Type ｜ 1-Subscription Type"` // channel type，null or 0-Payment Type ｜ 1-Subscription Type
	Channel          string      `json:"channel"          description:"channel name"`                                              // channel name
	Name             string      `json:"name"             description:"name"`                                                      // name
	SubChannel       string      `json:"subChannel"       description:"sub_channel_enum"`                                          // sub_channel_enum
	BrandData        string      `json:"brandData"        description:""`                                                          //
	Logo             string      `json:"logo"             description:"channel logo"`                                              // channel logo
	Host             string      `json:"host"             description:"pay host"`                                                  // pay host
	ChannelAccountId string      `json:"channelAccountId" description:"channel account id"`                                        // channel account id
	ChannelKey       string      `json:"channelKey"       description:""`                                                          //
	ChannelSecret    string      `json:"channelSecret"    description:"secret"`                                                    // secret
	Custom           string      `json:"custom"           description:"custom"`                                                    // custom
	GmtCreate        *gtime.Time `json:"gmtCreate"        description:"create time"`                                               // create time
	GmtModify        *gtime.Time `json:"gmtModify"        description:"update time"`                                               // update time
	Description      string      `json:"description"      description:"description"`                                               // description
	WebhookKey       string      `json:"webhookKey"       description:"webhook_key"`                                               // webhook_key
	WebhookSecret    string      `json:"webhookSecret"    description:"webhook_secret"`                                            // webhook_secret
	UniqueProductId  string      `json:"uniqueProductId"  description:"unique  channel productId, only stripe need"`               // unique  channel productId, only stripe need
}
