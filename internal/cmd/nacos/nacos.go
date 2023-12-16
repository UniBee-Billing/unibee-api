package nacos

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"os"
	"strconv"
	"strings"
)

func Init() {
	// 使用 flag 包声明命令行参数
	var (
		nacosDisableArg   string
		nacosIpArg        string
		nacosPortArg      string
		nacosNamespaceArg string
		nacosGroupArg     string
		nacosDataIdArg    string
	)
	flag.StringVar(&nacosDisableArg, "nacos-disable", "false", "true|false")
	flag.StringVar(&nacosIpArg, "nacos-ip", os.Getenv("nacos.ip"), "ip or domain")
	flag.StringVar(&nacosPortArg, "nacos-port", os.Getenv("nacos.port"), "port like 8848")
	flag.StringVar(&nacosNamespaceArg, "nacos-namespace", os.Getenv("nacos.namespace"), "port like 8848")
	flag.StringVar(&nacosGroupArg, "nacos-group", os.Getenv("nacos.group"), "nacos group")
	flag.StringVar(&nacosDataIdArg, "nacos-data-id", os.Getenv("nacos.data.id"), "nacos dataid like hk-go-settings.yaml")

	// 解析命令行参数
	flag.Parse()

	_ = deleteFile(nacosConfigSyncFilePath) //删除原 nacos 文件
	if g.IsEmpty(nacosDisableArg) || !strings.EqualFold(nacosDisableArg, "true") {
		uPort, err := strconv.ParseUint(nacosPortArg, 10, 64)
		if err != nil {
			fmt.Println("port 转换失败:", err)
			panic(err)
		}
		fmt.Printf("nacos ip:%s \n", nacosIpArg)
		fmt.Printf("nacos port:%d \n", uPort)
		fmt.Printf("nacos namespace:%s \n", nacosNamespaceArg)
		fmt.Printf("nacos group:%s \n", nacosGroupArg)
		fmt.Printf("nacos dataId:%s \n", nacosDataIdArg)

		//获取nacos配置文件
		nacosObj, _ := loadNacosConfig(strings.Trim(nacosIpArg, " "), uPort, strings.Trim(nacosNamespaceArg, " "), strings.Trim(nacosDataIdArg, " "), strings.Trim(nacosGroupArg, " "))
		fmt.Println("nacos config filepath:", nacosObj.GetConfigFilePath())
	} else {
		fmt.Println("nacos disabled")
	}
}

const nacosConfigSyncFilePath = "./nacos.config.yaml"

type Nacos struct {
	ip                                       string
	namespace, dataId, group, configFilePath string
	port                                     uint64
}

// 初始化
func loadNacosConfig(ip string, port uint64, namespace, dataId, group string) (n *Nacos, err error) {

	n = &Nacos{
		ip:        ip,
		port:      port,
		namespace: namespace,
		dataId:    dataId,
		group:     group,
	}
	err = n.syncToFile()
	return
}

func (n Nacos) GetConfigFilePath() string {
	if len(n.configFilePath) == 0 {
		panic("nacos config to save local file is not found!")
	}
	return n.configFilePath
}

func (n *Nacos) syncToFile() (err error) {
	config, err := GetNacosConfig(n.ip, n.port, n.namespace, n.group, n.dataId)
	if err != nil {
		fmt.Println("nacos config load failure")
		panic(err)
	}
	//创建文件
	file, err := createFile(nacosConfigSyncFilePath)
	//关闭
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic("file close error")
		}
	}(file)

	if file == nil {
		panic("create or read file error")
	}
	//写入
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(config)
	err = writer.Flush()
	//返回
	n.configFilePath = nacosConfigSyncFilePath

	return
}

// 创建文件
func createFile(path string) (file *os.File, err error) {
	file, err = os.Create(path)
	if err != nil {
		panic("create file " + err.Error())
	}
	return
}

func deleteFile(path string) (err error) {
	err = os.Remove(path)
	return
}
