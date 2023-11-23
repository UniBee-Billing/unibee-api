// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAccount is the golang structure of table user_account for DAO operations like Where/Data.
type UserAccount struct {
	g.Meta          `orm:"table:user_account, do:true"`
	Id              interface{} // userId
	Creator         interface{} // 创建人
	GmtCreate       *gtime.Time // 创建时间
	Modifier        interface{} // 修改人
	GmtModify       *gtime.Time // 修改时间
	IsDeleted       interface{} // 逻辑删除
	MobilePhone     interface{} // 手机号（非登录态为openId）
	Password        interface{} // 密码，md5加密
	UserName        interface{} // 用户名
	Email           interface{} // 邮箱
	Gender          interface{} // 性别
	AvatarUrl       interface{} // 头像url
	GrowthValue     interface{} // 成长值(使用vip_information)
	IsVip           interface{} // （废弃）是否开通vip，0没有，1开通，默认0
	OpenTime        *gtime.Time // （废弃）会员开通时间
	ExpireTime      *gtime.Time // （废弃）会员到期时间
	FromUserId      interface{} // 分享来源，用户id（上游用户id）
	ShareCount      interface{} // 分享下游数量，分享成功应+1
	ReMark          interface{} // 备注
	ParentReward    interface{} // 上游是否已经领过邀请人红包, 否0，是1,默认0
	InviteCode      interface{} // 邀请码
	IsSpecial       interface{} // 是否是特殊账号（0.否，1.是）
	ProxyId         interface{} // 代理商id
	Birthday        interface{} // 生日
	Profession      interface{} // 职业
	School          interface{} // 学校
	Custom          interface{} // 其他
	NearTime        *gtime.Time // 最近登录时间
	JobNumber       interface{} // 工号(新华导入用)
	Name            interface{} // 员工姓名(新华导入用)
	ShopId          interface{} // 店铺ID(新华导入用)
	OrganId         interface{} // 组织ID(新华导入用)
	PasswordMengyou interface{} // 盟有账号密码，md5加密
	Wid             interface{} // 盟有wid
	IsRisk          interface{} // 风控：0.低风险，1.中风险，2.高风险
	IsNomobileUser  interface{} // 是否未登录态账号（无手机号），0-否，1-是
	PlainPassword   interface{} // 明文密码
	Channel         interface{} // 渠道
}
