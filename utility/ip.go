package utility

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
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

func ExtractFirstIPAddresses(text string) string {
	var ips []string

	// Remove \b
	ipv4Regex := regexp.MustCompile(`(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

	matchesV4 := ipv4Regex.FindAllString(text, -1)
	fmt.Println("V4 matches:", matchesV4)

	for _, match := range matchesV4 {
		if net.ParseIP(match) != nil {
			ips = append(ips, match)
		}
	}

	if len(ips) > 0 {
		return ips[0]
	}
	return ""
}
