package role

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"math"
	"strconv"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/member"
	"unibee/internal/logic/operation_log"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type CreateRoleInternalReq struct {
	Id             uint64                         `json:"id"           description:"RoleId"`                  // id
	MerchantId     uint64                         `json:"merchantId"           description:"MerchantId"`      // role
	Role           string                         `json:"role"           description:"role"`                  // role
	PermissionData []*bean.MerchantRolePermission `json:"permissionData" description:"permission_data（json）"` // permission_data（json）
}

func NewMerchantRole(ctx context.Context, req *CreateRoleInternalReq) error {
	utility.Assert(req.Role != "Owner", "Invalid Role, Role 'Owner' is reserved")
	one := query.GetRoleByName(ctx, req.MerchantId, req.Role)
	utility.Assert(one == nil, "exist role:"+req.Role)
	one = &entity.MerchantRole{
		MerchantId:     req.MerchantId,
		Role:           req.Role,
		PermissionData: utility.MarshalToJsonString(req.PermissionData),
		CreateTime:     gtime.Now().Timestamp(),
	}
	result, err := dao.MerchantRole.Ctx(ctx).Data(one).OmitNil().Insert(one)
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Role(%v)", one.Id),
		Content:        "New",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return err
}

func EditMerchantRole(ctx context.Context, req *CreateRoleInternalReq) error {
	utility.Assert(req.Id > 0, "Invalid Id")
	utility.Assert(req.Role != "Owner", "Invalid Role, Role 'Owner' is reserved")
	one := query.GetRoleById(ctx, req.Id)
	utility.Assert(one != nil, "role not found :"+req.Role)
	utility.Assert(one.MerchantId == req.MerchantId, "Permission Deny")
	one.PermissionData = utility.MarshalToJsonString(req.PermissionData)
	_, err := dao.MerchantRole.Ctx(ctx).Data(g.Map{
		dao.MerchantRole.Columns().Role:           req.Role,
		dao.MerchantRole.Columns().PermissionData: utility.MarshalToJsonString(req.PermissionData),
		dao.MerchantRole.Columns().GmtModify:      gtime.Now(),
	}).Where(dao.MerchantRole.Columns().Id, one.Id).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Role(%v)", one.Id),
		Content:        "Edit",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return err
}

func DeleteMerchantRole(ctx context.Context, merchantId uint64, id uint64) error {
	one := query.GetRoleById(ctx, id)
	utility.Assert(one != nil, "role not found :"+strconv.FormatUint(id, 10))
	utility.Assert(one.MerchantId == merchantId, "Permission Deny")
	memberList, _ := GetMemberListByRoleId(ctx, merchantId, id)
	utility.Assert(len(memberList) == 0, fmt.Sprintf("Role has %v members binding can not delete", len(memberList)))
	_, err := dao.MerchantRole.Ctx(ctx).Data(g.Map{
		dao.MerchantRole.Columns().IsDeleted: gtime.Now().Timestamp(),
		dao.MerchantRole.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantRole.Columns().Id, one.Id).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Role(%v)", one.Id),
		Content:        "Delete",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return err
}

func HardDeleteMerchantRole(ctx context.Context, merchantId uint64, role string) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(len(role) > 0, "invalid role")
	_, err := dao.MerchantRole.Ctx(ctx).Where(dao.MerchantRole.Columns().Role, role).Where(dao.MerchantRole.Columns().MerchantId, merchantId).Delete()
	return err
}

func GetMemberListByRoleId(ctx context.Context, merchantId uint64, roleId uint64) ([]*detail.MerchantMemberDetail, int) {
	resultList := make([]*detail.MerchantMemberDetail, 0)
	totalList, _ := member.MerchantMemberTotalList(ctx, merchantId)
	for _, one := range totalList {
		var found = false
		for _, role := range one.MemberRoles {
			if role.Id == roleId {
				found = true
			}
		}
		if found {
			resultList = append(resultList, one)
		}
	}
	return resultList, len(resultList)
}

func GetMemberListByRoleIds(ctx context.Context, merchantId uint64, roleIds []uint64, page int, count int) ([]*detail.MerchantMemberDetail, int) {
	if count <= 0 {
		count = 20
	}
	if page < 0 {
		page = 0
	}
	resultList := make([]*detail.MerchantMemberDetail, 0)
	totalList, _ := member.MerchantMemberTotalList(ctx, merchantId)
	for _, one := range totalList {
		var found = false
		for _, role := range one.MemberRoles {
			for _, roleId := range roleIds {
				if role.Id == roleId {
					found = true
				}
			}
		}
		if found {
			resultList = append(resultList, one)
		}
	}
	sliceStart, sliceEnd := SlicePage(page+1, count, len(resultList))
	return resultList[sliceStart:sliceEnd], len(resultList)
}

func SlicePage(page, pageSize, nums int) (sliceStart, sliceEnd int) {
	if page < 0 {
		page = 1
	}

	if pageSize < 0 {
		pageSize = 20
	}

	if pageSize > nums {
		return 0, nums
	}

	pageCount := int(math.Ceil(float64(nums) / float64(pageSize)))
	if page > pageCount {
		return 0, 0
	}
	sliceStart = (page - 1) * pageSize
	sliceEnd = sliceStart + pageSize

	if sliceEnd > nums {
		sliceEnd = nums
	}
	return sliceStart, sliceEnd
}
