package invoice

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"io"
	"net/http"
	"os"
	"unibee/internal/logic/invoice/service"
)

func LinkEntry(r *ghttp.Request) {
	invoiceId := r.Get("invoiceId").String()
	res := service.LinkCheck(r.Context(), invoiceId, gtime.Now().Timestamp())
	if len(res.Link) > 0 {
		r.Response.RedirectTo(res.Link)
	} else if len(res.FileName) > 0 {
		r.Response.Header().Add("Content-type", "application/octet-stream")
		r.Response.Header().Add("content-disposition", "attachment; filename=\""+res.FileName+"\"")
		file, err := os.Open(res.FileName)
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
	} else if len(res.Message) > 0 {
		r.Response.Writeln(res.Message)
	} else {
		r.Response.Writeln("Server Error")
	}
}
