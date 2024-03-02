package merchant

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) NewAdminNote(ctx context.Context, req *subscription.NewAdminNoteReq) (res *subscription.NewAdminNoteRes, err error) {
	note := &entity.SubscriptionAdminNote{
		SubscriptionId:   req.SubscriptionId,
		MerchantMemberId: req.MerchantMemberId,
		Note:             req.Note,
		CreateTime:       gtime.Now().Timestamp(),
	}

	_, err = dao.SubscriptionAdminNote.Ctx(ctx).Data(note).OmitNil().Insert(note)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	return
}
