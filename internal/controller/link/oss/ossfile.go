package oss

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"io"
	"net/http"
	"os"
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
	var exist bool
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) && err == nil {
		exist = true
	}
	if len(one.Data) > 0 && !exist {
		err := os.WriteFile(filename, one.Data, 0644)
		if err != nil {
			g.Log().Errorf(r.Context(), "LinkEntry error:%s", err.Error())
			r.Response.WriteHeader(http.StatusBadRequest)
			r.Response.Writeln("Bad request")
			return
		}
	}
	if download {
		r.Response.Header().Add("Content-type", "application/octet-stream")
		r.Response.Header().Add("content-disposition", "attachment; filename=\""+filename+"\"")
	} else {
		r.Response.Header().Add("Content-type", "image/"+strings.ReplaceAll(extension, ".", ""))
	}
	file, err := os.Open(filename)
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
