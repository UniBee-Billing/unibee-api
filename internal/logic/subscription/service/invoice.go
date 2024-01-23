package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	"go-oversea-pay/api/merchant/invoice"
	v1 "go-oversea-pay/api/open/payment"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/gateway"
	"go-oversea-pay/internal/logic/gateway/ro"
	"go-oversea-pay/internal/logic/payment/service"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strings"
)

func CreateInvoice(ctx context.Context, req *invoice.NewInvoiceCreateReq) (res *invoice.NewInvoiceCreateRes, err error) {
	user := query.GetUserAccountById(ctx, uint64(req.UserId))
	utility.Assert(user != nil, fmt.Sprintf("send user not found:%d", req.UserId))
	utility.Assert(len(user.Email) > 0, fmt.Sprintf("send user email not found:%d", req.UserId))
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, req.ChannelId)
	utility.Assert(payChannel != nil, "payChannel not found")

	var invoiceItems []*ro.ChannelDetailInvoiceItem
	var totalAmountExcludingTax int64 = 0
	var totalTax int64 = 0
	for _, line := range req.Lines {
		amountExcludingTax := line.UnitAmountExcludingTax * line.Quantity
		tax := int64(float64(amountExcludingTax) * req.TaxPercentage) // 精度损失问题 todo mark
		invoiceItems = append(invoiceItems, &ro.ChannelDetailInvoiceItem{
			Currency:               req.Currency,
			Amount:                 amountExcludingTax + tax,
			AmountExcludingTax:     amountExcludingTax,
			UnitAmountExcludingTax: line.UnitAmountExcludingTax,
			Description:            line.Description,
		})
		totalTax = totalTax + tax
		totalAmountExcludingTax = totalAmountExcludingTax + amountExcludingTax
	}
	var totalAmount = totalTax + totalAmountExcludingTax

	//创建
	invoiceId := utility.CreateInvoiceId()
	one := &entity.Invoice{
		MerchantId:                     req.MerchantId,
		InvoiceId:                      invoiceId,
		InvoiceName:                    req.Name,
		UniqueId:                       invoiceId,
		TotalAmount:                    totalAmount,
		TotalAmountExcludingTax:        totalAmountExcludingTax,
		TaxAmount:                      totalTax,
		SubscriptionAmount:             totalAmount,
		SubscriptionAmountExcludingTax: totalAmountExcludingTax,
		Currency:                       strings.ToUpper(req.Currency),
		Lines:                          utility.MarshalToJsonString(req.Lines),
		ChannelId:                      req.ChannelId,
		Status:                         consts.InvoiceStatusPending,
		SendStatus:                     0,
		SendEmail:                      user.Email,
		UserId:                         req.UserId,
	}

	result, err := dao.Invoice.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`CreateInvoice record insert failure %s`, err)
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	return &invoice.NewInvoiceCreateRes{Invoice: one}, nil
}

func EditInvoice(ctx context.Context, req *invoice.NewInvoiceEditReq) (res *invoice.NewInvoiceEditRes, err error) {
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, fmt.Sprintf("invoice not found:%s", req.InvoiceId))
	utility.Assert(one.Status == consts.InvoiceStatusPending, "invoice not in pending status")
	utility.Assert(one.IsDeleted == 0, "invoice is deleted")
	if req.ChannelId > 0 {
		payChannel := query.GetSubscriptionTypePayChannelById(ctx, req.ChannelId)
		utility.Assert(payChannel != nil, "payChannel not found")
	} else {
		req.ChannelId = one.ChannelId
	}
	if len(req.Currency) == 0 {
		req.Currency = one.Currency
	}

	var invoiceItems []*ro.ChannelDetailInvoiceItem
	var totalAmountExcludingTax int64 = 0
	var totalTax int64 = 0
	for _, line := range req.Lines {
		amountExcludingTax := line.UnitAmountExcludingTax * line.Quantity
		tax := int64(float64(amountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
		invoiceItems = append(invoiceItems, &ro.ChannelDetailInvoiceItem{
			Currency:               req.Currency,
			Amount:                 amountExcludingTax + tax,
			AmountExcludingTax:     amountExcludingTax,
			UnitAmountExcludingTax: line.UnitAmountExcludingTax,
			Description:            line.Description,
		})
		totalTax = totalTax + tax
		totalAmountExcludingTax = totalAmountExcludingTax + amountExcludingTax
	}
	var totalAmount = totalTax + totalAmountExcludingTax

	//更新 Subscription
	_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().InvoiceName:                    req.Name,
		dao.Invoice.Columns().TotalAmount:                    totalAmount,
		dao.Invoice.Columns().TotalAmountExcludingTax:        totalAmountExcludingTax,
		dao.Invoice.Columns().TaxAmount:                      totalTax,
		dao.Invoice.Columns().SubscriptionAmount:             totalAmount,
		dao.Invoice.Columns().SubscriptionAmountExcludingTax: totalAmountExcludingTax,
		dao.Invoice.Columns().Currency:                       strings.ToUpper(req.Currency),
		dao.Invoice.Columns().Currency:                       req.Currency,
		dao.Invoice.Columns().TaxPercentage:                  req.TaxPercentage,
		dao.Invoice.Columns().ChannelId:                      req.ChannelId,
		dao.Invoice.Columns().Lines:                          utility.MarshalToJsonString(req.Lines),
		dao.Invoice.Columns().GmtModify:                      gtime.Now(),
	}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	//rowAffected, err := update.RowsAffected()
	//if rowAffected != 1 {
	//	return nil, gerror.Newf("EditInvoice update err:%s", update)
	//}
	one.Currency = req.Currency
	one.TaxPercentage = req.TaxPercentage
	one.ChannelId = req.ChannelId
	one.Lines = utility.MarshalToJsonString(req.Lines)
	return &invoice.NewInvoiceEditRes{Invoice: one}, nil
}

