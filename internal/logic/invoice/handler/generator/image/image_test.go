package image

import (
	"bytes"
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"testing"
)

func TestImageEncode(t *testing.T) {
	logoBytes, err := os.ReadFile("test.jpg")
	if err != nil {
		g.Log().Errorf(context.Background(), "read file test.jpg error:%s", err.Error())
	}
	// Get image format
	_, format, _ := image.DecodeConfig(bytes.NewReader(logoBytes))
	g.Log().Infof(context.Background(), "TestImageEncode:%s", format)
}
