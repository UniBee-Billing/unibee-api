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
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"io/ioutil"
	"strconv"
	"strings"
)

// FileUploadInput 上传文件输入参数
type FileUploadInput struct {
	File       *ghttp.UploadFile // 上传文件对象
	Path       string            // 上传目录
	Name       string            // 自定义文件名称
	UserId     string            // UserId
	RandomName bool              // 是否随机命名文件
}

// FileUploadOutput 上传文件返回参数
type FileUploadOutput struct {
	Id   uint   // 数据表ID
	Name string // 文件名称
	Path string // 本地路径
	Url  string // 访问URL，可能只是URI
}

func Upload(ctx context.Context, in FileUploadInput) (*FileUploadOutput, error) {
	minioClient, err := minio.New(consts.GetConfigInstance().MinioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(consts.GetConfigInstance().MinioConfig.AccessKey, consts.GetConfigInstance().MinioConfig.SecretKey, ""),
		Secure: false, // 如果是 HTTPS 连接，请将其设置为 true
	})
	if err != nil {
		return nil, err
	}

	tempFileName, err := in.File.Save(".", true)

	// 读取本地文件
	data, err := ioutil.ReadFile(tempFileName)
	if err != nil {
		return nil, err
	}

	var path string

	if len(in.Path) > 0 {
		path = in.Path
	} else {
		path = "cm"
	}

	var fileName string
	if in.RandomName || len(in.Name) == 0 {
		fileName = strings.ToLower(strconv.FormatInt(gtime.TimestampNano(), 36) + grand.S(6))
		fileName = fileName + gfile.Ext(in.File.Filename)
	} else {
		fileName = in.Name
	}

	_, err = minioClient.PutObject(ctx, consts.GetConfigInstance().MinioConfig.BucketName, gfile.Join(path, fileName), bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}

	// 记录到数据表
	toSave := entity.FileUpload{
		UserId:   in.UserId,
		Url:      consts.GetConfigInstance().MinioConfig.Domain + "/invoice/" + gfile.Join(path, fileName),
		FileName: fileName,
		Tag:      path,
	}
	result, err := dao.FileUpload.Ctx(ctx).Data(toSave).OmitEmpty().Insert()
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
