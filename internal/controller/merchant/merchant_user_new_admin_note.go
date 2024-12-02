package merchant

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"

	"unibee/api/merchant/user"
)

func (c *ControllerUser) NewAdminNote(ctx context.Context, req *user.NewAdminNoteReq) (res *user.NewAdminNoteRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "invalid token")
	utility.Assert(len(req.Note) < 20000, "note too long")
	note := &entity.UserAdminNote{
		UserId:           req.UserId,
		MerchantMemberId: int64(_interface.Context().Get(ctx).MerchantMember.Id),
		Note:             req.Note,
		CreateTime:       gtime.Now().Timestamp(),
	}

	_, err = dao.UserAdminNote.Ctx(ctx).Data(note).OmitNil().Insert(note)
	utility.AssertError(err, "Save Error")
	return &user.NewAdminNoteRes{}, nil
}
