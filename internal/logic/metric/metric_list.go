package metric

import (
	"context"
	"strings"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type ListInternalReq struct {
	MerchantId      uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	SortField       string `json:"sortField" dc:"Sort, gmt_create" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page,Start 0" `
	Count           int    `json:"count" dc:"Count Of Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
	SkipTotal       bool
}

func MerchantMetricList(ctx context.Context, req *ListInternalReq) ([]*bean.MerchantMetric, int) {
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	var total = 0
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}

	var isDeletes = []int{0}
	utility.Assert(req.MerchantId > 0, "merchantId not found")
	var sortKey = "gmt_create desc"
	if len(req.SortField) > 0 {
		utility.Assert(strings.Contains("gmt_create", req.SortField), "sortField should one of gmt_create")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}

	var list = make([]*bean.MerchantMetric, 0)
	var entities []*entity.MerchantMetric
	q := dao.MerchantMetric.Ctx(ctx).
		Where(dao.MerchantMetric.Columns().MerchantId, req.MerchantId).
		WhereIn(dao.MerchantMetricEvent.Columns().IsDeleted, isDeletes)
	if req.CreateTimeStart > 0 {
		q = q.WhereGTE(dao.MerchantMetric.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		q = q.WhereLTE(dao.MerchantMetric.Columns().CreateTime, req.CreateTimeEnd)
	}
	q = q.Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty()
	var err error
	if req.SkipTotal {
		err = q.Scan(&entities)
	} else {
		err = q.ScanAndCount(&entities, &total, true)
	}
	if err == nil && len(entities) > 0 {
		for _, one := range entities {
			list = append(list, bean.SimplifyMerchantMetric(one))
		}
	}

	return list, total
}
