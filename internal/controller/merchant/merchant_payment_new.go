package merchant

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"unibee/api/bean"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/payment/service"
	user2 "unibee/internal/logic/user"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) New(ctx context.Context, req *payment.NewReq) (res *payment.NewRes, err error) {
	utility.Assert(req != nil, "request req is nil")
	utility.Assert(req.GatewayId > 0, "gatewayId is nil")
	req.Currency = strings.ToUpper(req.Currency)
	merchantInfo := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	gateway := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	utility.Assert(gateway.MerchantId == merchantInfo.Id, "merchant gateway not match")

	if len(req.CancelUrl) > 0 {
		if req.Metadata == nil {
			req.Metadata = make(map[string]interface{})
		}
		req.Metadata["CancelUrl"] = req.CancelUrl
	}

	var user *entity.UserAccount
	if _interface.Context().Get(ctx).IsOpenApiCall {
		if req.UserId == 0 {
			utility.Assert(len(req.ExternalUserId) > 0, "ExternalUserId|UserId is nil")
			utility.Assert(len(req.Email) > 0, "Email|UserId is nil")
			user, err = user2.QueryOrCreateUser(ctx, &user2.NewUserInternalReq{
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
	if req.PlanId == 0 {
		utility.Assert(req.TotalAmount > 0, "amount value is nil")
		utility.Assert(len(req.Currency) > 0, "amount currency is nil")
	} else {
		plan := query.GetPlanById(ctx, req.PlanId)
		utility.Assert(plan != nil, "plan not found")
		if req.TotalAmount > 0 {
			utility.Assert(req.TotalAmount == plan.Amount, "TotalAmount not match plan's amount")
		}
		if len(req.Currency) > 0 {
			utility.Assert(req.Currency == plan.Currency, "Currency not match plan's amount")
		}
		req.TotalAmount = plan.Amount
		req.Currency = plan.Currency
		if req.Metadata == nil {
			req.Metadata = make(map[string]interface{})
		}
		req.Metadata["PlanId"] = strconv.FormatUint(req.PlanId, 10)
	}
	currencyNumberCheck(req.TotalAmount, req.Currency)
	utility.Assert(len(req.ExternalPaymentId) > 0, "ExternalPaymentId is nil")
	utility.Assert(user != nil, "User Not Found")

	if len(req.Email) == 0 {
		req.Email = user.Email
	}

	var name = req.Name
	if len(name) == 0 {
		name = merchantInfo.Name
	}
	if len(name) == 0 {
		name = merchantInfo.CompanyName
	}
	if len(name) == 0 {
		name = "Default Product"
	}
	var sendStatus = consts.InvoiceSendStatusUnnecessary
	if req.SendInvoice {
		sendStatus = consts.InvoiceSendStatusUnSend
	}
	var invoice *bean.Invoice
	if req.Items != nil && len(req.Items) > 0 {
		var invoiceItems []*bean.InvoiceItemSimplify
		var totalAmountExcludingTax int64 = 0
		var totalAmount int64 = 0
		var totalTax int64 = 0
		for _, line := range req.Items {
			utility.Assert(line.Amount > 0, "Item Amount invalid, should > 0")
			utility.Assert(len(line.Description) > 0, "Item Description invalid")
			invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
				Currency:               line.Currency,
				TaxPercentage:          line.TaxPercentage,
				Tax:                    line.Tax,
				OriginAmount:           line.Amount,
				DiscountAmount:         0,
				Amount:                 line.Amount,
				AmountExcludingTax:     line.AmountExcludingTax,
				UnitAmountExcludingTax: line.UnitAmountExcludingTax,
				Quantity:               line.Quantity,
				Name:                   line.Name,
				Description:            line.Description,
			})
			totalTax = totalTax + line.Tax
			totalAmountExcludingTax = totalAmountExcludingTax + line.AmountExcludingTax
			totalAmount = totalAmount + line.Amount
		}
		utility.Assert(totalAmount == req.TotalAmount, "sum(items.amount) should match totalAmount")
		invoice = &bean.Invoice{
			OriginAmount:            req.TotalAmount,
			TotalAmount:             req.TotalAmount,
			Currency:                req.Currency,
			InvoiceName:             name,
			ProductName:             name,
			TotalAmountExcludingTax: totalAmountExcludingTax,
			TaxAmount:               totalTax,
			DiscountAmount:          0,
			SendStatus:              sendStatus,
			DayUtilDue:              consts.DEFAULT_DAY_UTIL_DUE,
			Lines:                   invoiceItems,
			CountryCode:             req.CountryCode,
			Metadata:                req.Metadata,
		}
	} else {
		invoice = &bean.Invoice{
			OriginAmount:            req.TotalAmount,
			TotalAmount:             req.TotalAmount,
			TotalAmountExcludingTax: req.TotalAmount,
			InvoiceName:             name,
			ProductName:             name,
			Currency:                req.Currency,
			TaxAmount:               0,
			DiscountAmount:          0,
			CountryCode:             req.CountryCode,
			SendStatus:              sendStatus,
			DayUtilDue:              consts.DEFAULT_DAY_UTIL_DUE,
			Lines: []*bean.InvoiceItemSimplify{{
				Currency:               req.Currency,
				OriginAmount:           req.TotalAmount,
				Amount:                 req.TotalAmount,
				Tax:                    0,
				DiscountAmount:         0,
				AmountExcludingTax:     req.TotalAmount,
				TaxPercentage:          0,
				UnitAmountExcludingTax: req.TotalAmount,
				Name:                   name,
				Description:            req.Description,
				Quantity:               1,
			}},
			Metadata: req.Metadata,
		}
	}
	var uniqueId = fmt.Sprintf("%d_%s", merchantInfo.Id, req.ExternalPaymentId)
	exsitPayment := query.GetPaymentByUniqueId(ctx, uniqueId)
	utility.Assert(exsitPayment == nil, "same ExternalPaymentId exist:"+req.ExternalPaymentId)

	resp, err := service.GatewayPaymentCreate(ctx, &gateway_bean.GatewayNewPaymentReq{
		CheckoutMode: true,
		Gateway:      gateway,
		Pay: &entity.Payment{
			ExternalPaymentId: req.ExternalPaymentId,
			BizType:           consts.BizTypeOneTime,
			UserId:            user.Id,
			GatewayId:         gateway.Id,
			TotalAmount:       req.TotalAmount,
			Currency:          req.Currency,
			CountryCode:       req.CountryCode,
			MerchantId:        merchantInfo.Id,
			CompanyId:         merchantInfo.CompanyId,
			ReturnUrl:         req.RedirectUrl,
			GasPayer:          req.GasPayer,
			UniqueId:          uniqueId,
		},
		ExternalUserId: req.ExternalUserId,
		Email:          req.Email,
		Metadata:       req.Metadata,
		Invoice:        invoice,
	})
	utility.Assert(err == nil, fmt.Sprintf("%+v", err))
	res = &payment.NewRes{
		Status:            consts.PaymentCreated,
		PaymentId:         resp.Payment.PaymentId,
		ExternalPaymentId: req.ExternalPaymentId,
		Link:              resp.Link,
		Action:            resp.Action,
	}
	return
}
