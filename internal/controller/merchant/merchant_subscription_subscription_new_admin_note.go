package merchant

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) NewAdminNote(ctx context.Context, req *subscription.NewAdminNoteReq) (res *subscription.NewAdminNoteRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "invalid token")
	note := &entity.SubscriptionAdminNote{
		SubscriptionId:   req.SubscriptionId,
		MerchantMemberId: int64(_interface.Context().Get(ctx).MerchantMember.Id),
		Note:             req.Note,
		CreateTime:       gtime.Now().Timestamp(),
	}

	_, err = dao.SubscriptionAdminNote.Ctx(ctx).Data(note).OmitNil().Insert(note)
	utility.AssertError(err, "Save Error")
	return &subscription.NewAdminNoteRes{}, nil
}
