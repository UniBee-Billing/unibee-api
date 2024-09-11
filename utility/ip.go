package utility

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

func DetectLocalIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			for _, addr := range addrs {
				ip, _, err := net.ParseCIDR(addr.String())
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}

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
