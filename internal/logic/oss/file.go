package oss

import (
	"bytes"
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unibee/internal/cmd/config"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

// FileUploadInput
type FileUploadInput struct {
	File       *ghttp.UploadFile //
	Path       string            //
	Name       string            //
	UserId     string            // UserId
	RandomName bool              //
}

// FileUploadOutput
type FileUploadOutput struct {
	Id   uint   //
	Name string //
	Path string //
	Url  string //
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

	minioClient, err := minio.New(config.GetConfigInstance().MinioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.GetConfigInstance().MinioConfig.AccessKey, config.GetConfigInstance().MinioConfig.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	_, err = minioClient.PutObject(ctx, config.GetConfigInstance().MinioConfig.BucketName, gfile.Join(uploadPath, uploadFileName), bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: http.DetectContentType(data),
	})
	if err != nil {
		return nil, err
	}

	toSave := entity.FileUpload{
		UserId:     uploadUserId,
		Url:        config.GetConfigInstance().MinioConfig.Domain + "/invoice/" + gfile.Join(uploadPath, uploadFileName),
		FileName:   uploadFileName,
		Tag:        uploadPath,
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
