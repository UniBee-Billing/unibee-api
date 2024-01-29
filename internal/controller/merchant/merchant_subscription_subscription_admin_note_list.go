package merchant

import (
	"context"
	"go-oversea-pay/api/merchant/subscription"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
)

func (c *ControllerSubscription) SubscriptionAdminNoteList(ctx context.Context, req *subscription.SubscriptionAdminNoteListReq) (res *subscription.SubscriptionAdminNoteListRes, err error) {
	var mainList []*entity.SubscriptionAdminNote
	if req.Count <= 0 {
		req.Count = 10 //每页数量Default 10
	}
	if req.Page < 0 {
		req.Page = 0
	}
	var sortKey = "gmt_create desc"
	err = dao.SubscriptionAdminNote.Ctx(ctx).
		Where(dao.SubscriptionAdminNote.Columns().SubscriptionId, req.SubscriptionId).
		Where(dao.SubscriptionAdminNote.Columns().IsDeleted, 0).
		Limit(req.Page*req.Count, req.Count).
		Order(sortKey).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil, err
	}
	var resultList []*subscription.SubscriptionAdminNoteRo
	for _, note := range mainList {
		merchantUser := query.GetMerchantAccountById(ctx, uint64(note.MerchantUserId))
		resultList = append(resultList, &subscription.SubscriptionAdminNoteRo{
			GmtCreate:      note.GmtCreate,
			GmtModify:      note.GmtModify,
			SubscriptionId: note.SubscriptionId,
			UserName:       merchantUser.UserName,
			Mobile:         merchantUser.Mobile,
			Email:          merchantUser.Email,
			FirstName:      merchantUser.FirstName,
			LastName:       merchantUser.LastName,
		})
	}
	return &subscription.SubscriptionAdminNoteListRes{NoteLists: resultList}, nil
}
