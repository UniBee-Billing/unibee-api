// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantChannelConfig is the golang structure of table merchant_channel_config for DAO operations like Where/Data.
type MerchantChannelConfig struct {
	g.Meta           `orm:"table:merchant_channel_config, do:true"`
	Id               interface{} // channel_id
	MerchantId       interface{} // merchant_id
	EnumKey          interface{} // enum key , match in channel implementation
	ChannelType      interface{} // channel type，null or 0-Payment Type ｜ 1-Subscription Type
	Channel          interface{} // channel name
	Name             interface{} // name
	SubChannel       interface{} // sub_channel_enum
	BrandData        interface{} //
	Logo             interface{} // channel logo
	Host             interface{} // pay host
	ChannelAccountId interface{} // channel account id
	ChannelKey       interface{} //
	ChannelSecret    interface{} // secret
	Custom           interface{} // custom
	GmtCreate        *gtime.Time // create time
	GmtModify        *gtime.Time // update time
	Description      interface{} // description
	WebhookKey       interface{} // webhook_key
	WebhookSecret    interface{} // webhook_secret
	UniqueProductId  interface{} // unique  channel productId, only stripe need
}
