package merchant

import (
	"context"
	"go-oversea-pay/api/merchant/search"
	"go-oversea-pay/internal/logic/auth"
	"go-oversea-pay/internal/logic/invoice/service"
	"go-oversea-pay/internal/query"
	"strconv"
	"strings"
)

func (c *ControllerSearch) Search(ctx context.Context, req *search.SearchReq) (res *search.SearchRes, err error) {
	if len(req.SearchKey) == 0 {
		return &search.SearchRes{}, nil
	}
	searchResult := &search.SearchRes{
		PrecisionMatchObject: nil,
		UserAccounts:         nil,
		Invoices:             nil,
	}
	if strings.HasPrefix(req.SearchKey, "in") {
		one := query.GetInvoiceByInvoiceId(ctx, req.SearchKey)
		if one != nil {
			searchResult.PrecisionMatchObject = &search.PrecisionMatchObject{
				Type: "invoice",
				Id:   req.SearchKey,
				Data: one,
			}
		}
	} else {
		searchInt, err := strconv.Atoi(req.SearchKey)
		if err == nil {
			one := query.GetUserAccountById(ctx, uint64(searchInt))
			if one != nil {
				searchResult.PrecisionMatchObject = &search.PrecisionMatchObject{
					Type: "user",
					Id:   req.SearchKey,
					Data: one,
				}
			}
		}
	}
	searchUser, err := auth.SearchUser(ctx, req.SearchKey)
	if err == nil {
		searchResult.UserAccounts = searchUser
	}
	searchInvoice, err := service.SearchInvoice(ctx, req.SearchKey)
	if err == nil {
		searchResult.Invoices = searchInvoice
	}
	return searchResult, nil
}
