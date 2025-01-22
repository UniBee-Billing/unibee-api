package service

import (
	"context"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

func IsNewSubscriptionUser(ctx context.Context, merchantId uint64, email string) bool {
	if merchantId <= 0 || len(email) == 0 {
		return true
	}
	user := query.GetUserAccountByEmail(ctx, merchantId, email)
	if user == nil {
		return true
	}
	var one *entity.Invoice
	_ = dao.Invoice.Ctx(ctx).
		Where(dao.Invoice.Columns().UserId, user.Id).
		WhereNotNull(dao.Invoice.Columns().SubscriptionId).
		Where(dao.Invoice.Columns().BizType, consts.BizTypeSubscription).
		Where(dao.Invoice.Columns().Status, consts.InvoiceStatusPaid).
		OrderDesc(dao.Invoice.Columns().Id).
		OmitEmpty().Scan(&one)
	return one == nil
}
