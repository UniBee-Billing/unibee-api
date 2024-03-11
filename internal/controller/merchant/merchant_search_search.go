package merchant

import (
	"context"
	"strconv"
	"strings"
	"unibee/api/merchant/search"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/auth"
	"unibee/internal/logic/invoice/service"
	"unibee/internal/query"
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
		if one != nil && one.MerchantId == _interface.GetMerchantId(ctx) {
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
			if one != nil && one.MerchantId == _interface.GetMerchantId(ctx) {
				searchResult.PrecisionMatchObject = &search.PrecisionMatchObject{
					Type: "user",
					Id:   req.SearchKey,
					Data: one,
				}
			}
		}
	}
	searchUser, err := auth.SearchUser(ctx, _interface.GetMerchantId(ctx), req.SearchKey)
	if err == nil {
		searchResult.UserAccounts = searchUser
	}
	searchInvoice, err := service.SearchInvoice(ctx, _interface.GetMerchantId(ctx), req.SearchKey)
	if err == nil {
		searchResult.Invoices = searchInvoice
	}
	return searchResult, nil
}
