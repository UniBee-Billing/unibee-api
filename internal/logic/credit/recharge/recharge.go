package recharge

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/credit/account"
	"unibee/internal/logic/credit/config"
	"unibee/internal/logic/credit/recharge/handler"
	"unibee/internal/logic/payment/service"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type CreditRechargeInternalReq struct {
	UserId         uint64 `json:"userId"`
	MerchantId     uint64 `json:"merchantId"`
	GatewayId      uint64 `json:"gatewayId"`
	RechargeAmount int64  `json:"rechargeAmount"`
	Currency       string `json:"currency"`
	Name           string `json:"name"             description:"recharge name"`         // recharge name
	Description    string `json:"description"       description:"recharge description"` // recharge description
	ReturnUrl      string `json:"returnUrl"  dc:"ReturnUrl"  `
	CancelUrl      string `json:"cancelUrl" dc:"CancelUrl"`
}

type CreditRechargeInternalRes struct {
	User           *bean.UserAccount    `json:"user"`
	Merchant       *bean.Merchant       `json:"merchant"`
	Gateway        *detail.Gateway      `json:"gateway"`
	CreditAccount  *bean.CreditAccount  `json:"creditAccount"`
	CreditRecharge *bean.CreditRecharge `json:"creditRecharge"`
	Invoice        *bean.Invoice        `json:"invoice"`
	Payment        *bean.Payment        `json:"payment"`
	Link           string               `json:"link"`
	Paid           bool                 `json:"paid" dc:"Paid，true|false"`
}

