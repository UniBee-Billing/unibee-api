package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"strings"
	"unibee/api/merchant/payment"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/payment/service"
	user2 "unibee/internal/logic/user"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) NewPayment(ctx context.Context, req *subscription.NewPaymentReq) (res *subscription.NewPaymentRes, err error) {
	utility.Assert(req != nil, "request req is nil")
	utility.Assert(req.GatewayId > 0, "gatewayId is nil")
	req.Currency = strings.ToUpper(req.Currency)
	merchantInfo := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	gateway := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	utility.Assert(gateway.MerchantId == merchantInfo.Id, "merchant gateway not match")

	var user *entity.UserAccount
	if _interface.Context().Get(ctx).IsOpenApiCall {
		if req.UserId == 0 {
			utility.Assert(len(req.ExternalUserId) > 0, "ExternalUserId|UserId is nil")
			utility.Assert(len(req.Email) > 0, "Email|UserId is nil")
			user, err = user2.QueryOrCreateUser(ctx, &user2.NewReq{
				ExternalUserId: req.ExternalUserId,
				Email:          req.Email,
				MerchantId:     merchantInfo.Id,
			})
			utility.AssertError(err, "Server Error")
		} else {
			user = query.GetUserAccountById(ctx, req.UserId)
		}
		utility.Assert(user != nil, "User Not Found")
		if len(req.ExternalPaymentId) == 0 {
			req.ExternalPaymentId = uuid.New().String()
		}
	} else {
		user = query.GetUserAccountById(ctx, _interface.Context().Get(ctx).User.Id)
		utility.Assert(user != nil, "User Not Found")
		if req.UserId > 0 {
			utility.Assert(user.Id == req.UserId, "user not match")
		}
		if len(req.ExternalPaymentId) == 0 {
			req.ExternalPaymentId = uuid.New().String()
		}
	}
	// cancel other pending payment
	var oldList = make([]*entity.Payment, 0)
	_ = dao.Payment.Ctx(ctx).
		Where(dao.Payment.Columns().UserId, user.Id).
		Where(dao.Payment.Columns().Status, consts.PaymentCreated).
		Where(dao.Payment.Columns().BizType, consts.BizTypeOneTime).
		OmitEmpty().Scan(&oldList)
	go func() {
		defer func() {
			if exception := recover(); exception != nil {
				fmt.Printf("SubscriptionNewPayment PaymentGatewayCancel panic error:%s\n", exception)
				return
			}
		}()
		backgroundCtx := context.Background()
		for _, oldOne := range oldList {
			err = service.PaymentGatewayCancel(backgroundCtx, oldOne)
			if err != nil {
				g.Log().Errorf(backgroundCtx, "SubscriptionNewPayment NewPayment error:%s", err.Error())
			}
		}
	}()

	controllerPayment := ControllerPayment{}
	paymentRes, paymentErr := controllerPayment.New(ctx, &payment.NewReq{
		ExternalPaymentId: req.ExternalPaymentId,
		ExternalUserId:    req.ExternalUserId,
		Email:             req.Email,
		UserId:            user.Id,
		Currency:          req.Currency,
		TotalAmount:       req.TotalAmount,
		PlanId:            req.PlanId,
		GatewayId:         req.GatewayId,
		RedirectUrl:       req.RedirectUrl,
		CancelUrl:         req.CancelUrl,
		CountryCode:       req.CountryCode,
		Name:              req.Name,
		Description:       req.Description,
		Items:             req.Items,
		Metadata:          req.Metadata,
		GasPayer:          req.GasPayer,
		SendInvoice:       true,
	})

	if paymentErr != nil {
		return nil, paymentErr
	}
	return &subscription.NewPaymentRes{
		Status:            paymentRes.Status,
		PaymentId:         paymentRes.PaymentId,
		ExternalPaymentId: paymentRes.ExternalPaymentId,
		Link:              paymentRes.Link,
		Action:            paymentRes.Action,
	}, nil
}
