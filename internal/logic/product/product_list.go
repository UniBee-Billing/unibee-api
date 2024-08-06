package product

import (
	"context"
	"unibee/internal/query"

	"strings"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type ListInternalReq struct {
	MerchantId uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	Status     []int  `json:"status" dc:"Filter, Default All，,Status，1-active，2-inactive" `
	SortField  string `json:"sortField" dc:"Sort Field，id|create_time|gmt_modify，Default id" `
	SortType   string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page       int    `json:"page"  dc:"Page, Start 0" `
	Count      int    `json:"count"  dc:"Count Of Per Page" `
}

func ProductList(ctx context.Context, req *ListInternalReq) (list []*bean.Product, total int) {
	var mainList []*entity.Product
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}
	var sortKey = "id desc"
	if len(req.SortField) > 0 {
		utility.Assert(strings.Contains("id|create_time|gmt_modify", req.SortField), "sortField should one of create_time|gmt_modify")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}
	q := dao.Product.Ctx(ctx).
		Where(dao.Product.Columns().MerchantId, req.MerchantId).
		Where(dao.Product.Columns().IsDeleted, 0)
	if len(req.Status) > 0 {
		q = q.WhereIn(dao.Product.Columns().Status, req.Status)
	}
	err := q.OmitEmpty().
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		ScanAndCount(&mainList, &total, true)
	if err != nil {
		return nil, 0
	}
	if len(req.Status) == 0 || utility.IsIntInArray(req.Status, 1) {
		list = append(list, bean.SimplifyProduct(query.GetDefaultProduct()))
		total = total + 1
	}
	for _, one := range mainList {
		list = append(list, bean.SimplifyProduct(one))
	}
	return list, total
}
