package merchant

import (
	"context"
	"strings"
	"unibee/api/merchant/user"
	_interface "unibee/internal/interface/context"
	user2 "unibee/internal/logic/user"
)

func (c *ControllerUser) List(ctx context.Context, req *user.ListReq) (res *user.ListRes, err error) {
	result, err := user2.UserList(ctx, &user2.UserListInternalReq{
		MerchantId:     _interface.GetMerchantId(ctx),
		UserId:         req.UserId,
		Email:          strings.Trim(req.Email, " "),
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		SubscriptionId: req.SubscriptionId,
		SubStatus:      req.SubStatus,
		Status:         req.Status,
		DeleteInclude:  req.DeleteInclude,
		SortField:      req.SortField,
		SortType:       req.SortType,
		Page:            req.Page,
		Count:           req.Count,
		CreateTimeStart: req.CreateTimeStart,
		CreateTimeEnd:   req.CreateTimeEnd,
		PlanIds:         req.PlanIds,
	})
	if err != nil {
		return nil, err
	}
	return &user.ListRes{UserAccounts: result.UserAccounts, Total: result.Total}, nil
}
