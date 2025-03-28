package handler

import (
	"context"
	"fmt"
	"github.com/go-pdf/fpdf"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/text/currency"
	"log"
	"math"
	"os"
	"strings"
	"testing"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/internal/consts"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	_ "unibee/test"
	"unibee/utility"
)

func TestGenerateInvoicePdf(t *testing.T) {
	ctx := context.Background()
	one := query.GetInvoiceByInvoiceId(ctx, "81731133642950")
	utility.Assert(one != nil, "one not found")
	//one.RefundId = "refundId"
	//one.SendNote = "iv20240202ERExKnb6OhMfyyY"
	//var savePath = fmt.Sprintf("%s.pdf", "pdf_test")
	//err := createInvoicePdf(detail.ConvertInvoiceToDetail(ctx, one), query.GetMerchantById(ctx, one.MerchantId), query.GetUserAccountById(ctx, one.UserId), query.GetGatewayById(ctx, one.GatewayId), savePath)
	//utility.AssertError(err, "Pdf Generator Error")
	//err = os.Remove("f18f4fce-802b-471c-9418-9640384594f6.jpg")
	//if err != nil {
	//	return
	//}
	//err = os.Remove("pdf_test.pdf")
	//if err != nil {
	//	return
	//}
}

func exampleFunction() {
	fmt.Println("FunctionName:", utility.ReflectCurrentFunctionName())
}

func TestInvoicePdfGenerateAndEmailSendBackground(t *testing.T) {
	exampleFunction()
}

// 210 300
func TestGenerate(t *testing.T) {
	//text := "Ü"
	//utf16Data, err := fpdf.Utf8ToUtf16(text)
	//if err != nil {
	//	fmt.Println("转换失败:", err)
	//} else {
	//	fmt.Println("UTF-16 数据:", utf16Data)
	//}
	var savePath = fmt.Sprintf("%s.pdf", "pdf_test")
	var lines []*bean.InvoiceItemSimplify
	err := utility.UnmarshalFromJsonString("[{\"currency\":\"USD\",\"amount\":100,\"amountExcludingTax\":100,\"tax\":12,\"unitAmountExcludingTax\":100,\"description\":\"1 * Custom Luxe 3 Months (2024-09-05-2024-12-05)\",\"proration\":false,\"quantity\":1,\"periodEnd\":1705108316,\"periodStart\":1705021916},{\"currency\":\"USD\",\"amount\":0,\"amountExcludingTax\":0,\"tax\":0,\"unitAmountExcludingTax\":0,\"description\":\"0 × 3 Dollar Addon(Test) (at $3.00 / day)\",\"proration\":false,\"quantity\":0,\"periodEnd\":1705108316,\"periodStart\":1705021916},{\"currency\":\"USD\",\"amount\":350,\"amountExcludingTax\":350,\"tax\":0,\"unitAmountExcludingTax\":350,\"description\":\"Remaining Time On 1 * Year 2 plan After 2024-12-20\",\"proration\":false,\"quantity\":1,\"periodEnd\":1705108316,\"periodStart\":1705021916}]", &lines)
	utility.Assert(err == nil, fmt.Sprintf("UnmarshalFromJsonString error:%v", err))
	err = createInvoicePdf(context.Background(), &detail.InvoiceDetail{
		InvoiceId:                      "81720768257606",
		GmtCreate:                      gtime.Now(),
		TotalAmount:                    20000,
		TaxAmount:                      2000,
		DiscountAmount:                 2000,
		DiscountCode:                   "code11",
		VatNumber:                      "xxxxxVat",
		CountryCode:                    "EE",
		SubscriptionAmountExcludingTax: 20000,
		Currency:                       "RUB",
		Lines:                          lines,
		Status:                         consts.InvoiceStatusPaid,
		GmtModify:                      gtime.Now(),
		Link:                           "http://unibee.top",
		TaxPercentage:                  2000,
		PromoCreditDiscountAmount:      2000,
		RefundId:                       "xxxx",
		OriginalPaymentInvoice: &bean.Invoice{
			InvoiceId: "R81720768257606",
		},
		SendNote:   "81732871446425 (Partial Refund)",
		CreateFrom: "Refund Requested: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		Metadata:   map[string]interface{}{"ShowDetailItem": true, "LocalizedCurrency": "EUR", "LocalizedExchangeRate": 4.0044715544, "IssueVatNumber": " EE101775690", "IssueRegNumber": "TRN: 104167485200003", "IssueCompanyName": "Multilogin Software OÜ", "IssueAddress": "Supluse pst 1 - 201A, Tallinn Harju maakond, 119112 Harju maakond, 11911  Harju maakond, 11911"},
	}, &entity.Merchant{
		CompanyName: "Multilogin",
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
		ZipCode:            "zipcode",
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

func TestMustParseCurrencySymbolValue(t *testing.T) {
	var symbol = fmt.Sprintf("%v ", currency.NarrowSymbol(currency.MustParseISO(strings.ToUpper("RUB"))))
	g.Log().Infof(context.Background(), "%v", symbol)
	doc := fpdf.New("P", "mm", "A4", "")
	//doc.AddFont("ArialUnicode", "", "path/to/your/font.ttf")
	doc.AddUTF8Font("dejavu", "", "./fonts/DejaVuSansCondensed.ttf")
	doc.AddUTF8Font("dejavu", "B", "./fonts/DejaVuSansCondensed-Bold.ttf")
	doc.AddUTF8Font("dejavu", "I", "./fonts/DejaVuSansCondensed-Oblique.ttf")
	doc.AddUTF8Font("dejavu", "BI", "./fonts/DejaVuSansCondensed-BoldOblique.ttf")

	doc.SetFont("dejavu", "", 16)
	doc.AddPage()
	doc.Cell(40, 10, "Цена: ₽1000")
	err := doc.OutputFileAndClose("output.pdf")
	if err != nil {
		log.Fatal(err)
	}
}
