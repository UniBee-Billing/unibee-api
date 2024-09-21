package discount

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type UserDiscountListInternalReq struct {
	MerchantId      uint64 `json:"merchantId" dc:"The discount's merchant" v:"required"`
	Id              uint64 `json:"id" dc:"The discount's Id" v:"required"`
	SortField       string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page, Start 0" `
	Count           int    `json:"count"  dc:"Count Of Per Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
	SkipTotal       bool
}

func MerchantUserDiscountCodeList(ctx context.Context, req *UserDiscountListInternalReq) ([]*detail.MerchantUserDiscountCodeDetail, int) {
	var mainList = make([]*detail.MerchantUserDiscountCodeDetail, 0)

	one := query.GetDiscountById(ctx, req.Id)
	if one == nil || one.MerchantId != req.MerchantId {
		return mainList, 0
	}

	var list []*entity.MerchantUserDiscountCode
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}

	var total = 0
	var sortKey = "gmt_create desc"
	if len(req.SortField) > 0 {
		utility.Assert(strings.Contains("gmt_create|gmt_modify", req.SortField), "sortField should one of gmt_create|gmt_modify")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}
	q := dao.MerchantUserDiscountCode.Ctx(ctx)
	if req.CreateTimeStart > 0 {
		q = q.WhereGTE(dao.MerchantUserDiscountCode.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		q = q.WhereLTE(dao.MerchantUserDiscountCode.Columns().CreateTime, req.CreateTimeEnd)
	}
	var err error
	q = q.
		Where(dao.MerchantUserDiscountCode.Columns().MerchantId, req.MerchantId).
		Where(dao.MerchantUserDiscountCode.Columns().IsDeleted, 0).
		Where(dao.MerchantUserDiscountCode.Columns().Status, 1).
		Where(dao.MerchantUserDiscountCode.Columns().Code, one.Code).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count)
	if req.SkipTotal {
		err = q.Scan(&list)
	} else {
		err = q.ScanAndCount(&list, &total, true)
	}

	if err != nil {
		g.Log().Errorf(ctx, "MerchantUserDiscountCodeList err:%s", err.Error())
		return mainList, total
	}
	for _, one := range list {
		mainList = append(mainList, detail.ConvertMerchantUserDiscountCodeDetail(ctx, one))
	}

	return mainList, total
}
