package payment

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type PaymentsReq struct {
	g.Meta                   `path:"/payments" tags:"Open-Payment-Controller" method:"post" summary:"1.1 用于接收商户段的请求（包括 token 请求）"`
	MerchantId               int64              `p:"merchantId" dc:"商户号" v:"required|length:4,30#请输入商户号长度为:{min}到:{max}位"`
	Reference                string             `p:"reference" dc:"取消单号" v:"required"`
	Amount                   *PayAmountVo       `json:"amount"   in:"query" dc:"具体金额" v:"required"`
	PaymentMethod            *PaymentMethodsReq `json:"paymentMethod"   in:"query" dc:"支付方式" v:"required"`
	PaymentBrandAddition     *gjson.Json        `p:"paymentBrandAddition" dc:"支付方式补充(Json结构）" v:"required"`
	StorePaymentMethod       bool               `p:"storePaymentMethod" dc:"是否创建令牌" v:"required"`
	ReturnUrl                string             `p:"returnUrl" dc:"支付完成跳转地址" v:"required"`
	ShopperLocale            string             `p:"shopperLocale" dc:"语言en_US" v:"required"`
	CountryCode              string             `p:"countryCode" dc:"国家代码" v:"required"`
	TelephoneNumber          string             `p:"telephoneNumber" dc:"手机号" v:"required"`
	ShopperEmail             string             `p:"shopperEmail" dc:"用户邮箱" v:"required"`
	ShopperReference         string             `p:"shopperReference" dc:"shopper唯一Id" v:"required"`
	Channel                  string             `p:"channel" dc:"设备类型（WEB，WAP，APP, MINI, INWALLET）" v:"required"`
	LineItems                []*OutLineItem     `p:"lineItems" dc:"订单物品" v:"required"`
	DeviceType               string             `p:"deviceType" dc:"手机类型（Android、iOS）" v:"required"`
	ShopperIP                string             `p:"shopperIP" dc:"用户ip（v4，v6）" v:"required"`
	BrowserInfo              string             `p:"browserInfo" dc:"browserInfo" v:""`
	ShopperInteraction       string             `p:"shopperInteraction" dc:"交易类型" v:""`
	RecurringProcessingModel string             `p:"recurringProcessingModel" dc:"令牌类型" v:""`
	ShopperName              *OutShopperName    `p:"shopperName" dc:"shopperName" v:""`
	BillingAddress           *OutPayAddress     `p:"billingAddress" dc:"账单地址" v:""`
	DetailAddress            *OutPayAddress     `p:"detailAddress" dc:"邮寄地址" v:""`
	Capture                  bool               `p:"capture" dc:"是否立即请款" v:""`
	CaptureDelayHours        int                `p:"captureDelayHours" dc:"请款延迟执⾏时间" v:""`
	DeviceFingerprint        string             `p:"deviceFingerprint" dc:"设备指纹信息" v:""`
	MerchantOrderReference   string             `p:"merchantOrderReference" dc:"订单关联子交易码" v:""`
	Metadata                 *gjson.Json        `p:"reference" dc:"预留字段，JSON结构" v:""`
	DateOfBrith              string             `p:"dateOfBrith" dc:"生日，YYYY-MM-DD" v:""`
}
type PaymentsRes struct {
	Status    string      `p:"status" dc:"交易状态"`
	PaymentId string      `p:"paymentId" dc:"系统交易唯一编码-平台订单号"`
	Reference string      `p:"reference" dc:"商户订单号"`
	Action    *gjson.Json `p:"action" dc:"action"`
}

type OutShopperName struct {
	FirstName string `p:"firstName" dc:"名" v:"required"`
	LastName  string `p:"lastName" dc:"姓" v:"required"`
	Gender    string `p:"gender" dc:"性别" v:"required"`
}

type OutPayAddress struct {
	City              string `p:"city" dc:"城市" v:"required"`
	Country           string `p:"country" dc:"国家代码" v:"required"`
	HouseNumberOrName string `p:"houseNumberOrName" dc:"公寓名或门牌号" v:"required"`
	PostalCode        string `p:"postalCode" dc:"邮编" v:"required"`
	StateOrProvince   string `p:"stateOrProvince" dc:"洲代码" v:"required"`
	Street            string `p:"street" dc:"街道名称" v:"required"`
}

type OutLineItem struct {
	AmountExcludingTax int64  `p:"amountExcludingTax" dc:"amountExcludingTax" v:""`
	AmountIncludingTax int64  `p:"amountIncludingTax" dc:"amountIncludingTax" v:""`
	Description        string `p:"description" dc:"description" v:""`
	Id                 string `p:"id" dc:"id" v:"required"`
	Quantity           int64  `p:"quantity" dc:"quantity" v:"required"`
	TaxAmount          int64  `p:"taxAmount" dc:"taxAmount" v:"required"`
	TaxPercentage      int64  `p:"taxPercentage" dc:"taxPercentage" v:""`
	ProductUrl         string `p:"productUrl" dc:"productUrl" v:""`
	ImageUrl           string `p:"imageUrl" dc:"imageUrl" v:""`
}

type PaymentMethodsReq struct {
	g.Meta  `path:"/paymentMethods" tags:"Out-Controller" method:"post" summary:"1.0 根据配置⽀付⽅式的信息，通过请求字段筛选可以返回的⽀付⽅式(Klarna、Evonet支持）"`
	TokenId string                   `p:"tokenId" dc:"令牌id，如果有绑定过" v:""`
	Channel string                   `p:"type" dc:"支付方式类型" v:"required"`
	Issuer  []*OutPaymentMethodIssur `p:"issuer" dc:"银行信息" v:""`
}
type PaymentMethodsRes struct {
}

type OutPaymentMethodIssur struct {
	Name     string `p:"name" dc:"名称" v:""`
	Id       string `p:"id" dc:"银行对应id" v:""`
	Disabled string `p:"disabled" dc:"如该支付方式 包含银行信息 返回：当前该银 行状态" v:""`
}

type PaymentDetailsReq struct {
	g.Meta    `path:"/paymentDetails/{PaymentId}" tags:"Out-Controller" method:"post" summary:"1.5 查询当前交易状态及详情"`
	PaymentId string `in:"path" dc:"平台支付单号" v:"required|length:4,30#请输入平台支付单号长度为:{min}到:{max}位"`
}
type PaymentDetailsRes struct {
}
