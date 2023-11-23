// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CompanyInfo is the golang structure for table company_info.
type CompanyInfo struct {
	Id                       uint64      `json:"id"                       ` // 公司ID
	CompanyName              string      `json:"companyName"              ` // 公司名称
	CompanyLogo              string      `json:"companyLogo"              ` // 公司logo
	TaxNumber                string      `json:"taxNumber"                ` // 公司税号
	Status                   int         `json:"status"                   ` // 状态（1: 审核中   2： 审核通过   3： 驳回          ）
	Address                  string      `json:"address"                  ` // 公司地址 省/市/区
	DetailAddress            string      `json:"detailAddress"            ` // 详细地址
	IsDelete                 int         `json:"isDelete"                 ` // 逻辑删除
	GmtCreate                *gtime.Time `json:"gmtCreate"                ` // 新增时间
	GmtModify                *gtime.Time `json:"gmtModify"                ` // 修改时间
	Name                     string      `json:"name"                     ` // 法人姓名
	Mobile                   string      `json:"mobile"                   ` // 联系电话
	IdCard                   string      `json:"idCard"                   ` // 证件号码
	IdCardFront              string      `json:"idCardFront"              ` // 证件照片正面
	IdCardBack               string      `json:"idCardBack"               ` // 证件照片反面
	BankcardName             string      `json:"bankcardName"             ` // 银行卡名字（开户名）
	BankcardNum              string      `json:"bankcardNum"              ` // 银行卡账号
	BankName                 string      `json:"bankName"                 ` // 开户行
	SubBankName              string      `json:"subBankName"              ` // 开户支行名称
	BusinessUrl              string      `json:"businessUrl"              ` // 营业执照图片
	BusinessNum              string      `json:"businessNum"              ` // 营业执照号码(系统默认创建设置为时间戳)
	RefuseStr                string      `json:"refuseStr"                ` // 驳回原因
	DocumentType             int         `json:"documentType"             ` // 证件类型，1,居民身份证, 2,中国护照, 3,港澳居民来往内地通行证, 4,台湾居民来往大陆通行证, 5,外国护照, 6,外国人永久居留身份证（外国人永久居留证）, 7,中华人民共和国港澳居民居住证, 8,中华人民共和国台湾居民居住证, 9,中华人民共和国外国人工作许可证（A类）, 10,中华人民共和国外国人工作许可证（B类）,11,中华人民共和国外国人工作许可证（C类）, 12,其他个人证件
	EMail                    string      `json:"eMail"                    ` // 邮箱地址
	ManageMobile             string      `json:"manageMobile"             ` // 管理员手机号
	DeductPoint              int         `json:"deductPoint"              ` // 品牌扣点（万分位，100即为扣点1%）
	CompanyFullName          string      `json:"companyFullName"          ` // 新增公司全称、原company_name用作品牌名
	SignKey                  string      `json:"signKey"                  ` // 云管家signKey
	AuthToken                string      `json:"authToken"                ` // 云管家authToken
	PayChannel               int         `json:"payChannel"               ` // 支付通道，利楚支付-0，微信支付-1，新华支付-10
	MemberName               string      `json:"memberName"               ` // 会员卡名称(公司下所有会员通用)
	MemberUrl                string      `json:"memberUrl"                ` // 会员图片URL
	MemberCopywriting        string      `json:"memberCopywriting"        ` // 会员文案描述
	Wid                      int64       `json:"wid"                      ` // 管家店铺wid
	MerchantName             string      `json:"merchantName"             ` // 管家店铺名称
	Pid                      int64       `json:"pid"                      ` // 云管家返回信息店铺Pid
	IsBind                   int         `json:"isBind"                   ` // 是否绑定云管家：0.未绑定或未创建，1.在云管家新建并绑定，2.绑定云管家已有品牌
	IsBindStatus             int         `json:"isBindStatus"             ` // 是否绑定云管家：0.未绑定或未创建，1.新建或绑定成功，2.新建或绑定失败
	WxpayMchId               string      `json:"wxpayMchId"               ` // 微信支付商户号配置（商家券等微信支付产品使用)
	AppSecret                string      `json:"appSecret"                ` // 云店开放平台密钥（已废弃）
	BalanceSaleCommission    int         `json:"balanceSaleCommission"    ` // 品牌储值卡分账销售提成，万分位
	BalanceSaleStatus        int         `json:"balanceSaleStatus"        ` // 品牌储值卡分账状态。0.未启用 1.已启用
	IsAppletShare            int         `json:"isAppletShare"            ` // 桌台管理小程序分享：默认为0.未设置，1.已设置
	WxMerchantLogoUrl        string      `json:"wxMerchantLogoUrl"        ` // 商家券商户logo
	WxBackgroundColor        string      `json:"wxBackgroundColor"        ` // 商家券背景颜色
	WxMiniProgramsInfo       int         `json:"wxMiniProgramsInfo"       ` // 商家券小程序入口是否开启：默认0不开启、1开启
	WxEntranceWords          string      `json:"wxEntranceWords"          ` // 商家券入口文案
	WxGuidingWords           string      `json:"wxGuidingWords"           ` // 商家券引导文案
	WxMiniProgramsPath       string      `json:"wxMiniProgramsPath"       ` // 商家券小程序路径
	WxMerchantName           string      `json:"wxMerchantName"           ` // 商家券商户名称
	IsWhite                  int         `json:"isWhite"                  ` // 是否开启白名单：默认0未开启， 1.已开启
	PassTime                 *gtime.Time `json:"passTime"                 ` // 审核通过时间,若无数据取创建时间
	BackgroundDetailUrl      string      `json:"backgroundDetailUrl"      ` // 腾讯吉市 品牌详情页背景
	MarketingCopy            string      `json:"marketingCopy"            ` // 腾讯吉市 品牌营销文案
	CompanyColor             string      `json:"companyColor"             ` // 腾讯吉市 品牌色
	CustomerServiceTelephone string      `json:"customerServiceTelephone" ` // 腾讯吉市 品牌客服电话
	ContactName              string      `json:"contactName"              ` // 腾讯吉市 联系人姓名
	DistinguishCompany       int         `json:"distinguishCompany"       ` // 区分腾讯吉市品牌和云店品牌 默认0，0云店，1腾讯吉市，2.云极只用于品牌注册驳回使用
	ItemCatId                int64       `json:"itemCatId"                ` // 腾讯吉市 商品品类ID
	Sort                     int         `json:"sort"                     ` // 品牌排序，默认NULL
	AreaSort                 int         `json:"areaSort"                 ` // 品牌区排序， 默认NULL, 置顶默认值为0
	WxAppId                  string      `json:"wxAppId"                  ` // 小程序AppId(吉市优惠券跳转使用)
	CompanyCat               int64       `json:"companyCat"               ` // 腾讯吉市品牌品类ID
	VerifyStatus             int         `json:"verifyStatus"             ` // 品牌要素校验状态，0-未校验通过，1-品牌名+税号二要素校验通过
	DyMerchantId             string      `json:"dyMerchantId"             ` // 抖音品牌进件商户号（抖音进件及支付使用）--废弃
	DyMerchantStatus         int         `json:"dyMerchantStatus"         ` // 抖音品牌进件状态，0-未完成，1-已完成（--废弃）
	SharerUrl                string      `json:"sharerUrl"                ` // 小程序分享URL
	BusinessLicense          string      `json:"businessLicense"          ` // 经营许可证URL
	TrilateralId             string      `json:"trilateralId"             ` // 储值第三方id(品牌维度，每个品牌有一个储值第三方id)
	PreferentialPaymentId    string      `json:"preferentialPaymentId"    ` // 优惠买单第三方id(品牌维度，每个品牌有一个优惠买单第三方id)
	CategoryQualification    string      `json:"categoryQualification"    ` // 类目资质
	BrandQualification       string      `json:"brandQualification"       ` // 品牌资质
	City                     string      `json:"city"                     ` //
	OptionNum                int         `json:"optionNum"                ` // 加料数量，默认0不限制
	DyAgentStatus            int         `json:"dyAgentStatus"            ` // 抖音分账方进件状态：0: 未进件 1: 进件成功 2: 进件失败 3: 审核中
	Stars                    int         `json:"stars"                    ` // 品牌星级,默认0,数值越大优先级越高(目前用于系统标签任务执行优先级,小于0的默认为僵尸品牌忽略处理)
	GoodsDetailSort          string      `json:"goodsDetailSort"          ` // 商品详情排序默认为空[]
	DineWay                  int         `json:"dineWay"                  ` // 就餐方式，默认0店内就餐 1.打包带走 2不默认就餐方式
	DyGuestId                string      `json:"dyGuestId"                ` // 抖音商家ID(就是来客ID)
}
