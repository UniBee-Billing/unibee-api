package merchant

import (
	"context"
	"unibee/api/merchant/subscription"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
)

func (c *ControllerSubscription) AdminNoteList(ctx context.Context, req *subscription.AdminNoteListReq) (res *subscription.AdminNoteListRes, err error) {
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
	var resultList []*subscription.AdminNoteRo
	for _, note := range mainList {
		merchantMember := query.GetMerchantMemberById(ctx, uint64(note.MerchantMemberId))
		if merchantMember != nil {
			resultList = append(resultList, &subscription.AdminNoteRo{
				GmtCreate:      note.GmtCreate,
				GmtModify:      note.GmtModify,
				SubscriptionId: note.SubscriptionId,
				UserName:       merchantMember.UserName,
				Mobile:         merchantMember.Mobile,
				Email:          merchantMember.Email,
				FirstName:      merchantMember.FirstName,
				LastName:       merchantMember.LastName,
			})
		}
	}
	return &subscription.AdminNoteListRes{NoteLists: resultList}, nil
}
