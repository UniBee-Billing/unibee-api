package merchant

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/merchant/subscription"
)

func (c *ControllerSubscription) SubscriptionNewAdminNote(ctx context.Context, req *subscription.SubscriptionNewAdminNoteReq) (res *subscription.SubscriptionNewAdminNoteRes, err error) {
	note := &entity.SubscriptionAdminNote{
		SubscriptionId: req.SubscriptionId,
		MerchantUserId: req.MerchantUserId,
		Note:           req.Note,
	}

	_, err = dao.SubscriptionAdminNote.Ctx(ctx).Data(note).OmitNil().Insert(note)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	return
}
