package metric_event

import (
	"context"
	"strings"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type EventListInternalReq struct {
	MerchantId      uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
	UserId          int64  `json:"userId" dc:"Filter UserId, Default All" `
	SortField       string `json:"sortField" dc:"Sort，user_id|gmt_create" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page,Start 0" `
	Count           int    `json:"count" dc:"Count Of Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
	SkipTotal       bool
}

type EventListInternalRes struct {
	Events []*detail.MerchantMetricEventDetail `json:"events" description:"Event Metric Event List" `
	Total  int                                 `json:"total" dc:"Total"`
}

func EventList(ctx context.Context, req *EventListInternalReq) (res *EventListInternalRes, err error) {
	var mainList []*entity.MerchantMetricEvent
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
		utility.Assert(strings.Contains("user_id|gmt_create", req.SortField), "sortField should one of user_id|gmt_create")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}
	q := dao.MerchantMetricEvent.Ctx(ctx).
		Where(dao.MerchantMetricEvent.Columns().MerchantId, req.MerchantId).
		WhereIn(dao.MerchantMetricEvent.Columns().IsDeleted, isDeletes)
	if req.UserId > 0 {
		q = q.Where(dao.MerchantMetricEvent.Columns().Id, req.UserId)
	}
	if req.CreateTimeStart > 0 {
		q = q.WhereGTE(dao.UserAccount.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		q = q.WhereLTE(dao.UserAccount.Columns().CreateTime, req.CreateTimeEnd)
	}
	q = q.Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty()
	if req.SkipTotal {
		err = q.Scan(&mainList)
	} else {
		err = q.ScanAndCount(&mainList, &total, true)
	}
	if err != nil {
		return nil, err
	}
	var resultList = make([]*detail.MerchantMetricEventDetail, 0)
	for _, one := range mainList {
		resultList = append(resultList, detail.ConvertMerchantMetricEventDetail(ctx, one))
	}
	return &EventListInternalRes{Events: resultList, Total: total}, nil
}
