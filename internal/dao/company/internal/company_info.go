// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CompanyInfoDao is the data access object for table company_info.
type CompanyInfoDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns CompanyInfoColumns // columns contains all the column names of Table for convenient usage.
}

// CompanyInfoColumns defines and stores column names for table company_info.
type CompanyInfoColumns struct {
	Id                       string // 公司ID
	CompanyName              string // 公司名称
	CompanyLogo              string // 公司logo
	TaxNumber                string // 公司税号
	Status                   string // 状态（1: 审核中   2： 审核通过   3： 驳回          ）
	Address                  string // 公司地址 省/市/区
	DetailAddress            string // 详细地址
	IsDelete                 string // 逻辑删除
	GmtCreate                string // 新增时间
	GmtModify                string // 修改时间
	Name                     string // 法人姓名
	Mobile                   string // 联系电话
	IdCard                   string // 证件号码
	IdCardFront              string // 证件照片正面
	IdCardBack               string // 证件照片反面
	BankcardName             string // 银行卡名字（开户名）
	BankcardNum              string // 银行卡账号
	BankName                 string // 开户行
	SubBankName              string // 开户支行名称
	BusinessUrl              string // 营业执照图片
	BusinessNum              string // 营业执照号码(系统默认创建设置为时间戳)
	RefuseStr                string // 驳回原因
	DocumentType             string // 证件类型，1,居民身份证, 2,中国护照, 3,港澳居民来往内地通行证, 4,台湾居民来往大陆通行证, 5,外国护照, 6,外国人永久居留身份证（外国人永久居留证）, 7,中华人民共和国港澳居民居住证, 8,中华人民共和国台湾居民居住证, 9,中华人民共和国外国人工作许可证（A类）, 10,中华人民共和国外国人工作许可证（B类）,11,中华人民共和国外国人工作许可证（C类）, 12,其他个人证件
	EMail                    string // 邮箱地址
	ManageMobile             string // 管理员手机号
	DeductPoint              string // 品牌扣点（万分位，100即为扣点1%）
	CompanyFullName          string // 新增公司全称、原company_name用作品牌名
	SignKey                  string // 云管家signKey
	AuthToken                string // 云管家authToken
	PayChannel               string // 支付通道，利楚支付-0，微信支付-1，新华支付-10
	MemberName               string // 会员卡名称(公司下所有会员通用)
	MemberUrl                string // 会员图片URL
	MemberCopywriting        string // 会员文案描述
	Wid                      string // 管家店铺wid
	MerchantName             string // 管家店铺名称
	Pid                      string // 云管家返回信息店铺Pid
	IsBind                   string // 是否绑定云管家：0.未绑定或未创建，1.在云管家新建并绑定，2.绑定云管家已有品牌
	IsBindStatus             string // 是否绑定云管家：0.未绑定或未创建，1.新建或绑定成功，2.新建或绑定失败
	WxpayMchId               string // 微信支付商户号配置（商家券等微信支付产品使用)
	AppSecret                string // 云店开放平台密钥（已废弃）
	BalanceSaleCommission    string // 品牌储值卡分账销售提成，万分位
	BalanceSaleStatus        string // 品牌储值卡分账状态。0.未启用 1.已启用
	IsAppletShare            string // 桌台管理小程序分享：默认为0.未设置，1.已设置
	WxMerchantLogoUrl        string // 商家券商户logo
	WxBackgroundColor        string // 商家券背景颜色
	WxMiniProgramsInfo       string // 商家券小程序入口是否开启：默认0不开启、1开启
	WxEntranceWords          string // 商家券入口文案
	WxGuidingWords           string // 商家券引导文案
	WxMiniProgramsPath       string // 商家券小程序路径
	WxMerchantName           string // 商家券商户名称
	IsWhite                  string // 是否开启白名单：默认0未开启， 1.已开启
	PassTime                 string // 审核通过时间,若无数据取创建时间
	BackgroundDetailUrl      string // 腾讯吉市 品牌详情页背景
	MarketingCopy            string // 腾讯吉市 品牌营销文案
	CompanyColor             string // 腾讯吉市 品牌色
	CustomerServiceTelephone string // 腾讯吉市 品牌客服电话
	ContactName              string // 腾讯吉市 联系人姓名
	DistinguishCompany       string // 区分腾讯吉市品牌和云店品牌 默认0，0云店，1腾讯吉市，2.云极只用于品牌注册驳回使用
	ItemCatId                string // 腾讯吉市 商品品类ID
	Sort                     string // 品牌排序，默认NULL
	AreaSort                 string // 品牌区排序， 默认NULL, 置顶默认值为0
	WxAppId                  string // 小程序AppId(吉市优惠券跳转使用)
	CompanyCat               string // 腾讯吉市品牌品类ID
	VerifyStatus             string // 品牌要素校验状态，0-未校验通过，1-品牌名+税号二要素校验通过
	DyMerchantId             string // 抖音品牌进件商户号（抖音进件及支付使用）--废弃
	DyMerchantStatus         string // 抖音品牌进件状态，0-未完成，1-已完成（--废弃）
	SharerUrl                string // 小程序分享URL
	BusinessLicense          string // 经营许可证URL
	TrilateralId             string // 储值第三方id(品牌维度，每个品牌有一个储值第三方id)
	PreferentialPaymentId    string // 优惠买单第三方id(品牌维度，每个品牌有一个优惠买单第三方id)
	CategoryQualification    string // 类目资质
	BrandQualification       string // 品牌资质
	City                     string //
	OptionNum                string // 加料数量，默认0不限制
	DyAgentStatus            string // 抖音分账方进件状态：0: 未进件 1: 进件成功 2: 进件失败 3: 审核中
	Stars                    string // 品牌星级,默认0,数值越大优先级越高(目前用于系统标签任务执行优先级,小于0的默认为僵尸品牌忽略处理)
	GoodsDetailSort          string // 商品详情排序默认为空[]
	DineWay                  string // 就餐方式，默认0店内就餐 1.打包带走 2不默认就餐方式
	DyGuestId                string // 抖音商家ID(就是来客ID)
}

