package utility

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(url string) string {
	// 获取图片文件名
	fileName := filepath.Base(url)

	currentDir, err := os.Getwd()
	// 构建本地文件路径
	localFilePath := filepath.Join(currentDir, fileName)
	// 创建文件
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return ""
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 发起 HTTP 请求获取图片
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading image:", err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	// 检查 HTTP 状态码
	if response.StatusCode != http.StatusOK {
		fmt.Println("Error:", response.Status)
		return ""
	}

	// 将内容写入文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return ""
	}
	fmt.Println("Image downloaded successfully:", localFilePath)
	return localFilePath
}
