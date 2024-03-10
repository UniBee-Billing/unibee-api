package service

import (
	"context"
	"strings"
	"unibee/api/merchant/payment"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/gateway/ro"
	entity "unibee/internal/model/entity/oversea_pay"
	query2 "unibee/internal/query"
	"unibee/utility"
)

type RefundListInternalReq struct {
	MerchantId uint64 `json:"merchantId"   dc:"MerchantId"`
	PaymentId  string `json:"paymentId" dc:"PaymentId" v:"required"`
	Status     int    `json:"status" dc:"Status,10-create|20-success|30-Failed|40-Reverse"`
	GatewayId  uint64 `json:"gatewayId"   dc:"GatewayId"`
	UserId     int64  `json:"userId" dc:"UserId"`
	Email      string `json:"email" dc:"Email"`
	Currency   string `json:"currency" dc:"Currency"`
}

func RefundList(ctx context.Context, req *RefundListInternalReq) (RefundDetails []*payment.RefundDetail, err error) {
	var mainList []*payment.RefundDetail
	utility.Assert(req.MerchantId > 0, "merchantId not found")
	utility.Assert(len(req.PaymentId) > 0, "PaymentId not found")
	var sortKey = "create_time desc"
	query := dao.Refund.Ctx(ctx).
		Where(dao.Refund.Columns().MerchantId, req.MerchantId)
	if req.GatewayId > 0 {
		query = query.Where(dao.Refund.Columns().GatewayId, req.GatewayId)
	}
	if req.UserId > 0 {
		query = query.Where(dao.Refund.Columns().UserId, req.UserId)
	} else if len(req.Email) > 0 {
		user := query2.GetUserAccountByEmail(ctx, req.MerchantId, req.Email)
		if user != nil {
			query = query.Where(dao.Refund.Columns().UserId, user.Id)
		} else {
			return mainList, nil
		}
	}
	if req.Status > 0 {
		query = query.Where(dao.Refund.Columns().Status, req.Status)
	}
	if len(req.Currency) > 0 {
		query = query.Where(dao.Refund.Columns().Currency, strings.ToUpper(req.Currency))
	}

	var list []*entity.Refund
	err = query.
		Order(sortKey).
		OmitEmpty().Scan(&list)
	if err != nil {
		return nil, err
	}
	for _, one := range list {
		mainList = append(mainList, &payment.RefundDetail{
			User:    ro.SimplifyUserAccount(query2.GetUserAccountById(ctx, uint64(one.UserId))),
			Payment: ro.SimplifyPayment(query2.GetPaymentByPaymentId(ctx, one.PaymentId)),
			Refund:  ro.SimplifyRefund(one),
		})
	}

	return mainList, nil
}