func CreateRechargePayment(ctx context.Context, req *CreditRechargeInternalReq) (*CreditRechargeInternalRes, error) {
	utility.Assert(req != nil, "request is nil")
	utility.Assert(req.UserId > 0, "invalid userId")
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	utility.Assert(req.GatewayId > 0, "invalid gatewayId")
	gateway := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	utility.Assert(req.RechargeAmount > 0, "invalid rechargeAmount")
	utility.Assert(req.Currency != "", "invalid currency")
	user := query.GetUserAccountById(ctx, uint64(req.UserId))
	utility.Assert(user != nil, "user not found")
	utility.Assert(user.MerchantId == req.MerchantId, "merchant not match")
	creditAccount := account.QueryOrCreateCreditAccount(ctx, req.UserId, req.Currency, consts.CreditAccountTypeMain)
	utility.Assert(creditAccount != nil, "credit creditAccount create failed")
	utility.Assert(creditAccount.Type == consts.CreditAccountTypeMain, "invalid credit account type, should be main account")
	utility.Assert(creditAccount.RechargeEnable == 1, "Credit account recharge disabled")
	utility.AssertError(config.CheckCreditConfigRecharge(ctx, req.MerchantId, creditAccount.Type, req.Currency), "Invalid Credit Config")

	one := &entity.CreditRecharge{
		UserId:            req.UserId,
		MerchantId:        req.MerchantId,
		CreditId:          creditAccount.Id,
		RechargeId:        utility.CreateCreditRechargeId(),
		RechargeStatus:    consts.CreditRechargeCreated,
		Currency:          req.Currency,
		TotalAmount:       req.RechargeAmount,
		PaymentAmount:     req.RechargeAmount,
		Name:              req.Name,
		Description:       req.Description,
		GatewayId:         req.GatewayId,
		TotalRefundAmount: 0,
		CreateTime:        gtime.Now().Timestamp(),
		AccountType:       creditAccount.Type,
	}
	result, err := dao.CreditRecharge.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`create credit recharge record insert failure %s`, err)
	}
	id, _ := result.LastInsertId()
	one.Id = id
	// create invoice and payment
	{
		//name and description
		if len(req.Name) == 0 {
			req.Name = "Credit Recharge"
		}
		if len(req.Description) == 0 {
			req.Description = "Credit Recharge"
		}
	}

	rechargeInvoice := &entity.Invoice{
		BizType:                        consts.BizTypeCreditRecharge,
		MerchantId:                     req.MerchantId,
		InvoiceId:                      utility.CreateInvoiceId(),
		InvoiceName:                    req.Name,
		ProductName:                    req.Description,
		UniqueId:                       one.RechargeId,
		TotalAmount:                    one.TotalAmount,
		TotalAmountExcludingTax:        one.TotalRefundAmount,
		TaxAmount:                      0,
		TaxPercentage:                  0,
		SubscriptionAmount:             one.TotalAmount,
		SubscriptionAmountExcludingTax: one.TotalAmount,
		Currency:                       strings.ToUpper(one.Currency),
		Lines: utility.MarshalToJsonString([]*bean.InvoiceItemSimplify{
			{
				Currency:               strings.ToUpper(one.Currency),
				OriginAmount:           one.TotalAmount,
				DiscountAmount:         0,
				Amount:                 one.TotalAmount,
				Tax:                    0,
				AmountExcludingTax:     one.TotalAmount,
				TaxPercentage:          0,
				UnitAmountExcludingTax: one.TotalAmount,
				Name:                   req.Name,
				Description:            req.Description,
				PdfDescription:         "",
				Proration:              false,
				Quantity:               1,
				PeriodEnd:              0,
				PeriodStart:            0,
				Plan:                   nil,
			},
		}),
		GatewayId:   req.GatewayId,
		Status:      consts.InvoiceStatusPending,
		SendStatus:  consts.InvoiceSendStatusUnSend,
		SendEmail:   user.Email,
		UserId:      user.Id,
		CreateTime:  gtime.Now().Timestamp(),
		CountryCode: user.CountryCode,
		CreateFrom:  "Admin",
	}

	result, err = dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		g.Log().Errorf(ctx, "CreateRechargePayment CreateInvoice Error:%s\n", err.Error())
		_ = handler.HandleCreditRechargeFailed(ctx, one.RechargeId)
		return nil, gerror.Newf(`CreateInvoice record insert failure %s`, err)
	}
	id, _ = result.LastInsertId()
	one.Id = id
	// update invoiceId to credit recharge
	_, err = dao.CreditRecharge.Ctx(ctx).Data(g.Map{
		dao.CreditRecharge.Columns().InvoiceId: one.InvoiceId,
		dao.CreditRecharge.Columns().GmtModify: gtime.Now(),
	}).Where(dao.CreditRecharge.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		g.Log().Errorf(ctx, "CreateRechargePayment Update Recharge invoiceId Error:%s\n", err.Error())
		_ = handler.HandleCreditRechargeFailed(ctx, one.RechargeId)
		return nil, err
	}

	createRes, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, rechargeInvoice, true, req.ReturnUrl, req.CancelUrl, "Credit Recharge", 0)
	if err != nil {
		g.Log().Errorf(ctx, "CreateRechargePayment CreateSubInvoicePaymentDefaultAutomatic Credit Recharge Error:%s\n", err.Error())
		_ = handler.HandleCreditRechargeFailed(ctx, one.RechargeId)
		return nil, err
	}
	// update payment to credit recharge if success
	if createRes.Status == consts.PaymentSuccess {
		err = handler.HandleCreditRechargeSuccess(ctx, one.RechargeId, rechargeInvoice, createRes.Payment)
		if err != nil {
			g.Log().Errorf(ctx, "CreateRechargePayment HandleCreditRechargeSuccess Error:%s\n", err.Error())
			return nil, err
		}
	}

	return &CreditRechargeInternalRes{
		User:           bean.SimplifyUserAccount(user),
		Gateway:        detail.ConvertGatewayDetail(ctx, gateway),
		CreditRecharge: bean.SimplifyCreditRecharge(one),
		Invoice:        bean.SimplifyInvoice(rechargeInvoice),
		Payment:        bean.SimplifyPayment(createRes.Payment),
		Link:           createRes.Link,
		Paid:           createRes.Status == consts.PaymentSuccess,
	}, nil
}
