package utility

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(url string) string {
	fileName := filepath.Base(url)

	currentDir, err := os.Getwd()
	localFilePath := filepath.Join(currentDir, fileName)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file:%s\n", err.Error())
		return ""
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Printf("Error creating file:%s\n", err.Error())
		}
	}(file)

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error downloading file:%s\n", err.Error())
		return ""
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Printf("Error creating file:%s\n", err.Error())
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Error:%s\n", response.Status)
		return ""
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Printf("Error writing to file:%s\n", err.Error())
		return ""
	}
	fmt.Println("File downloaded successfully:", localFilePath)
	return localFilePath
}
