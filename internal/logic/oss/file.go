package oss

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
	"os"
	"strconv"
	"strings"
	"unibee/internal/cmd/config"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

type FileUploadInput struct {
	File       *ghttp.UploadFile
	Path       string
	Name       string
	UserId     string
	RandomName bool
}

type FileUploadOutput struct {
	Id   uint
	Name string
	Path string
	Url  string
}

func Upload(ctx context.Context, in FileUploadInput) (*FileUploadOutput, error) {
	var path string
	if len(in.Path) > 0 {
		path = in.Path
	} else {
		path = "cm"
	}

	tempFileName, err := in.File.Save(".", true)
	if err != nil {
		return nil, err
	}

	userId := in.UserId

	var fileName string
	if in.RandomName || len(in.Name) == 0 {
		fileName = strings.ToLower(strconv.FormatInt(gtime.TimestampNano(), 36) + grand.S(6))
		fileName = fileName + gfile.Ext(in.File.Filename)
	} else {
		fileName = in.Name
	}

	return UploadLocalFile(ctx, tempFileName, path, fileName, userId)
}

func UploadLocalFile(ctx context.Context, localFilePath string, uploadPath string, uploadFileName string, uploadUserId string) (*FileUploadOutput, error) {
	data, err := os.ReadFile(localFilePath)
	if err != nil {
		return nil, err
	}
	if data == nil || len(data) == 0 {
		return nil, gerror.Newf("invalid file, size 0 or nil, localFilePath:%s", localFilePath)
	}

	toSave := entity.FileUpload{
		UserId:     uploadUserId,
		Url:        config.GetConfigInstance().Server.DomainPath + "/oss/file/" + uploadFileName,
		FileName:   uploadFileName,
		Tag:        uploadPath,
		Data:       data,
		CreateTime: gtime.Now().Timestamp(),
	}
	result, err := dao.FileUpload.Ctx(ctx).Data(toSave).OmitNil().Insert()
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()

	return &FileUploadOutput{
		Id:   uint(id),
		Name: toSave.FileName,
		Path: toSave.Tag,
		Url:  toSave.Url,
	}, nil
}
