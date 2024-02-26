package merchant

import (
	"context"
	"unibee/api/merchant/subscription"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

func (c *ControllerSubscription) SubscriptionAdminNoteList(ctx context.Context, req *subscription.SubscriptionAdminNoteListReq) (res *subscription.SubscriptionAdminNoteListRes, err error) {
	var mainList []*entity.SubscriptionAdminNote
	if req.Count <= 0 {
		req.Count = 20
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
		merchantUser := query.GetMerchantUserAccountById(ctx, uint64(note.MerchantUserId))
		if merchantUser != nil {
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
	}
	return &subscription.SubscriptionAdminNoteListRes{NoteLists: resultList}, nil
}
