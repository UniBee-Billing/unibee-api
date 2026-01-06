package merchant

import (
	"context"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/auth"
)

func (c *ControllerAuth) RegisterEmailCheck(ctx context.Context, req *auth.RegisterEmailCheckReq) (res *auth.RegisterEmailCheckRes, err error) {
	utility.Assert(len(req.Email) > 0, "Email Needed")
	utility.Assert(utility.IsEmailValid(req.Email), "Invalid Email Format")
	var newOne *entity.MerchantMember
	newOne = query.GetMerchantMemberByEmail(ctx, req.Email)
	utility.Assert(newOne == nil, "Email already existed")
	return &auth.RegisterEmailCheckRes{Valid: true}, nil
}
