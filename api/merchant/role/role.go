package role

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type ListReq struct {
	g.Meta `path:"/list" tags:"Role" method:"get" summary:"Get Merchant Role List"`
}

type ListRes struct {
	MerchantRoles []*bean.MerchantRoleSimplify `json:"merchantRoles" dc:"Merchant Roles"`
}

type NewReq struct {
	g.Meta      `path:"/new" tags:"Role" method:"post" summary:"New Merchant Role"`
	Role        string                         `json:"role" dc:"Code" v:"required"`
	Permissions []*bean.MerchantRolePermission `json:"permissions" dc:"Permissions" v:"required"`
}

type NewRes struct {
}

type EditReq struct {
	g.Meta      `path:"/edit" tags:"Role" method:"post" summary:"Edit Merchant Role"`
	Role        string                         `json:"role" dc:"Code" v:"required"`
	Permissions []*bean.MerchantRolePermission `json:"permissions" dc:"Permissions" v:"required"`
}

type EditRes struct {
}

type DeleteReq struct {
	g.Meta `path:"/delete" tags:"Role" method:"post" summary:"Delete Merchant Role"`
	Role   string `json:"role" dc:"Code" v:"required"`
}

type DeleteRes struct {
}
