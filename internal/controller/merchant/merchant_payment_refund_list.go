package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/payment/service"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) RefundList(ctx context.Context, req *payment.RefundListReq) (res *payment.RefundListRes, err error) {
	list, total, err := service.RefundList(ctx, &service.RefundListInternalReq{
		MerchantId: _interface.GetMerchantId(ctx),
		PaymentId:  req.PaymentId,
		Status:     req.Status,
		GatewayId:  req.GatewayId,
		UserId:     req.UserId,
		Email:      req.Email,
		Currency:   req.Currency,
	})
	if err != nil {
		return nil, err
	}
	return &payment.RefundListRes{RefundDetails: list, Total: total}, nil
}
