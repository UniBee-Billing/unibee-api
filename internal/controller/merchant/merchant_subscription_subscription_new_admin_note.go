package merchant

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

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
	if err != nil {
		g.Log().Printf(ctx, "NewAdminNote :%s", err.Error())
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	return
}
