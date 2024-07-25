package service

import (
	"context"
	"strings"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	detail2 "unibee/internal/logic/invoice/detail"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type PaymentListInternalReq struct {
	MerchantId  uint64 `json:"merchantId"   dc:"MerchantId"`
	GatewayId   uint64 `json:"gatewayId"   dc:"GatewayId"`
	UserId      uint64 `json:"userId" dc:"UserId " `
	Email       string `json:"email" dc:"Email"`
	Status      int    `json:"status" dc:"Status, 10-Created|20-Success|30-Failed|40-Cancelled"`
	Currency    string `json:"currency" dc:"Currency"`
	CountryCode string `json:"countryCode" dc:"CountryCode"`
	SortField   string `json:"sortField" dc:"Sort Field，user_id|create_time|status" `
	SortType    string `json:"sortType" dc:"Sort Type，asc|desc" `
	Page        int    `json:"page"  dc:"Page, Start With 0" `
	Count       int    `json:"count"  dc:"Count" dc:"Count Of Page" `
}

type PaymentListInternalRes struct {
	PaymentDetails []*detail.PaymentDetail `json:"paymentDetails" dc:"PaymentDetails"`
}

func PaymentList(ctx context.Context, req *PaymentListInternalReq) (PaymentDetails []*detail.PaymentDetail, total int, err error) {
	req.Currency = strings.ToUpper(req.Currency)
	var mainList []*detail.PaymentDetail
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}

	utility.Assert(req.MerchantId > 0, "merchantId not found")
	var sortKey = "id desc"
	if len(req.SortField) > 0 {
		utility.Assert(strings.Contains("user_id|create_time|status", req.SortField), "sortField should one of user_id|create_time|status")
		if len(req.SortType) > 0 {
			utility.Assert(strings.Contains("asc|desc", req.SortType), "sortType should one of asc|desc")
			sortKey = req.SortField + " " + req.SortType
		} else {
			sortKey = req.SortField + " desc"
		}
	}
	q := dao.Payment.Ctx(ctx).
		Where(dao.Payment.Columns().MerchantId, req.MerchantId)
	if req.GatewayId > 0 {
		q = q.Where(dao.Payment.Columns().GatewayId, req.GatewayId)
	}
	if req.UserId > 0 {
		q = q.Where(dao.Payment.Columns().UserId, req.UserId)
	} else if len(req.Email) > 0 {
		user := query.GetUserAccountByEmail(ctx, req.MerchantId, req.Email)
		if user != nil {
			q = q.Where(dao.Payment.Columns().UserId, user.Id)
		} else {
			return mainList, 0, nil
		}
	}
	if req.Status > 0 {
		q = q.Where(dao.Payment.Columns().Status, req.Status)
	}
	if len(req.Currency) > 0 {
		q = q.Where(dao.Payment.Columns().Currency, strings.ToUpper(req.Currency))
	}
	if len(req.CountryCode) > 0 {
		q = q.Where(dao.Payment.Columns().CountryCode, req.CountryCode)
	}
	var list []*entity.Payment
	err = q.
		Order(sortKey).
		Limit(req.Page*req.Count, req.Count).
		OmitEmpty().ScanAndCount(&list, &total, true)
	if err != nil {
		return nil, 0, err
	}
	for _, one := range list {
		mainList = append(mainList, &detail.PaymentDetail{
			User:    bean.SimplifyUserAccount(query.GetUserAccountById(ctx, one.UserId)),
			Payment: bean.SimplifyPayment(one),
			Gateway: bean.SimplifyGateway(query.GetGatewayById(ctx, one.GatewayId)),
			Invoice: detail2.InvoiceDetail(ctx, one.InvoiceId),
		})
	}

	return mainList, total, nil
}
