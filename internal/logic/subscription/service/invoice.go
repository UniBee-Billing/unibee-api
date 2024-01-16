package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/api/merchant/invoice"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/payment/gateway"
	"go-oversea-pay/internal/logic/payment/gateway/ro"
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
	//创建
	one := &entity.Invoice{
		MerchantId:                     req.MerchantId,
		InvoiceId:                      utility.CreateInvoiceOrderNo(),
		TotalAmount:                    req.TotalAmount,
		TotalAmountExcludingTax:        req.TotalAmount,
		TaxAmount:                      0,
		SubscriptionAmount:             req.TotalAmount,
		SubscriptionAmountExcludingTax: req.TotalAmount,
		Currency:                       strings.ToUpper(req.Currency),
		Lines:                          utility.MarshalToJsonString(req.Lines),
		ChannelId:                      req.ChannelId,
		Status:                         consts.InvoiceStatusPending,
		SendStatus:                     0,
		SendEmail:                      user.Email,
		UserId:                         req.UserId,
	}

	result, err := dao.Invoice.Ctx(ctx).Data(one).OmitEmpty().Insert(one)
	if err != nil {
		return nil, gerror.Newf(`CreateInvoice record insert failure %s`, err)
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	createRes, err := gateway.GetPayChannelServiceProvider(ctx, req.ChannelId).DoRemoteChannelInvoiceCreate(ctx, payChannel, &ro.ChannelCreateInvoiceInternalReq{
		Invoice:     one,
		PayMethod:   2,
		DaysUtilDue: 1, //todo 默认值
	})
	if err != nil {
		return nil, gerror.Newf(`CreateChannelInvoice failure %v`, err)
	}
	//更新 Subscription
	update, err := dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().ChannelUserId:     createRes.ChannelUserId,
		dao.Invoice.Columns().ChannelInvoiceId:  createRes.ChannelInvoiceId,
		dao.Invoice.Columns().ChannelInvoicePdf: createRes.ChannelInvoicePdf,
		dao.Invoice.Columns().Status:            int(createRes.Status),
		dao.Invoice.Columns().Link:              createRes.Link,
		dao.Invoice.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().Id, one.Id).OmitEmpty().Update()
	if err != nil {
		return nil, err
	}
	rowAffected, err := update.RowsAffected()
	if rowAffected != 1 {
		return nil, gerror.Newf("ChannelInvoice update err:%s", update)
	}
	one.Status = int(createRes.Status)
	one.Link = createRes.Link
	one.ChannelUserId = createRes.ChannelUserId
	//todo mark 下面的流程

	return &invoice.NewInvoiceCreateRes{Invoice: one}, nil
}
