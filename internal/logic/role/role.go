package role

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type CreateRoleInternalReq struct {
	MerchantId     uint64                         `json:"merchantId"           description:"MerchantId"`      // role
	Role           string                         `json:"role"           description:"role"`                  // role
	PermissionData []*bean.MerchantRolePermission `json:"permissionData" description:"permission_data（json）"` // permission_data（json）
}

func NewMerchantRole(ctx context.Context, req *CreateRoleInternalReq) error {
	one := query.GetRoleByName(ctx, req.MerchantId, req.Role)
	utility.Assert(one == nil, "exist role:"+req.Role)
	one = &entity.MerchantRole{
		MerchantId:     req.MerchantId,
		Role:           req.Role,
		PermissionData: utility.MarshalToJsonString(req.PermissionData),
		CreateTime:     gtime.Now().Timestamp(),
	}
	_, err := dao.MerchantMetricEvent.Ctx(ctx).Data(one).OmitNil().Insert(one)
	return err
}

func EditMerchantRole(ctx context.Context, req *CreateRoleInternalReq) error {
	one := query.GetRoleByName(ctx, req.MerchantId, req.Role)
	utility.Assert(one != nil, "role not found :"+req.Role)
	one.PermissionData = utility.MarshalToJsonString(req.PermissionData)
	_, err := dao.MerchantRole.Ctx(ctx).Data(g.Map{
		dao.MerchantRole.Columns().PermissionData: utility.MarshalToJsonString(req.PermissionData),
		dao.MerchantRole.Columns().GmtModify:      gtime.Now(),
	}).Where(dao.MerchantRole.Columns().Id, one.Id).OmitNil().Update()
	return err
}

func DeleteMerchantRole(ctx context.Context, merchantId uint64, role string) error {
	one := query.GetRoleByName(ctx, merchantId, role)
	utility.Assert(one != nil, "role not found :"+role)
	_, err := dao.MerchantRole.Ctx(ctx).Data(g.Map{
		dao.MerchantRole.Columns().IsDeleted: gtime.Now().Timestamp(),
		dao.MerchantRole.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantRole.Columns().Id, one.Id).OmitNil().Update()
	return err
}
