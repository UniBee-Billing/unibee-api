package role

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type ListReq struct {
	g.Meta `path:"/list" tags:"Role" method:"get" summary:"RoleList"`
}

type ListRes struct {
	MerchantRoles []*bean.MerchantRole `json:"merchantRoles" dc:"Merchant Roles"`
	Total         int                  `json:"total" dc:"Total"`
}

type NewReq struct {
	g.Meta      `path:"/new" tags:"Role" method:"post" summary:"NewRole"`
	Role        string                         `json:"role" dc:"Role" v:"required"`
	Permissions []*bean.MerchantRolePermission `json:"permissions" dc:"Permissions" v:"required"`
}

type NewRes struct {
}

type EditReq struct {
	g.Meta      `path:"/edit" tags:"Role" method:"post" summary:"EditRole"`
	Id          uint64                         `json:"id" dc:"id" v:"required"`
	Role        string                         `json:"role" dc:"Role" v:"required"`
	Permissions []*bean.MerchantRolePermission `json:"permissions" dc:"Permissions" v:"required"`
}

type EditRes struct {
}

type DeleteReq struct {
	g.Meta `path:"/delete" tags:"Role" method:"post" summary:"DeleteRole"`
	Id     uint64 `json:"id" dc:"id" v:"required"`
}

type DeleteRes struct {
}
