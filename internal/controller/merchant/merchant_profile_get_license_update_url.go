package merchant

import (
	"context"
	"fmt"
	go_redismq "github.com/jackyang-hk/go-redismq"
	_interface "unibee/internal/interface/context"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/profile"
)

type NewLicenseUpdateUrlReq struct {
	MerchantId uint64 `json:"merchantId"`
	PlanId     int64  `json:"planId" dc:"Id of plan to update" dc:"Id of plan to update"`
	ReturnUrl  string `json:"returnUrl"  dc:"ReturnUrl"`
	CancelUrl  string `json:"cancelUrl" dc:"CancelUrl"`
}

func (c *ControllerProfile) GetLicenseUpdateUrl(ctx context.Context, req *profile.GetLicenseUpdateUrlReq) (res *profile.GetLicenseUpdateUrlRes, err error) {
	updateUrlRes := go_redismq.Invoke(ctx, &go_redismq.InvoiceRequest{
		Group:  "GID_UniBee_License",
		Method: "GetLicenseUpdateUrl",
		Request: utility.MarshalToJsonString(&NewLicenseUpdateUrlReq{
			MerchantId: _interface.GetMerchantId(ctx),
			PlanId:     req.PlanId,
			ReturnUrl:  req.ReturnUrl,
			CancelUrl:  req.CancelUrl,
		}),
	}, 0)
	if !updateUrlRes.Status {
		return nil, gerror.New(fmt.Sprintf("Server error:%s", updateUrlRes.Response))
	}
	return &profile.GetLicenseUpdateUrlRes{Url: fmt.Sprintf("%s", updateUrlRes.Response)}, nil
}
