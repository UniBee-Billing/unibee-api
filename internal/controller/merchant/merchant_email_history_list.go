package merchant

import (
	"context"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"

	"unibee/api/merchant/email"
	_interface "unibee/internal/interface/context"
	email2 "unibee/internal/logic/email"
)

func (c *ControllerEmail) HistoryList(ctx context.Context, req *email.HistoryListReq) (res *email.HistoryListRes, err error) {
	merchantId := _interface.GetMerchantId(ctx)

	// Convert API request to internal request
	internalReq := &email2.EmailHistoryListInternalReq{
		MerchantId:      merchantId,
		Email:           req.Email,
		Status:          req.Status,
		SearchKey:       req.SearchKey,
		SortField:       req.SortField,
		SortType:        req.SortType,
		Page:            req.Page,
		Count:           req.Count,
		CreateTimeStart: req.CreateTimeStart,
		CreateTimeEnd:   req.CreateTimeEnd,
	}

	// Call the internal list function
	emailHistories, total := email2.MerchantEmailHistoryList(ctx, internalReq)

	statistics := &detail.MerchantEmailHistoryStatistics{
		TotalSend:    0,
		TotalSuccess: 0,
		TotalFail:    0,
	}
	totalSend, _ := dao.MerchantEmailHistory.Ctx(ctx).
		Where(dao.MerchantEmailHistory.Columns().MerchantId, merchantId).
		Count()
	statistics.TotalSend = int64(totalSend)
	totalSuccess, _ := dao.MerchantEmailHistory.Ctx(ctx).
		Where(dao.MerchantEmailHistory.Columns().MerchantId, merchantId).
		Where(dao.MerchantEmailHistory.Columns().Status, 1).
		Count()
	statistics.TotalSuccess = int64(totalSuccess)
	totalFailure, _ := dao.MerchantEmailHistory.Ctx(ctx).
		Where(dao.MerchantEmailHistory.Columns().MerchantId, merchantId).
		Where(dao.MerchantEmailHistory.Columns().Status, 2).
		Count()
	statistics.TotalFail = int64(totalFailure)

	return &email.HistoryListRes{
		EmailHistoryStatistics: statistics,
		EmailHistories:         emailHistories,
		Total:                  total,
	}, nil
}
