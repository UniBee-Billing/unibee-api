package query

import (
	"context"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetOssFileByFileName(ctx context.Context, filename string) (one *entity.FileUpload) {
	if len(filename) == 0 {
		return nil
	}
	err := dao.FileUpload.Ctx(ctx).Where(dao.FileUpload.Columns().FileName, filename).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
