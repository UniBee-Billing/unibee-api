package member

import (
	"context"
	"fmt"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"

	"github.com/gogf/gf/v2/frame/g"
)

type MemberListReq struct {
	MerchantId      uint64
	SearchKey       string   `json:"searchKey" dc:"Search Key, FirstName,LastName or Email"  `
	Email           string   `json:"email" dc:"Search Filter Email" `
	RoleIds         []uint64 `json:"roleIds" description:"The member roleId if specified'"`
	Page            int      `json:"page"  description:"Page, Start With 0" `
	Count           int      `json:"count"  description:"Count Of Page"`
	CreateTimeStart int64    `json:"createTimeStart" dc:"CreateTimeStart，UTC timestamp，seconds" `
	CreateTimeEnd   int64    `json:"createTimeEnd" dc:"CreateTimeEnd，UTC timestamp，seconds" `
}

func MerchantMemberList(ctx context.Context, req *MemberListReq) ([]*detail.MerchantMemberDetail, int) {
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}
	var total = 0
	var resultList = make([]*detail.MerchantMemberDetail, 0)
	var mainList = make([]*entity.MerchantMember, 0)

	q := dao.MerchantMember.Ctx(ctx).
		Where(dao.MerchantMember.Columns().MerchantId, req.MerchantId).
		Where(dao.MerchantMember.Columns().IsDeleted, 0)
	if req.CreateTimeStart > 0 {
		q = q.WhereGTE(dao.MerchantMember.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		q = q.WhereLTE(dao.MerchantMember.Columns().CreateTime, req.CreateTimeEnd)
	}
	if len(req.SearchKey) > 0 {
		q = q.Where(q.Builder().
			WhereOrLike(dao.MerchantMember.Columns().Email, "%"+req.SearchKey+"%").
			WhereOrLike(dao.MerchantMember.Columns().FirstName, "%"+req.SearchKey+"%").
			WhereOrLike(dao.MerchantMember.Columns().LastName, "%"+req.SearchKey+"%"))
	}
	if req.Email != "" {
		q = q.WhereLike(dao.MerchantMember.Columns().Email, "%"+req.Email+"%")
	}

	if len(req.RoleIds) > 0 {
		orq := q.Builder()
		for _, roleId := range req.RoleIds {
			orq = orq.WhereOrLike(dao.MerchantMember.Columns().Role, "%"+fmt.Sprintf("%d", roleId)+"%")
		}
		q = q.Where(orq)
	}

	err := q.Limit(req.Page*req.Count, req.Count).
		ScanAndCount(&mainList, &total, true)
	if err != nil {
		g.Log().Errorf(ctx, "MerchantMemberList err:%s", err.Error())
		return resultList, len(resultList)
	}
	for _, one := range mainList {
		resultList = append(resultList, detail.ConvertMemberToDetail(ctx, one))
	}
	return resultList, total
}

func MerchantMemberTotalList(ctx context.Context, merchantId uint64, email string) ([]*detail.MerchantMemberDetail, int) {
	var resultList = make([]*detail.MerchantMemberDetail, 0)
	var mainList = make([]*entity.MerchantMember, 0)

	query := dao.MerchantMember.Ctx(ctx).
		Where(dao.MerchantMember.Columns().MerchantId, merchantId).
		Where(dao.MerchantMember.Columns().IsDeleted, 0)

	if email != "" {
		query = query.WhereLike(dao.MerchantMember.Columns().Email, "%"+email+"%")
	}

	err := query.Scan(&mainList)
	if err != nil {
		g.Log().Errorf(ctx, "MerchantMemberList err:%s", err.Error())
		return resultList, len(resultList)
	}
	for _, one := range mainList {
		resultList = append(resultList, detail.ConvertMemberToDetail(ctx, one))
	}
	return resultList, len(resultList)
}
