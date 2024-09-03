package utility

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"os"
)

func ReadBuildVersionInfo(ctx context.Context) string {
	buildInfo, err := os.ReadFile("./version.txt")
	if err != nil {
		glog.Errorf(ctx, "ReadBuildVersionInfo error:%s", err.Error())
	}
	return string(buildInfo)
}
