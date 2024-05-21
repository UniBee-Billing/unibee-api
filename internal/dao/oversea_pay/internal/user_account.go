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
	GatewayId          string // gateway_id
	PaymentMethod      string //
	MerchantId         string // merchant_id
	GmtCreate          string // create time
	GmtModify          string // update time
	IsDeleted          string // 0-UnDeleted，1-Deleted
	Password           string // password , encrypt
	UserName           string // user name
	Mobile             string // mobile
	Email              string // email
	Gender             string // gender
	AvatarUrl          string // avator url
	ReMark             string // note
	IsSpecial          string // is special account（0.no，1.yes）- deperated
	Birthday           string // brithday
	Profession         string // profession
	School             string // school
	Custom             string // custom
	LastLoginAt        string // last login time, utc time
	IsRisk             string // is risk account (deperated)
	Version            string // version
	Phone              string // phone
	Address            string // address
	FirstName          string // first name
	LastName           string // last name
	CompanyName        string // company name
	VATNumber          string // vat number
	Telegram           string // telegram
	WhatsAPP           string // whats app
	WeChat             string // wechat
	TikTok             string // tictok
	LinkedIn           string // linkedin
	Facebook           string // facebook
	OtherSocialInfo    string //
	CountryCode        string // country_code
	CountryName        string // country_name
	SubscriptionName   string // subscription name
	SubscriptionId     string // subscription id
	SubscriptionStatus string // sub status，0-Init | 1-Pending｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete
	RecurringAmount    string // total recurring amount, cent
	BillingType        string // 1-recurring,2-one-time
	TimeZone           string //
	CreateTime         string // create utc time
	ExternalUserId     string // external_user_id
	Status             string // 0-Active, 2-Suspend
}

// userAccountColumns holds the columns for table user_account.
var userAccountColumns = UserAccountColumns{
	Id:                 "id",
	GatewayId:          "gateway_id",
	PaymentMethod:      "payment_method",
	MerchantId:         "merchant_id",
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
	LastLoginAt:        "last_login_at",
	IsRisk:             "is_risk",
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
	CountryCode:        "country_code",
	CountryName:        "country_name",
	SubscriptionName:   "subscription_name",
	SubscriptionId:     "subscription_id",
	SubscriptionStatus: "subscription_status",
	RecurringAmount:    "recurring_amount",
	BillingType:        "billing_type",
	TimeZone:           "time_zone",
	CreateTime:         "create_time",
	ExternalUserId:     "external_user_id",
	Status:             "status",
}

// NewUserAccountDao creates and returns a new DAO object for table data access.
func NewUserAccountDao() *UserAccountDao {
	return &UserAccountDao{
		group:   "default",
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
