package config

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go-oversea-pay/internal/consts"
	"gopkg.in/yaml.v3"
)

func GetNacosConfig(ip string, port uint64, namespace string, group string, dataId string) (string, error) {
	sc := []constant.ServerConfig{{
		IpAddr: ip,
		Port:   port,
	}}

	cc := constant.ClientConfig{
		NamespaceId:         namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "log",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})

	if err != nil {
		return "", err
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})

	if err != nil {
		return "", err
	}

	err = yaml.Unmarshal([]byte(content), consts.GetConfigInstance())
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
	} else {
		fmt.Printf(`Nacos Sync Config: %+v`, consts.GetConfigInstance())
		fmt.Println("")
	}

	return content, err
}
