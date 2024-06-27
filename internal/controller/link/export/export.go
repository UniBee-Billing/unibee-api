package export

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"io"
	"net/http"
	"os"
	"unibee/internal/query"
	"unibee/utility"
)

func LinkExportEntry(r *ghttp.Request) {
	taskId := r.Get("taskId").Int64()
	if taskId <= 0 {
		r.Response.Writeln("TaskId not found")
		return
	}

	one := query.GetMerchantBatchTask(r.Context(), uint64(taskId))
	if one == nil {
		r.Response.Writeln("Task not found")
		return
	}
	if len(one.DownloadUrl) == 0 || one.Status != 2 {
		g.Log().Errorf(r.Context(), "LinkEntry task not success")
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.Writeln("Bad request")
		return
	}

	fileName := utility.DownloadFile(one.DownloadUrl)
	if len(fileName) == 0 {
		g.Log().Errorf(r.Context(), "LinkEntry pdfFile download or generate error")
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.Writeln("Bad request")
		return
	}

	r.Response.Header().Add("Content-type", "application/octet-stream")
	r.Response.Header().Add("content-disposition", "attachment; filename=\""+fileName+"\"")
	file, err := os.Open(fileName)
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
