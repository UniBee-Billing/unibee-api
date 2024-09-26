package oss

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"unibee/internal/query"
)

func FileEntry(r *ghttp.Request) {
	filename := r.Get("filename").String()
	if len(filename) == 0 {
		r.Response.Writeln("filename not found")
		return
	}

	one := query.GetOssFileByFileName(r.Context(), filename)
	if one == nil {
		r.Response.Writeln("filename not found")
		return
	}
	extension := filepath.Ext(filename)
	var download bool
	if extension == ".jpg" || extension == ".jpeg" || extension == ".png" {
		download = false
	} else {
		download = true
	}
	r.Response.Header().Add("Content-Length", fmt.Sprintf("%d", len(one.Data)))
	if download {
		r.Response.Header().Add("Content-type", "application/octet-stream")
		r.Response.Header().Add("content-disposition", "attachment; filename=\""+filename+"\"")
	} else {
		r.Response.Header().Add("Content-type", "image/"+strings.ReplaceAll(extension, ".", ""))
	}
	_, err := io.Copy(r.Response.ResponseWriter, bytes.NewReader(one.Data))
	if err != nil {
		g.Log().Errorf(r.Context(), "LinkEntry error:%s", err.Error())
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.Writeln("Bad request")
	}
}
