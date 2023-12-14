// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAccount is the golang structure for table user_account.
type UserAccount struct {
	Id              uint64      `json:"id"              ` // userId
	Creator         int64       `json:"creator"         ` // 创建人
	GmtCreate       *gtime.Time `json:"gmtCreate"       ` // 创建时间
	Modifier        int64       `json:"modifier"        ` // 修改人
	GmtModify       *gtime.Time `json:"gmtModify"       ` // 修改时间
	IsDeleted       int         `json:"isDeleted"       ` // 逻辑删除
	MobilePhone     string      `json:"mobilePhone"     ` // 手机号（非登录态为openId）
	Password        string      `json:"password"        ` // 密码，md5加密
	UserName        string      `json:"userName"        ` // 用户名
	Email           string      `json:"email"           ` // 邮箱
	Gender          string      `json:"gender"          ` // 性别
	AvatarUrl       string      `json:"avatarUrl"       ` // 头像url
	GrowthValue     int64       `json:"growthValue"     ` // 成长值(使用vip_information)
	IsVip           int         `json:"isVip"           ` // （废弃）是否开通vip，0没有，1开通，默认0
	OpenTime        *gtime.Time `json:"openTime"        ` // （废弃）会员开通时间
	ExpireTime      *gtime.Time `json:"expireTime"      ` // （废弃）会员到期时间
	FromUserId      int64       `json:"fromUserId"      ` // 分享来源，用户id（上游用户id）
	ShareCount      int         `json:"shareCount"      ` // 分享下游数量，分享成功应+1
	ReMark          string      `json:"reMark"          ` // 备注
	ParentReward    int         `json:"parentReward"    ` // 上游是否已经领过邀请人红包, 否0，是1,默认0
	InviteCode      string      `json:"inviteCode"      ` // 邀请码
	IsSpecial       int         `json:"isSpecial"       ` // 是否是特殊账号（0.否，1.是）
	ProxyId         int64       `json:"proxyId"         ` // 代理商id
	Birthday        string      `json:"birthday"        ` // 生日
	Profession      string      `json:"profession"      ` // 职业
	School          string      `json:"school"          ` // 学校
	Custom          string      `json:"custom"          ` // 其他
	NearTime        *gtime.Time `json:"nearTime"        ` // 最近登录时间
	JobNumber       string      `json:"jobNumber"       ` // 工号(新华导入用)
	Name            string      `json:"name"            ` // 员工姓名(新华导入用)
	ShopId          int64       `json:"shopId"          ` // 店铺ID(新华导入用)
	OrganId         int64       `json:"organId"         ` // 组织ID(新华导入用)
	PasswordMengyou string      `json:"passwordMengyou" ` // 盟有账号密码，md5加密
	Wid             string      `json:"wid"             ` // 盟有wid
	IsRisk          int         `json:"isRisk"          ` // 风控：0.低风险，1.中风险，2.高风险
	IsNomobileUser  int         `json:"isNomobileUser"  ` // 是否未登录态账号（无手机号），0-否，1-是
	PlainPassword   string      `json:"plainPassword"   ` // 明文密码
	Channel         string      `json:"outchannel"         ` // 渠道
}
