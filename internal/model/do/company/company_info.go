// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CompanyInfo is the golang structure of table company_info for DAO operations like Where/Data.
type CompanyInfo struct {
	g.Meta                   `orm:"table:company_info, do:true"`
	Id                       interface{} // 公司ID
	CompanyName              interface{} // 公司名称
	CompanyLogo              interface{} // 公司logo
	TaxNumber                interface{} // 公司税号
	Status                   interface{} // 状态（1: 审核中   2： 审核通过   3： 驳回          ）
	Address                  interface{} // 公司地址 省/市/区
	DetailAddress            interface{} // 详细地址
	IsDelete                 interface{} // 逻辑删除
	GmtCreate                *gtime.Time // 新增时间
	GmtModify                *gtime.Time // 修改时间
	Name                     interface{} // 法人姓名
	Mobile                   interface{} // 联系电话
	IdCard                   interface{} // 证件号码
	IdCardFront              interface{} // 证件照片正面
	IdCardBack               interface{} // 证件照片反面
	BankcardName             interface{} // 银行卡名字（开户名）
	BankcardNum              interface{} // 银行卡账号
	BankName                 interface{} // 开户行
	SubBankName              interface{} // 开户支行名称
	BusinessUrl              interface{} // 营业执照图片
	BusinessNum              interface{} // 营业执照号码(系统默认创建设置为时间戳)
	RefuseStr                interface{} // 驳回原因
	DocumentType             interface{} // 证件类型，1,居民身份证, 2,中国护照, 3,港澳居民来往内地通行证, 4,台湾居民来往大陆通行证, 5,外国护照, 6,外国人永久居留身份证（外国人永久居留证）, 7,中华人民共和国港澳居民居住证, 8,中华人民共和国台湾居民居住证, 9,中华人民共和国外国人工作许可证（A类）, 10,中华人民共和国外国人工作许可证（B类）,11,中华人民共和国外国人工作许可证（C类）, 12,其他个人证件
	EMail                    interface{} // 邮箱地址
	ManageMobile             interface{} // 管理员手机号
	DeductPoint              interface{} // 品牌扣点（万分位，100即为扣点1%）
	CompanyFullName          interface{} // 新增公司全称、原company_name用作品牌名
	SignKey                  interface{} // 云管家signKey
	AuthToken                interface{} // 云管家authToken
	PayChannel               interface{} // 支付通道，利楚支付-0，微信支付-1，新华支付-10
	MemberName               interface{} // 会员卡名称(公司下所有会员通用)
	MemberUrl                interface{} // 会员图片URL
	MemberCopywriting        interface{} // 会员文案描述
	Wid                      interface{} // 管家店铺wid
	MerchantName             interface{} // 管家店铺名称
	Pid                      interface{} // 云管家返回信息店铺Pid
	IsBind                   interface{} // 是否绑定云管家：0.未绑定或未创建，1.在云管家新建并绑定，2.绑定云管家已有品牌
	IsBindStatus             interface{} // 是否绑定云管家：0.未绑定或未创建，1.新建或绑定成功，2.新建或绑定失败
	WxpayMchId               interface{} // 微信支付商户号配置（商家券等微信支付产品使用)
	AppSecret                interface{} // 云店开放平台密钥（已废弃）
	BalanceSaleCommission    interface{} // 品牌储值卡分账销售提成，万分位
	BalanceSaleStatus        interface{} // 品牌储值卡分账状态。0.未启用 1.已启用
	IsAppletShare            interface{} // 桌台管理小程序分享：默认为0.未设置，1.已设置
	WxMerchantLogoUrl        interface{} // 商家券商户logo
	WxBackgroundColor        interface{} // 商家券背景颜色
	WxMiniProgramsInfo       interface{} // 商家券小程序入口是否开启：默认0不开启、1开启
	WxEntranceWords          interface{} // 商家券入口文案
	WxGuidingWords           interface{} // 商家券引导文案
	WxMiniProgramsPath       interface{} // 商家券小程序路径
	WxMerchantName           interface{} // 商家券商户名称
	IsWhite                  interface{} // 是否开启白名单：默认0未开启， 1.已开启
	PassTime                 *gtime.Time // 审核通过时间,若无数据取创建时间
	BackgroundDetailUrl      interface{} // 腾讯吉市 品牌详情页背景
	MarketingCopy            interface{} // 腾讯吉市 品牌营销文案
	CompanyColor             interface{} // 腾讯吉市 品牌色
	CustomerServiceTelephone interface{} // 腾讯吉市 品牌客服电话
	ContactName              interface{} // 腾讯吉市 联系人姓名
	DistinguishCompany       interface{} // 区分腾讯吉市品牌和云店品牌 默认0，0云店，1腾讯吉市，2.云极只用于品牌注册驳回使用
	ItemCatId                interface{} // 腾讯吉市 商品品类ID
	Sort                     interface{} // 品牌排序，默认NULL
	AreaSort                 interface{} // 品牌区排序， 默认NULL, 置顶默认值为0
	WxAppId                  interface{} // 小程序AppId(吉市优惠券跳转使用)
	CompanyCat               interface{} // 腾讯吉市品牌品类ID
	VerifyStatus             interface{} // 品牌要素校验状态，0-未校验通过，1-品牌名+税号二要素校验通过
	DyMerchantId             interface{} // 抖音品牌进件商户号（抖音进件及支付使用）--废弃
	DyMerchantStatus         interface{} // 抖音品牌进件状态，0-未完成，1-已完成（--废弃）
	SharerUrl                interface{} // 小程序分享URL
	BusinessLicense          interface{} // 经营许可证URL
	TrilateralId             interface{} // 储值第三方id(品牌维度，每个品牌有一个储值第三方id)
	PreferentialPaymentId    interface{} // 优惠买单第三方id(品牌维度，每个品牌有一个优惠买单第三方id)
	CategoryQualification    interface{} // 类目资质
	BrandQualification       interface{} // 品牌资质
	City                     interface{} //
	OptionNum                interface{} // 加料数量，默认0不限制
	DyAgentStatus            interface{} // 抖音分账方进件状态：0: 未进件 1: 进件成功 2: 进件失败 3: 审核中
	Stars                    interface{} // 品牌星级,默认0,数值越大优先级越高(目前用于系统标签任务执行优先级,小于0的默认为僵尸品牌忽略处理)
	GoodsDetailSort          interface{} // 商品详情排序默认为空[]
	DineWay                  interface{} // 就餐方式，默认0店内就餐 1.打包带走 2不默认就餐方式
	DyGuestId                interface{} // 抖音商家ID(就是来客ID)
}
