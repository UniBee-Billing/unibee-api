// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MerchantInfoDao is the data access object for table merchant_info.
type MerchantInfoDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns MerchantInfoColumns // columns contains all the column names of Table for convenient usage.
}

// MerchantInfoColumns defines and stores column names for table merchant_info.
type MerchantInfoColumns struct {
	Id                     string // 用户的ID
	CompanyId              string // 公司ID
	UserId                 string // 用户ID
	Type                   string // 类型，0-个人，1-企业
	CompanyName            string // 公司名称
	BusinessNum            string // 税号
	Name                   string // 个人或法人姓名
	Idcard                 string // 个人或法人身份证号
	Location               string // 省市区地址
	Address                string // 详细地址
	IdcardFrontPic         string // 个人或法人身份证正面
	IdcardBackPic          string // 个人或法人身份证背面
	BusinessLicensePic     string // 营业执照图片
	Tag                    string // 标签，或业务类型，比如trtc,im
	GmtCreate              string // 创建时间
	GmtModify              string // 修改时间
	BecomePractitionerTime string // 成为练习者时间
	IsDeleted              string // 是否删除，0-未删除，1-删除
	CompanyLogo            string // 账号头像
	Mobile                 string // 登录手机号
	Mark                   string // 备注
	SettleCurrency         string // 结算币种，null 代表结算用 CNY
	ServiceRate            string // 服务费比例，万分位，百分比[0，10000)，精度为0.01%，如3即为0.03%
}

// merchantInfoColumns holds the columns for table merchant_info.
var merchantInfoColumns = MerchantInfoColumns{
	Id:                     "id",
	CompanyId:              "company_id",
	UserId:                 "user_id",
	Type:                   "type",
	CompanyName:            "company_name",
	BusinessNum:            "business_num",
	Name:                   "name",
	Idcard:                 "idcard",
	Location:               "location",
	Address:                "address",
	IdcardFrontPic:         "idcard_front_pic",
	IdcardBackPic:          "idcard_back_pic",
	BusinessLicensePic:     "business_license_pic",
	Tag:                    "tag",
	GmtCreate:              "gmt_create",
	GmtModify:              "gmt_modify",
	BecomePractitionerTime: "become_practitioner_time",
	IsDeleted:              "is_deleted",
	CompanyLogo:            "company_logo",
	Mobile:                 "mobile",
	Mark:                   "mark",
	SettleCurrency:         "settle_currency",
	ServiceRate:            "service_rate",
}

// NewMerchantInfoDao creates and returns a new DAO object for table data access.
func NewMerchantInfoDao() *MerchantInfoDao {
	return &MerchantInfoDao{
		group:   "oversea_pay",
		table:   "merchant_info",
		columns: merchantInfoColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MerchantInfoDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MerchantInfoDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MerchantInfoDao) Columns() MerchantInfoColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MerchantInfoDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MerchantInfoDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MerchantInfoDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
