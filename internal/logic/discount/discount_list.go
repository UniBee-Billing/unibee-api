package discount

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
	"unibee/api/bean/detail"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type ListInternalReq struct {
	MerchantId      uint64
	DiscountType    []int  `json:"discountType"  dc:"discount_type, 1-percentage, 2-fixed_amount" `
	BillingType     []int  `json:"billingType"  dc:"billing_type, 1-one-time, 2-recurring" `
	Status          []int  `json:"status" dc:"status, 1-editable, 2-active, 3-deactive, 4-expire" `
	Code            string `json:"code" dc:"Filter Code"  `
	SearchKey       string `json:"searchKey" dc:"Search Key, code or name"  `
	Currency        string `json:"currency" dc:"Filter Currency"  `
	SortField       string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page, Start 0" `
	Count           int    `json:"count"  dc:"Count Of Per Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd" `
	SkipTotal       bool
}

func MerchantDiscountCodeList(ctx context.Context, req *ListInternalReq) ([]*detail.MerchantDiscountCodeDetail, int) {
	var mainList = make([]*detail.MerchantDiscountCodeDetail, 0)
	var list []*entity.MerchantDiscountCode
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
	q := dao.MerchantDiscountCode.Ctx(ctx)
	if len(req.DiscountType) > 0 {
		q = q.WhereIn(dao.MerchantDiscountCode.Columns().DiscountType, req.DiscountType)
	}
	if len(req.BillingType) > 0 {
		q = q.WhereIn(dao.MerchantDiscountCode.Columns().BillingType, req.BillingType)
	}
	if len(req.Status) > 0 {
		if len(req.Status) == 1 && utility.IsIntInArray(req.Status, consts.DiscountStatusArchived) {
			q = q.WhereGT(dao.MerchantDiscountCode.Columns().IsDeleted, 0)
		} else if utility.IsIntInArray(req.Status, consts.DiscountStatusArchived) {
			//q = q.WhereIn(dao.MerchantDiscountCode.Columns().Status, req.Status)
			q = q.Where(q.Builder().WhereOrIn(dao.MerchantDiscountCode.Columns().Status, req.Status).
				WhereOrGT(dao.MerchantDiscountCode.Columns().IsDeleted, 0))
		} else {
			q = q.WhereIn(dao.MerchantDiscountCode.Columns().Status, req.Status).Where(dao.MerchantDiscountCode.Columns().IsDeleted, 0)
		}
	}
	if len(req.Currency) > 0 {
		q = q.WhereIn(dao.MerchantDiscountCode.Columns().Currency, strings.ToUpper(req.Currency))
	}
	if len(req.SearchKey) > 0 {
		q = q.Where(q.Builder().WhereOrLike(dao.MerchantDiscountCode.Columns().Code, "%"+req.SearchKey+"%").
			WhereOrLike(dao.MerchantDiscountCode.Columns().Name, "%"+req.SearchKey+"%"))
	} else if len(req.Code) > 0 {
		q = q.WhereLike(dao.MerchantDiscountCode.Columns().Code, "%"+req.Code+"%")
	}
	if req.CreateTimeStart > 0 {
		q = q.WhereGTE(dao.MerchantDiscountCode.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		q = q.WhereLTE(dao.MerchantDiscountCode.Columns().CreateTime, req.CreateTimeEnd)
	}
	var err error
	q = q.
		Where(dao.MerchantDiscountCode.Columns().MerchantId, req.MerchantId).
		Where(dao.MerchantDiscountCode.Columns().Type, 0).
		Order(fmt.Sprintf("is_deleted asc, %s", sortKey)).
		Limit(req.Page*req.Count, req.Count)
	if req.SkipTotal {
		err = q.Scan(&list)
	} else {
		err = q.ScanAndCount(&list, &total, true)
	}
	if err != nil {
		g.Log().Errorf(ctx, "MerchantDiscountCodeList err:%s", err.Error())
		return mainList, total
	}
	for _, one := range list {
		mainList = append(mainList, detail.ConvertMerchantDiscountCodeDetail(ctx, one))
	}

	return mainList, total
}
