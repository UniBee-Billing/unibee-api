package invoice

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"io"
	"net/http"
	"os"
	"unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/invoice/service"
	"unibee/internal/query"
)

func VerifyInvoiceLinkSecurityToken(ctx context.Context, invoiceId string, token string) bool {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	if one == nil {
		return false
	}
	if token == one.SendTerms {
		return true
	}
	return false
}

func LinkEntry(r *ghttp.Request) {
	invoiceId := r.Get("invoiceId").String()
	if len(invoiceId) == 0 {
		r.Response.Writeln("InvoiceId not found")
		return
	}
	st := r.Get("st").String()
	if !VerifyInvoiceLinkSecurityToken(r.Context(), invoiceId, st) {
		r.Response.Writeln("Invalid link")
		return
	}
	res := service.LinkCheck(r.Context(), invoiceId, gtime.Now().Timestamp())
	if len(res.Link) > 0 {
		r.Response.RedirectTo(res.Link)
	} else if len(res.Message) > 0 {
		r.Response.Writeln(res.Message)
	} else {
		r.Response.Writeln("Server Error")
	}
}

func LinkPdfEntry(r *ghttp.Request) {
	invoiceId := r.Get("invoiceId").String()
	if len(invoiceId) == 0 {
		r.Response.Writeln("InvoiceId not found")
		return
	}
	st := r.Get("st").String()
	download := r.Get("download").Bool()
	if !VerifyInvoiceLinkSecurityToken(r.Context(), invoiceId, st) {
		r.Response.Writeln("Invalid link")
		return
	}
	one := query.GetInvoiceByInvoiceId(r.Context(), invoiceId)
	if one == nil {
		r.Response.Writeln("Invoice not found")
		return
	}
	var pdfFileName string
	//if len(one.SendPdf) > 0 {
	//	pdfFileName = utility.DownloadFile(one.SendPdf)
	//} else {
	pdfFileName = handler.GenerateInvoicePdf(r.Context(), one)
	//}
	if len(pdfFileName) == 0 {
		g.Log().Errorf(r.Context(), "LinkEntry pdfFile download or generate error")
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.Writeln("Bad request")
		return
	}
	//r.Response.Header().Add("Access-Control-Allow-Origin", "*")
	//r.Response.Header().Add("Content-Security-Policy", "block-all-mixed-content")
	//r.Response.Header().Add("Strict-Transport-Security", "includeSubDomains")
	//r.Response.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	//r.Response.Header().Add("Access-Control-Allow-Credentials", "true")
	//r.Response.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//r.Response.CORSDefault()
	//corsOptions := r.Response.DefaultCORSOptions()
	//corsOptions.AllowDomain = []string{"user.unibee.top", "merchant.unibee.top"}
	//r.Response.CORS(corsOptions)
	//merchant := query.GetMerchantById(r.Context(), one.MerchantId)
	//r.Response.Header().Set("X-Frame-Options", fmt.Sprintf("ALLOW-FROM %s", merchant.Host))
	if download {
		r.Response.Header().Add("Content-type", "application/octet-stream")
		r.Response.Header().Add("content-disposition", "attachment; filename=\""+pdfFileName+"\"")
	} else {
		r.Response.Header().Add("Content-type", "application/pdf")
	}
	file, err := os.Open(pdfFileName)
	if err != nil {
		g.Log().Errorf(r.Context(), "LinkEntry error:%s", err.Error())
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.Writeln("Bad request")
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			g.Log().Errorf(r.Context(), "LinkEntry error:%s", err.Error())
		}
	}(file)

	_, err = io.Copy(r.Response.ResponseWriter, file)
	if err != nil {
		g.Log().Errorf(r.Context(), "LinkEntry error:%s", err.Error())
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.Writeln("Bad request")
	}
}
