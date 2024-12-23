package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	email2 "unibee/internal/logic/email"
	"unibee/internal/logic/email/sender"
	"unibee/utility"

	"unibee/api/merchant/email"
)

func (c *ControllerEmail) SenderSetup(ctx context.Context, req *email.SenderSetupReq) (res *email.SenderSetupRes, err error) {
	utility.Assert(len(req.Name) > 0, "Invalid name")
	utility.Assert(len(req.Address) > 0, "Invalid address")
	err = email2.SetupMerchantEmailSender(ctx, _interface.GetMerchantId(ctx), &sender.Sender{
		Name:    req.Name,
		Address: req.Address,
	})
	if err != nil {
		return nil, err
	}
	return &email.SenderSetupRes{}, nil
}
