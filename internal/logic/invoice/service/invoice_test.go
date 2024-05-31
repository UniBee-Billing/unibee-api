package service

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/stretchr/testify/require"
	"testing"
	"unibee/api/bean/detail"
	"unibee/api/merchant/invoice"
	"unibee/internal/consts"
	detail2 "unibee/internal/logic/invoice/detail"
	"unibee/test"
)

func TestInvoice(t *testing.T) {
	ctx := context.Background()
	var one *detail.InvoiceDetail
	var err error
	t.Run("Test for invoice Create|Edit|List|Delete", func(t *testing.T) {
		res, err := CreateInvoice(ctx, test.TestMerchant.Id, &invoice.NewReq{
			UserId:        test.TestUser.Id,
			TaxPercentage: 1000,
			GatewayId:     test.TestGateway.Id,
			Currency:      "USD",
			Name:          "test_invoice",
			Lines: []*invoice.NewInvoiceItemParam{{
				UnitAmountExcludingTax: 100,
				Description:            "test",
				Quantity:               1,
			}},
			Finish: false,
		})
		require.Nil(t, err)
		require.NotNil(t, res)
		one = res.Invoice
		require.NotNil(t, one)
		one = detail2.InvoiceDetail(ctx, one.InvoiceId)
		require.Equal(t, "USD", one.Currency)
		require.Equal(t, "test_invoice", one.InvoiceName)
		require.Equal(t, int64(110), one.TotalAmount)
		_, err = EditInvoice(ctx, &invoice.EditReq{
			InvoiceId:     one.InvoiceId,
			TaxPercentage: 0,
			GatewayId:     test.TestGateway.Id,
			Currency:      "EUR",
			Name:          "test_invoice_2",
			Lines: []*invoice.NewInvoiceItemParam{{
				UnitAmountExcludingTax: 100,
				Description:            "test",
				Quantity:               1,
			}},
			Finish: false,
		})
		require.Nil(t, err)
		one = detail2.InvoiceDetail(ctx, one.InvoiceId)
		require.Equal(t, "EUR", one.Currency)
		require.Equal(t, "test_invoice_2", one.InvoiceName)
		require.Equal(t, int64(100), one.TotalAmount)
		list, err := InvoiceList(ctx, &InvoiceListInternalReq{
			MerchantId: test.TestMerchant.Id,
			SortField:  "gmt_create",
			SortType:   "desc",
			FirstName:  test.TestUser.FirstName,
			SendEmail:  test.TestUser.Email,
			Page:       -1,
			Count:      0,
		})
		require.Nil(t, err)
		require.Equal(t, true, len(list.Invoices) > 0)
		searchInvoice, err := SearchInvoice(ctx, test.TestMerchant.Id, one.InvoiceId)
		require.Nil(t, err)
		require.NotNil(t, searchInvoice)
		require.Equal(t, 1, len(searchInvoice))
		err = DeletePendingInvoice(ctx, one.InvoiceId)
		require.Nil(t, err)
		list, err = InvoiceList(ctx, &InvoiceListInternalReq{
			MerchantId: test.TestMerchant.Id,
			Page:       -1,
			Count:      0,
		})
		require.Nil(t, err)
		require.Equal(t, true, len(list.Invoices) > 0)
	})
	t.Run("Test for invoice HardDelete", func(t *testing.T) {
		err = HardDeleteInvoice(ctx, one.MerchantId, one.InvoiceId)
		require.Nil(t, err)
	})

	t.Run("Test for invoice Create|Finish|Link", func(t *testing.T) {
		res, err := CreateInvoice(ctx, test.TestMerchant.Id, &invoice.NewReq{
			UserId:        test.TestUser.Id,
			TaxPercentage: 1000,
			GatewayId:     test.TestGateway.Id,
			Currency:      "USD",
			Name:          "test_invoice",
			Lines: []*invoice.NewInvoiceItemParam{{
				UnitAmountExcludingTax: 100,
				Description:            "test",
				Quantity:               1,
			}},
			Finish: false,
		})
		require.Nil(t, err)
		require.NotNil(t, res)
		one = res.Invoice
		require.NotNil(t, one)
		one = detail2.InvoiceDetail(ctx, one.InvoiceId)
		require.Equal(t, "USD", one.Currency)
		require.Equal(t, "test_invoice", one.InvoiceName)
		require.Equal(t, int64(110), one.TotalAmount)
		finishInvoice, err := FinishInvoice(ctx, &invoice.FinishReq{
			InvoiceId: one.InvoiceId,
			//PayMethod:   2,
			DaysUtilDue: 2,
		})
		require.Nil(t, err)
		require.NotNil(t, finishInvoice)
		require.Equal(t, consts.InvoiceStatusProcessing, finishInvoice.Invoice.Status)
		require.NotNil(t, finishInvoice.Invoice.Link)
		checkRes := LinkCheck(ctx, one.InvoiceId, gtime.Now().Timestamp())
		require.NotNil(t, checkRes)
		require.Equal(t, true, len(checkRes.Link) > 0)
		checkRes = LinkCheck(ctx, one.InvoiceId, gtime.Now().AddDate(0, 0, 3).Timestamp())
		require.NotNil(t, checkRes)
		require.Equal(t, true, len(checkRes.Link) == 0)
		require.Equal(t, true, len(checkRes.Message) > 0)
		err = CancelProcessingInvoice(ctx, one.InvoiceId, "AutoTestCancel")
		require.Nil(t, err)
		checkRes = LinkCheck(ctx, one.InvoiceId, gtime.Now().Timestamp())
		require.NotNil(t, checkRes)
		require.Equal(t, true, len(checkRes.Link) == 0)
		require.Equal(t, true, len(checkRes.Message) > 0)
		one = detail2.InvoiceDetail(ctx, one.InvoiceId)
		require.Equal(t, "USD", one.Currency)
		require.Equal(t, "test_invoice", one.InvoiceName)
		require.Equal(t, int64(110), one.TotalAmount)
		require.Equal(t, consts.InvoiceStatusCancelled, one.Status)
	})
	t.Run("Test for invoice HardDelete", func(t *testing.T) {
		err = HardDeleteInvoice(ctx, one.MerchantId, one.InvoiceId)
		require.Nil(t, err)
	})
}