func DeletePendingInvoice(ctx context.Context, invoiceId string) error {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	utility.Assert(one != nil, fmt.Sprintf("invoice not found:%s", invoiceId))
	utility.Assert(one.Status == consts.InvoiceStatusPending, "invoice not in pending status")
	if one.IsDeleted == 1 {
		return nil
	} else {
		//更新 Subscription
		_, err := dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().IsDeleted: 0,
			dao.Invoice.Columns().GmtModify: gtime.Now(),
		}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		//rowAffected, err := update.RowsAffected()
		//if rowAffected != 1 {
		//	return gerror.Newf("EditInvoice update err:%s", update)
		//}
		return nil
	}
}

func CancelProcessingInvoice(ctx context.Context, invoiceId string) error {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	utility.Assert(one != nil, fmt.Sprintf("invoice not found:%s", invoiceId))
	if one.Status == consts.InvoiceStatusCancelled {
		return nil
	}
	utility.Assert(one.Status == consts.InvoiceStatusProcessing, "invoice not in processing status")
	utility.Assert(one.IsDeleted == 0, "invoice is deleted")
	if one.IsDeleted == 1 {
		return nil
	} else {
		payChannel := query.GetSubscriptionTypePayChannelById(ctx, one.ChannelId)
		utility.Assert(payChannel != nil, "payChannel not found")
		_, err := gateway.GetPayChannelServiceProvider(ctx, one.ChannelId).DoRemoteChannelInvoiceCancel(ctx, payChannel, &ro.ChannelCancelInvoiceInternalReq{
			ChannelInvoiceId: one.ChannelInvoiceId,
		})
		if err != nil {
			return gerror.Newf(`FinishInvoice failure %v`, err)
		}
		// todo mark 重新生成 cancel 状态的 pdf 并发送邮件
		//更新 Subscription
		_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
			dao.Invoice.Columns().Status:    consts.InvoiceStatusCancelled,
			dao.Invoice.Columns().GmtModify: gtime.Now(),
		}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		//rowAffected, err := update.RowsAffected()
		//if rowAffected != 1 {
		//	return gerror.Newf("EditInvoice update err:%s", update)
		//}
		return nil
	}
}

func FinishInvoice(ctx context.Context, req *invoice.ProcessInvoiceForPayReq) (*invoice.ProcessInvoiceForPayRes, error) {
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, fmt.Sprintf("invoice not found:%s", req.InvoiceId))
	utility.Assert(one.Status == consts.InvoiceStatusPending, "invoice not in pending status")
	utility.Assert(one.IsDeleted == 0, "invoice is deleted")
	payChannel := query.GetSubscriptionTypePayChannelById(ctx, one.ChannelId)
	utility.Assert(payChannel != nil, "payChannel not found")
	var lines []*ro.NewInvoiceItem
	err := utility.UnmarshalFromJsonString(one.Lines, &lines)
	if err != nil {
		return nil, err
	}
	createRes, err := gateway.GetPayChannelServiceProvider(ctx, one.ChannelId).DoRemoteChannelInvoiceCreateAndPay(ctx, payChannel, &ro.ChannelCreateInvoiceInternalReq{
		Invoice:      one,
		InvoiceLines: lines,
		PayMethod:    req.PayMethod,
		DaysUtilDue:  req.DaysUtilDue,
	})
	if err != nil {
		return nil, gerror.Newf(`FinishInvoice failure %v`, err)
	}
	//更新 Subscription
	_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().ChannelUserId:     createRes.ChannelUserId,
		dao.Invoice.Columns().ChannelInvoiceId:  createRes.ChannelInvoiceId,
		dao.Invoice.Columns().ChannelInvoicePdf: createRes.ChannelInvoicePdf,
		dao.Invoice.Columns().Status:            int(createRes.Status),
		dao.Invoice.Columns().Link:              createRes.Link,
		dao.Invoice.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	//rowAffected, err := update.RowsAffected()
	//if rowAffected != 1 {
	//	return nil, gerror.Newf("FinishInvoice update err:%s", update)
	//}
	one.Status = int(createRes.Status)
	one.Link = createRes.Link
	one.ChannelUserId = createRes.ChannelUserId
	// todo mark 下面的流程
	// todo mark 生成 pdf 并发送邮件

	return &invoice.ProcessInvoiceForPayRes{Invoice: one}, nil
}

func CreateInvoiceRefund(ctx context.Context, req *invoice.NewInvoiceRefundReq) (*entity.Refund, error) {
	utility.Assert(req.RefundAmount > 0, "refundFee should > 0")
	utility.Assert(len(req.InvoiceId) > 0, "invoiceId invalid")
	utility.Assert(len(req.Reason) > 0, "reason should not be blank")
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, "invoice not found")
	utility.Assert(one.TotalAmount > req.RefundAmount, "not enough amount to refund")
	refund, err := service.DoChannelRefund(ctx, consts.PAYMENT_BIZ_TYPE_INVOICE, &v1.RefundsReq{
		PaymentId:  one.PaymentId,
		MerchantId: one.MerchantId,
		Reference:  uuid.New().String(), //todo make internal refund reference
		Reason:     req.Reason,
		Amount: &v1.PayAmountVo{
			Currency: one.Currency,
			Value:    req.RefundAmount,
		},
	}, 0)
	if err != nil {
		return nil, err
	}
	return refund, nil
}
