package utility

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

func DetectLocalIP() string {
	// 获取本机的所有网络接口信息
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	// 遍历所有网络接口
	for _, iface := range interfaces {
		// 排除 lo（loopback）接口
		if iface.Flags&net.FlagLoopback == 0 {
			// 获取该网络接口的所有地址信息
			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			// 遍历所有地址
			for _, addr := range addrs {
				// 将地址转换为 IP 地址
				ip, _, err := net.ParseCIDR(addr.String())
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}

				// 判断是否为 IPv4 地址
				if ip.To4() != nil {
					fmt.Printf("Detect IPv4 Address: %s\n", ip)
					return ip.String()
				}
			}
		}
	}
	return ""
}

var publicIP = ""

func GetPublicIP() string {
	if len(publicIP) > 0 {
		return publicIP
	}
	url := "https://api.ipify.org" // or "https://api.ipify.org?format=text"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("GetPublicIP Error:%s", err.Error())
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("GetPublicIP Error:%s", err.Error())
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("GetPublicIP Error:%s", err.Error())
		return ""
	}
	publicIP = string(body)
	return publicIP
}
