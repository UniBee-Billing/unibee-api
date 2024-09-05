package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"math"
	"os"
	"testing"
	"unibee/internal/consts"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	_ "unibee/test"
	"unibee/utility"
)

func TestGenerateInvoicePdf(t *testing.T) {
	ctx := context.Background()
	one := query.GetInvoiceByInvoiceId(ctx, "iv20240202ERExKnb6OhMfyyY")
	utility.Assert(one != nil, "one not found")
	one.RefundId = "refundId"
	one.SendNote = "iv20240202ERExKnb6OhMfyyY"
	var savePath = fmt.Sprintf("%s.pdf", "pdf_test")
	err := createInvoicePdf(one, query.GetMerchantById(ctx, one.MerchantId), query.GetUserAccountById(ctx, one.UserId), query.GetGatewayById(ctx, one.GatewayId), savePath)
	utility.AssertError(err, "Pdf Generator Error")
	err = os.Remove("f18f4fce-802b-471c-9418-9640384594f6.jpg")
	if err != nil {
		return
	}
	err = os.Remove("pdf_test.pdf")
	if err != nil {
		return
	}
}

func TestInvoicePdfGenerateAndEmailSendBackground(t *testing.T) {
	_ = SendInvoiceEmailToUser(context.Background(), "iv20240316QYuw5DQGcABgHDn", true, "")
}

func TestGenerate(t *testing.T) {
	var savePath = fmt.Sprintf("%s.pdf", "pdf_test")
	err := createInvoicePdf(&entity.Invoice{
		InvoiceId:                      "81720768257606",
		GmtCreate:                      gtime.Now(),
		TotalAmount:                    20000,
		TaxAmount:                      2000,
		DiscountAmount:                 2000,
		DiscountCode:                   "code11",
		VatNumber:                      "xxxxxVat",
		CountryCode:                    "EE",
		SubscriptionAmountExcludingTax: 20000,
		Currency:                       "USD",
		Lines:                          "[{\"currency\":\"USD\",\"amount\":100,\"amountExcludingTax\":100,\"tax\":12,\"unitAmountExcludingTax\":100,\"description\":\"1 × 1美金计划(测试专用) (at $1.00 / day)\",\"proration\":false,\"quantity\":1,\"periodEnd\":1705108316,\"periodStart\":1705021916},{\"currency\":\"USD\",\"amount\":0,\"amountExcludingTax\":0,\"tax\":0,\"unitAmountExcludingTax\":0,\"description\":\"0 × 3美金Addon(测试专用) (at $3.00 / day)\",\"proration\":false,\"quantity\":0,\"periodEnd\":1705108316,\"periodStart\":1705021916},{\"currency\":\"USD\",\"amount\":350,\"amountExcludingTax\":350,\"tax\":0,\"unitAmountExcludingTax\":350,\"description\":\"1 × testUpgrade (at $3.50 / day)\",\"proration\":false,\"quantity\":1,\"periodEnd\":1705108316,\"periodStart\":1705021916}]",
		Status:                         consts.InvoiceStatusPaid,
		GmtModify:                      gtime.Now(),
		Link:                           "http://unibee.top",
		TaxPercentage:                  2000,
		RefundId:                       "dddd",
		CreateFrom:                     "Refund Requested: xxxxxxxxx",
		MetaData:                       utility.MarshalToJsonString(map[string]interface{}{"ShowDetailItem": true, "LocalizedCurrency": "EUR", "LocalizedExchangeRate": 1.5, "IssueVatNumber": " EE101775690", "IssueRegNumber": "12660871", "IssueCompanyName": "Multilogin Software OÜ", "IssueAddress": "Supluse pst 1 - 201A, Tallinn Harju maakond, 119112 Harju maakond, 11911  Harju maakond, 11911"}),
	}, &entity.Merchant{
		CompanyName: "Multilogin OÜ",
		BusinessNum: "EE101775690",
		Name:        "UniBee",
		Idcard:      "12660871",
		Location:    "Supluse",
		Address:     "Supluse ",
		IsDeleted:   0,
		CompanyLogo: "http://unibee.top/files/invoice/cm/czi8o0j0jqd87mqwta.png",
	}, &entity.UserAccount{
		IsDeleted:          0,
		Email:              "jack.fu@wowow.io",
		Address:            "Best Billing Team Ltd Dubai Hills, Duai, UAE 12345",
		FirstName:          "jack",
		LastName:           "fu",
		ZipCode:            "zipCode",
		City:               "Hangzhou",
		RegistrationNumber: "Regxxxddd",
		VATNumber:          "EE101775690",
	}, nil, savePath)
	if err != nil {
		fmt.Printf("err :%s", err.Error())
	}
	err = os.Remove("f18f4fce-802b-471c-9418-9640384594f6.jpg")
	if err != nil {
		return
	}
	//err = os.Remove("pdf_test.pdf")
	//if err != nil {
	//	return
	//}
	//fmt.Println(fmt.Sprintf("%v ", currency.NarrowSymbol(currency.ParseISO(strings.ToUpper("DD")))))
}

func TestTimeFormat(t *testing.T) {
	v := 1 - (1 / (1 + utility.ConvertTaxPercentageToInternalFloat(2000)))
	fmt.Println(int(math.Floor(float64(-12000) * v)))
}
