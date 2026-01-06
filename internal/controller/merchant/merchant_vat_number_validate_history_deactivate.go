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

func (c *ControllerVat) NumberValidateHistoryDeactivate(ctx context.Context, req *vat.NumberValidateHistoryDeactivateReq) (res *vat.NumberValidateHistoryDeactivateRes, err error) {
	utility.Assert(req.HistoryId > 0, "Invalid History Id")
	one := query.GetVatNumberValidateHistoryById(ctx, req.HistoryId)
	utility.Assert(one != nil, "Invalid History Id")
	validateMessage := one.ValidateMessage
	if !strings.HasPrefix(validateMessage, "[Manual]") {
		validateMessage = fmt.Sprintf("[Manual]%s", validateMessage)
	}
	_, err = dao.MerchantVatNumberVerifyHistory.Ctx(ctx).Data(g.Map{
		dao.MerchantVatNumberVerifyHistory.Columns().Valid:           0,
		dao.MerchantVatNumberVerifyHistory.Columns().ValidateMessage: validateMessage,
	}).Where(dao.MerchantVatNumberVerifyHistory.Columns().Id, req.HistoryId).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     _interface.Context().Get(ctx).MerchantId,
		Target:         fmt.Sprintf("VatNumberValidateHistory(%v)", req.HistoryId),
		Content:        "Deactivate",
		UserId:         0,
		SubscriptionId: "",
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return &vat.NumberValidateHistoryDeactivateRes{}, nil
}
