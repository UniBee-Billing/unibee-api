package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"os"
	"testing"
	"unibee/internal/consts"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	_ "unibee/test"
	"unibee/utility"
)

func TestGenerateInvoicePdf(t *testing.T) {
	ctx := context.Background()
	one := query.GetInvoiceByInvoiceId(ctx, "iv20240202ERExKnb6OhMfyyY")
	utility.Assert(one != nil, "one not found")
	var savePath = fmt.Sprintf("%s.pdf", "pdf_test")
	err := createInvoicePdf(one, query.GetMerchantById(ctx, one.MerchantId), query.GetUserAccountById(ctx, one.UserId), savePath)
	utility.AssertError(err, "Pdf")
	err = os.Remove("f18f4fce-802b-471c-9418-9640384594f6.jpg")
	if err != nil {
		return
	}
	err = os.Remove("pdf_test.pdf")
	if err != nil {
		return
	}
}

func TestGenerate(t *testing.T) {
	var savePath = fmt.Sprintf("%s.pdf", "pdf_test")
	err := createInvoicePdf(&entity.Invoice{
		InvoiceId:                      "in20240111j91EsJ8qGR9gBjI",
		GmtCreate:                      gtime.Now(),
		TotalAmount:                    20000,
		TaxAmount:                      2000,
		DiscountAmount:                 2000,
		SubscriptionAmountExcludingTax: 20000,
		Currency:                       "USD",
		Lines:                          "[{\"currency\":\"USD\",\"amount\":100,\"amountExcludingTax\":100,\"tax\":0,\"unitAmountExcludingTax\":100,\"description\":\"1 × 1美金计划(测试专用) (at $1.00 / day)\",\"proration\":false,\"quantity\":1,\"periodEnd\":1705108316,\"periodStart\":1705021916},{\"currency\":\"USD\",\"amount\":0,\"amountExcludingTax\":0,\"tax\":0,\"unitAmountExcludingTax\":0,\"description\":\"0 × 3美金Addon(测试专用) (at $3.00 / day)\",\"proration\":false,\"quantity\":0,\"periodEnd\":1705108316,\"periodStart\":1705021916},{\"currency\":\"USD\",\"amount\":350,\"amountExcludingTax\":350,\"tax\":0,\"unitAmountExcludingTax\":350,\"description\":\"1 × testUpgrade (at $3.50 / day)\",\"proration\":false,\"quantity\":1,\"periodEnd\":1705108316,\"periodStart\":1705021916}]",
		Status:                         consts.InvoiceStatusPaid,
		GmtModify:                      gtime.Now(),
		Link:                           "http://unibee.top",
		TaxPercentage:                  2000,
	}, &entity.Merchant{
		CompanyName: "UniBee.inc",
		BusinessNum: "EE101775690",
		Name:        "UniBee",
		Idcard:      "12660871",
		Location:    "Supluse",
		Address:     "pst 1-201A, Tallinn Harju maakond, 11911",
		IsDeleted:   0,
		CompanyLogo: "https://imagesize.hknet-inc.com/sp/files/f18f4fce-802b-471c-9418-9640384594f6.jpg",
	}, &entity.UserAccount{
		IsDeleted: 0,
		Address:   "Best Billing Team Ltd Dubai Hills, Duai, UAE 12345",
		FirstName: "Yvonne",
		LastName:  "Wang",
	}, savePath)
	if err != nil {
		fmt.Printf("err :%s", err.Error())
	}
	err = os.Remove("f18f4fce-802b-471c-9418-9640384594f6.jpg")
	if err != nil {
		return
	}
	err = os.Remove("pdf_test.pdf")
	if err != nil {
		return
	}
}

func TestTimeFormat(t *testing.T) {
	fmt.Println(gtime.Now().Layout("2006-01-02 15:04:05"))
}
