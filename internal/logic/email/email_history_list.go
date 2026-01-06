package email

import (
	"context"
	"strings"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
)

type EmailHistoryListInternalReq struct {
	MerchantId      uint64
	Email           string `json:"email" dc:"Filter Email" `
	SearchKey       string `json:"searchKey" dc:"Search Key, email or title" `
	Status          []int  `json:"status" dc:"status, 0-pending, 1-success, 2-failure" `
	SortField       string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_create" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page, Start 0" `
	Count           int    `json:"count"  dc:"Count Of Per Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart，UTC timestamp，seconds" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
	SkipTotal       bool
}

func MerchantEmailHistoryList(ctx context.Context, req *EmailHistoryListInternalReq) ([]*detail.MerchantEmailHistoryDetail, int) {
	var mainList = make([]*detail.MerchantEmailHistoryDetail, 0)
	var list []*entity.MerchantEmailHistory
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
	q := dao.MerchantEmailHistory.Ctx(ctx)

	if len(req.SearchKey) > 0 {
		q = q.Where(q.Builder().
			WhereOr("LOWER(email) like LOWER(?)", "%"+req.SearchKey+"%").
			WhereOrLike(dao.MerchantEmailHistory.Columns().Title, "%"+req.SearchKey+"%").
			WhereOrLike(dao.MerchantEmailHistory.Columns().Content, "%"+req.SearchKey+"%"))
	} else if len(req.Email) > 0 {
		q = q.Where("LOWER(email) like LOWER(?)", "%"+req.Email+"%")
	}

	if len(req.Status) > 0 {
		q = q.WhereIn(dao.MerchantEmailHistory.Columns().Status, req.Status)
	}

	if req.CreateTimeStart > 0 {
		q = q.WhereGTE(dao.MerchantEmailHistory.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		q = q.WhereLTE(dao.MerchantEmailHistory.Columns().CreateTime, req.CreateTimeEnd)
	}
	var err error
	q = q.
		Where(dao.MerchantEmailHistory.Columns().MerchantId, req.MerchantId).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count)
	if req.SkipTotal {
		err = q.Scan(&list)
	} else {
		err = q.ScanAndCount(&list, &total, true)
	}
	if err != nil {
		g.Log().Errorf(ctx, "MerchantEmailHistoryList err:%s", err.Error())
		return mainList, total
	}
	for _, one := range list {
		mainList = append(mainList, detail.ConvertMerchantEmailHistoryDetail(ctx, one))
	}

	return mainList, total
}
