// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserAccountDao is the data access object for table user_account.
type UserAccountDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns UserAccountColumns // columns contains all the column names of Table for convenient usage.
}

// UserAccountColumns defines and stores column names for table user_account.
type UserAccountColumns struct {
	Id                 string // userId
	GmtCreate          string // create time
	GmtModify          string // 修改时间
	IsDeleted          string // 0-UnDeleted，1-Deleted
	Password           string // 密码，加密存储
	UserName           string // 用户名
	Mobile             string // 手机号
	Email              string // 邮箱
	Gender             string // 性别
	AvatarUrl          string // 头像url
	ReMark             string // 备注
	IsSpecial          string // 是否是特殊账号（0.否，1.是）
	Birthday           string // 生日
	Profession         string // 职业
	School             string // 学校
	Custom             string // 其他
	NearTime           string // 最近登录时间
	Wid                string // 盟有wid
	IsRisk             string // 风控：0.低风险，1.中风险，2.高风险
	Channel            string // 渠道
	Version            string // 版本
	Phone              string //
	Address            string //
	FirstName          string //
	LastName           string //
	CompanyName        string //
	VATNumber          string //
	Telegram           string //
	WhatsAPP           string //
	WeChat             string //
	TikTok             string //
	LinkedIn           string //
	Facebook           string //
	OtherSocialInfo    string //
	PaymentMethod      string //
	CountryCode        string // country_code
	CountryName        string // country_name
	SubscriptionName   string //
	SubscriptionId     string //
	SubscriptionStatus string // sub status，0-Init | 1-Create｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete
	RecurringAmount    string //
	BillingType        string // 1-recurring,2-one-time
}

// userAccountColumns holds the columns for table user_account.
var userAccountColumns = UserAccountColumns{
	Id:                 "id",
	GmtCreate:          "gmt_create",
	GmtModify:          "gmt_modify",
	IsDeleted:          "is_deleted",
	Password:           "password",
	UserName:           "user_name",
	Mobile:             "mobile",
	Email:              "email",
	Gender:             "gender",
	AvatarUrl:          "avatar_url",
	ReMark:             "re_mark",
	IsSpecial:          "is_special",
	Birthday:           "birthday",
	Profession:         "profession",
	School:             "school",
	Custom:             "custom",
	NearTime:           "near_time",
	Wid:                "wid",
	IsRisk:             "is_risk",
	Channel:            "channel",
	Version:            "version",
	Phone:              "phone",
	Address:            "address",
	FirstName:          "first_name",
	LastName:           "last_name",
	CompanyName:        "company_name",
	VATNumber:          "VAT_number",
	Telegram:           "Telegram",
	WhatsAPP:           "WhatsAPP",
	WeChat:             "WeChat",
	TikTok:             "TikTok",
	LinkedIn:           "LinkedIn",
	Facebook:           "Facebook",
	OtherSocialInfo:    "other_social_info",
	PaymentMethod:      "payment_method",
	CountryCode:        "country_code",
	CountryName:        "country_name",
	SubscriptionName:   "subscription_name",
	SubscriptionId:     "subscription_id",
	SubscriptionStatus: "subscription_status",
	RecurringAmount:    "recurring_amount",
	BillingType:        "billing_type",
}

// NewUserAccountDao creates and returns a new DAO object for table data access.
func NewUserAccountDao() *UserAccountDao {
	return &UserAccountDao{
		group:   "oversea_pay",
		table:   "user_account",
		columns: userAccountColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UserAccountDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UserAccountDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UserAccountDao) Columns() UserAccountColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UserAccountDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UserAccountDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UserAccountDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
