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
	Id              string // userId
	Creator         string // 创建人
	GmtCreate       string // 创建时间
	Modifier        string // 修改人
	GmtModify       string // 修改时间
	IsDeleted       string // 逻辑删除
	MobilePhone     string // 手机号（非登录态为openId）
	Password        string // 密码，md5加密
	UserName        string // 用户名
	Email           string // 邮箱
	Gender          string // 性别
	AvatarUrl       string // 头像url
	GrowthValue     string // 成长值(使用vip_information)
	IsVip           string // （废弃）是否开通vip，0没有，1开通，默认0
	OpenTime        string // （废弃）会员开通时间
	ExpireTime      string // （废弃）会员到期时间
	FromUserId      string // 分享来源，用户id（上游用户id）
	ShareCount      string // 分享下游数量，分享成功应+1
	ReMark          string // 备注
	ParentReward    string // 上游是否已经领过邀请人红包, 否0，是1,默认0
	InviteCode      string // 邀请码
	IsSpecial       string // 是否是特殊账号（0.否，1.是）
	ProxyId         string // 代理商id
	Birthday        string // 生日
	Profession      string // 职业
	School          string // 学校
	Custom          string // 其他
	NearTime        string // 最近登录时间
	JobNumber       string // 工号(新华导入用)
	Name            string // 员工姓名(新华导入用)
	ShopId          string // 店铺ID(新华导入用)
	OrganId         string // 组织ID(新华导入用)
	PasswordMengyou string // 盟有账号密码，md5加密
	Wid             string // 盟有wid
	IsRisk          string // 风控：0.低风险，1.中风险，2.高风险
	IsNomobileUser  string // 是否未登录态账号（无手机号），0-否，1-是
	PlainPassword   string // 明文密码
	Channel         string // 渠道
}

// userAccountColumns holds the columns for table user_account.
var userAccountColumns = UserAccountColumns{
	Id:              "id",
	Creator:         "creator",
	GmtCreate:       "gmt_create",
	Modifier:        "modifier",
	GmtModify:       "gmt_modify",
	IsDeleted:       "is_deleted",
	MobilePhone:     "mobile_phone",
	Password:        "password",
	UserName:        "user_name",
	Email:           "email",
	Gender:          "gender",
	AvatarUrl:       "avatar_url",
	GrowthValue:     "growth_value",
	IsVip:           "is_vip",
	OpenTime:        "open_time",
	ExpireTime:      "expire_time",
	FromUserId:      "from_user_id",
	ShareCount:      "share_count",
	ReMark:          "re_mark",
	ParentReward:    "parent_reward",
	InviteCode:      "invite_code",
	IsSpecial:       "is_special",
	ProxyId:         "proxy_id",
	Birthday:        "birthday",
	Profession:      "profession",
	School:          "school",
	Custom:          "custom",
	NearTime:        "near_time",
	JobNumber:       "job_number",
	Name:            "name",
	ShopId:          "shop_id",
	OrganId:         "organ_id",
	PasswordMengyou: "password_mengyou",
	Wid:             "wid",
	IsRisk:          "is_risk",
	IsNomobileUser:  "is_nomobile_user",
	PlainPassword:   "plain_password",
	Channel:         "outchannel",
}

// NewUserAccountDao creates and returns a new DAO object for table data access.
func NewUserAccountDao() *UserAccountDao {
	return &UserAccountDao{
		group:   "auth",
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
