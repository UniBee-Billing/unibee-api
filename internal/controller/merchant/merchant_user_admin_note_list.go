package merchant

import (
	"context"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"

	"unibee/api/merchant/user"
)

func (c *ControllerUser) AdminNoteList(ctx context.Context, req *user.AdminNoteListReq) (res *user.AdminNoteListRes, err error) {
	var mainList []*entity.UserAdminNote
	if req.Count <= 0 {
		req.Count = 20
	}
	if req.Page < 0 {
		req.Page = 0
	}
	var sortKey = "gmt_create desc"
	err = dao.UserAdminNote.Ctx(ctx).
		Where(dao.UserAdminNote.Columns().UserId, req.UserId).
		Where(dao.UserAdminNote.Columns().IsDeleted, 0).
		Limit(req.Page*req.Count, req.Count).
		Order(sortKey).
		OmitEmpty().Scan(&mainList)
	if err != nil {
		return nil, err
	}
	var resultList []*detail.UserAdminNoteDetail
	for _, note := range mainList {
		merchantMember := query.GetMerchantMemberById(ctx, uint64(note.MerchantMemberId))
		if merchantMember != nil {
			resultList = append(resultList, detail.ConvertUserAdminNoteDetail(ctx, note))
		}
	}
	return &user.AdminNoteListRes{NoteLists: resultList}, nil
}
