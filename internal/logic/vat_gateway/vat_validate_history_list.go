package vat_gateway

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type NumberValidateHistoryInternalReq struct {
	MerchantId      uint64
	SearchKey       string `json:"searchKey" dc:"Search Key, vatNumber, validateGateway, company, company address, message"  `
	VatNumber       string `json:"vatNumber" dc:"Filter Vat Number"`
	CountryCode     string `json:"countryCode" dc:"CountryCode"`
	ValidateGateway string `json:"validateGateway" dc:"Filter Validate Gateway, vatsense"`
	Status          []int  `json:"status" dc:"status, 0-Invalid，1-Valid" `
	SortField       string `json:"sortField" dc:"Sort Field，gmt_create|gmt_modify，Default gmt_modify" `
	SortType        string `json:"sortType" dc:"Sort Type，asc|desc，Default desc" `
	Page            int    `json:"page"  dc:"Page, Start 0" `
	Count           int    `json:"count"  dc:"Count Of Per Page" `
	CreateTimeStart int64  `json:"createTimeStart" dc:"CreateTimeStart，UTC timestamp，seconds" `
	CreateTimeEnd   int64  `json:"createTimeEnd" dc:"CreateTimeEnd，UTC timestamp，seconds" `
	SkipTotal       bool
}

func MerchantNumberValidateHistoryList(ctx context.Context, req *NumberValidateHistoryInternalReq) ([]*bean.MerchantVatNumberVerifyHistory, int) {
	var mainList = make([]*bean.MerchantVatNumberVerifyHistory, 0)
	var list []*entity.MerchantVatNumberVerifyHistory
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
	q := dao.MerchantVatNumberVerifyHistory.Ctx(ctx)

	if len(req.SearchKey) > 0 {
		q = q.Where(q.Builder().
			WhereOr("LOWER(vat_number) like LOWER(?)", "%"+req.SearchKey+"%").
			WhereOrLike(dao.MerchantVatNumberVerifyHistory.Columns().CompanyName, "%"+req.SearchKey+"%").
			WhereOrLike(dao.MerchantVatNumberVerifyHistory.Columns().CompanyAddress, "%"+req.SearchKey+"%").
			WhereOrLike(dao.MerchantVatNumberVerifyHistory.Columns().ValidateMessage, "%"+req.SearchKey+"%"))
	} else if len(req.VatNumber) > 0 {
		q = q.Where("LOWER(vat_number) like LOWER(?)", "%"+req.VatNumber+"%")
	}

	if len(req.CountryCode) > 0 {
		q = q.WhereIn(dao.MerchantVatNumberVerifyHistory.Columns().CountryCode, req.CountryCode)
	}

	if len(req.Status) > 0 {
		q = q.WhereIn(dao.MerchantVatNumberVerifyHistory.Columns().Valid, req.Status)
	}

	if req.CreateTimeStart > 0 {
		q = q.WhereGTE(dao.MerchantVatNumberVerifyHistory.Columns().CreateTime, req.CreateTimeStart)
	}
	if req.CreateTimeEnd > 0 {
		q = q.WhereLTE(dao.MerchantVatNumberVerifyHistory.Columns().CreateTime, req.CreateTimeEnd)
	}
	var err error
	q = q.
		Where(dao.MerchantVatNumberVerifyHistory.Columns().MerchantId, req.MerchantId).
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count)
	if req.SkipTotal {
		err = q.Scan(&list)
	} else {
		err = q.ScanAndCount(&list, &total, true)
	}
	if err != nil {
		g.Log().Errorf(ctx, "MerchantNumberValidateHistoryList err:%s", err.Error())
		return mainList, total
	}
	for _, one := range list {
		mainList = append(mainList, bean.SimplifyMerchantVatNumberVerifyHistory(one))
	}

	return mainList, total
}
