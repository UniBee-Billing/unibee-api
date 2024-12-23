package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/payment/service"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) List(ctx context.Context, req *payment.ListReq) (res *payment.ListRes, err error) {
	paymentDetails, total, err := service.PaymentList(ctx, &service.PaymentListInternalReq{
		MerchantId:  _interface.GetMerchantId(ctx),
		GatewayId:   req.GatewayId,
		UserId:      req.UserId,
		Email:       req.Email,
		Status:      req.Status,
		Currency:    req.Currency,
		CountryCode: req.CountryCode,
		SortField:   req.SortField,
		SortType:    req.SortType,
		Page:        req.Page,
		Count:       req.Count,
	})
	if err != nil {
		return nil, err
	}
	return &payment.ListRes{PaymentDetails: paymentDetails, Total: total}, nil
}
