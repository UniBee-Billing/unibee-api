// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantGateway is the golang structure of table merchant_gateway for DAO operations like Where/Data.
type MerchantGateway struct {
	g.Meta           `orm:"table:merchant_gateway, do:true"`
	Id               interface{} // gateway_id
	MerchantId       interface{} // merchant_id
	EnumKey          interface{} // enum key , match in gateway implementation
	GatewayType      interface{} // gateway type，null or 0-Payment Type ｜ 1-Subscription Type
	GatewayName      interface{} // gateway name
	Name             interface{} // name
	SubGateway       interface{} // sub_gateway_enum
	BrandData        interface{} //
	Logo             interface{} // gateway logo
	Host             interface{} // pay host
	GatewayAccountId interface{} // gateway account id
	GatewayKey       interface{} //
	GatewaySecret    interface{} // secret
	Custom           interface{} // custom
	GmtCreate        *gtime.Time // create time
	GmtModify        *gtime.Time // update time
	Description      interface{} // description
	WebhookKey       interface{} // webhook_key
	WebhookSecret    interface{} // webhook_secret
	UniqueProductId  interface{} // unique  gateway productId, only stripe need
	CreateTime       interface{} // create utc time
}
