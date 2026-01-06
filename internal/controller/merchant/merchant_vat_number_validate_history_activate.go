package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/operation_log"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/vat"
)

func (c *ControllerVat) NumberValidateHistoryActivate(ctx context.Context, req *vat.NumberValidateHistoryActivateReq) (res *vat.NumberValidateHistoryActivateRes, err error) {
	utility.Assert(req.HistoryId > 0, "Invalid History Id")
	utility.Assert(req.CountryCode != "", "Invalid Country Code")
	err = utility.ValidateCountryCode(req.CountryCode)
	utility.AssertError(err, "Invalid Country Code")
	utility.Assert(req.CompanyName != "", "Invalid Company Name")
	utility.Assert(req.CompanyAddress != "", "Invalid Company Address")
	one := query.GetVatNumberValidateHistoryById(ctx, req.HistoryId)
	utility.Assert(one != nil, "Invalid History Id")
	validateMessage := one.ValidateMessage
	if !strings.HasPrefix(validateMessage, "[Manual]") {
		validateMessage = fmt.Sprintf("[Manual]%s", validateMessage)
	}
	_, err = dao.MerchantVatNumberVerifyHistory.Ctx(ctx).Data(g.Map{
		dao.MerchantVatNumberVerifyHistory.Columns().Valid:           1,
		dao.MerchantVatNumberVerifyHistory.Columns().IsDeleted:       0,
		dao.MerchantVatNumberVerifyHistory.Columns().CompanyAddress:  req.CompanyAddress,
		dao.MerchantVatNumberVerifyHistory.Columns().CompanyName:     req.CompanyName,
		dao.MerchantVatNumberVerifyHistory.Columns().CountryCode:     req.CountryCode,
		dao.MerchantVatNumberVerifyHistory.Columns().ValidateMessage: validateMessage,
	}).Where(dao.MerchantVatNumberVerifyHistory.Columns().Id, req.HistoryId).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     _interface.Context().Get(ctx).MerchantId,
		Target:         fmt.Sprintf("VatNumberValidateHistory(%v)", req.HistoryId),
		Content:        "Activate",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return &vat.NumberValidateHistoryActivateRes{}, nil
}