// companyInfoColumns holds the columns for table company_info.
var companyInfoColumns = CompanyInfoColumns{
	Id:                       "id",
	CompanyName:              "company_name",
	CompanyLogo:              "company_logo",
	TaxNumber:                "tax_number",
	Status:                   "status",
	Address:                  "address",
	DetailAddress:            "detail_address",
	IsDelete:                 "is_delete",
	GmtCreate:                "gmt_create",
	GmtModify:                "gmt_modify",
	Name:                     "name",
	Mobile:                   "mobile",
	IdCard:                   "id_card",
	IdCardFront:              "id_card_front",
	IdCardBack:               "id_card_back",
	BankcardName:             "bankcard_name",
	BankcardNum:              "bankcard_num",
	BankName:                 "bank_name",
	SubBankName:              "sub_bank_name",
	BusinessUrl:              "business_url",
	BusinessNum:              "business_num",
	RefuseStr:                "refuse_str",
	DocumentType:             "document_type",
	EMail:                    "e_mail",
	ManageMobile:             "manage_mobile",
	DeductPoint:              "deduct_point",
	CompanyFullName:          "company_full_name",
	SignKey:                  "sign_key",
	AuthToken:                "auth_token",
	PayChannel:               "pay_channel",
	MemberName:               "member_name",
	MemberUrl:                "member_url",
	MemberCopywriting:        "member_copywriting",
	Wid:                      "wid",
	MerchantName:             "merchant_name",
	Pid:                      "pid",
	IsBind:                   "is_bind",
	IsBindStatus:             "is_bind_status",
	WxpayMchId:               "wxpay_mch_id",
	AppSecret:                "app_secret",
	BalanceSaleCommission:    "balance_sale_commission",
	BalanceSaleStatus:        "balance_sale_status",
	IsAppletShare:            "is_applet_share",
	WxMerchantLogoUrl:        "wx_merchant_logo_url",
	WxBackgroundColor:        "wx_background_color",
	WxMiniProgramsInfo:       "wx_mini_programs_info",
	WxEntranceWords:          "wx_entrance_words",
	WxGuidingWords:           "wx_guiding_words",
	WxMiniProgramsPath:       "wx_mini_programs_path",
	WxMerchantName:           "wx_merchant_name",
	IsWhite:                  "is_white",
	PassTime:                 "pass_time",
	BackgroundDetailUrl:      "background_detail_url",
	MarketingCopy:            "marketing_copy",
	CompanyColor:             "company_color",
	CustomerServiceTelephone: "customer_service_telephone",
	ContactName:              "contact_name",
	DistinguishCompany:       "distinguish_company",
	ItemCatId:                "item_cat_id",
	Sort:                     "sort",
	AreaSort:                 "area_sort",
	WxAppId:                  "wx_app_id",
	CompanyCat:               "company_cat",
	VerifyStatus:             "verify_status",
	DyMerchantId:             "dy_merchant_id",
	DyMerchantStatus:         "dy_merchant_status",
	SharerUrl:                "sharer_url",
	BusinessLicense:          "business_license",
	TrilateralId:             "trilateral_id",
	PreferentialPaymentId:    "preferential_payment_id",
	CategoryQualification:    "category_qualification",
	BrandQualification:       "brand_qualification",
	City:                     "city",
	OptionNum:                "option_num",
	DyAgentStatus:            "dy_agent_status",
	Stars:                    "stars",
	GoodsDetailSort:          "goods_detail_sort",
	DineWay:                  "dine_way",
	DyGuestId:                "dy_guest_id",
}

// NewCompanyInfoDao creates and returns a new DAO object for table data access.
func NewCompanyInfoDao() *CompanyInfoDao {
	return &CompanyInfoDao{
		group:   "company",
		table:   "company_info",
		columns: companyInfoColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CompanyInfoDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CompanyInfoDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CompanyInfoDao) Columns() CompanyInfoColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CompanyInfoDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CompanyInfoDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CompanyInfoDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
