package merchant

import (
	"context"
	"strings"
	_interface "unibee/internal/interface"
	member2 "unibee/internal/logic/operation_log"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) OperationLogList(ctx context.Context, req *member.OperationLogListReq) (res *member.OperationLogListRes, err error) {
	list, total := member2.MerchantOperationLogList(ctx, &member2.OperationLogListInternalReq{
		MerchantId:      _interface.GetMerchantId(ctx),
		MemberFirstName: req.MemberFirstName,
		MemberLastName:  req.MemberLastName,
		MemberEmail:     req.MemberEmail,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           strings.Trim(req.Email, " "),
		SubscriptionId:  req.SubscriptionId,
		InvoiceId:       req.InvoiceId,
		PlanId:          req.PlanId,
		DiscountCode:    req.DiscountCode,
		Page:            req.Page,
		Count:           req.Count,
	})
	return &member.OperationLogListRes{MerchantOperationLogs: list, Total: total}, nil
}
