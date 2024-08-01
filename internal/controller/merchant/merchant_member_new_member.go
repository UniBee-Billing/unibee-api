package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface"
	member2 "unibee/internal/logic/member"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) NewMember(ctx context.Context, req *member.NewMemberReq) (res *member.NewMemberRes, err error) {
	if config.GetConfigInstance().Mode != "standalone" || config.GetConfigInstance().Mode != "cloud" {
		return nil, gerror.New("Not Support")
	}
	err = member2.AddMerchantMember(ctx, _interface.GetMerchantId(ctx), req.Email, req.FirstName, req.LastName, req.RoleIds)
	if err != nil {
		return nil, err
	}
	return &member.NewMemberRes{}, nil
}
