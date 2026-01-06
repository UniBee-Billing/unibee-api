package utility

import (
	"fmt"
	"regexp"
	"strings"
)

func ExtractBrowserOS(userAgent string) string {
	browser := "Unknown Browser"
	os := "Unknown OS"

	// macOS
	if strings.Contains(userAgent, "Mac OS X") {
		os = "Mac"
	} else if strings.Contains(userAgent, "Windows NT") {
		// Windows
		if strings.Contains(userAgent, "Windows NT 10.0") {
			os = "Windows 10"
		} else if strings.Contains(userAgent, "Windows NT 6.3") {
			os = "Windows 8.1"
		} else if strings.Contains(userAgent, "Windows NT 6.2") {
			os = "Windows 8"
		} else if strings.Contains(userAgent, "Windows NT 6.1") {
			os = "Windows 7"
		} else if strings.Contains(userAgent, "Windows NT 6.0") {
			os = "Windows Vista"
		} else {
			os = "Windows" // Fallback for older/unspecified Windows versions
		}
	} else if strings.Contains(userAgent, "Android") {
		os = "Android"
	} else if strings.Contains(userAgent, "iOS") || strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad") {
		os = "iOS"
	} else if strings.Contains(userAgent, "Linux") {
		os = "Linux"
	} else if strings.Contains(userAgent, "CrOS") { // Chrome OS
		os = "Chrome OS"
	}

	// Chrome
	if strings.Contains(userAgent, "Chrome") && !strings.Contains(userAgent, "Edg") { // Edge based on Chromium also contains "Chrome"
		re := regexp.MustCompile(`Chrome/([0-9.]+)`)
		matches := re.FindStringSubmatch(userAgent)
		if len(matches) > 1 {
			browser = "Chrome " + matches[1]
		} else {
			browser = "Chrome"
		}
	} else if strings.Contains(userAgent, "Edg/") || strings.Contains(userAgent, "Edge/") { // Microsoft Edge (Chromium-based or Legacy)
		re := regexp.MustCompile(`(Edg|Edge)/([0-9.]+)`)
		matches := re.FindStringSubmatch(userAgent)
		if len(matches) > 2 {
			browser = "Edge " + matches[2]
		} else {
			browser = "Edge"
		}
	} else if strings.Contains(userAgent, "Firefox") {
		re := regexp.MustCompile(`Firefox/([0-9.]+)`)
		matches := re.FindStringSubmatch(userAgent)
		if len(matches) > 1 {
			browser = "Firefox " + matches[1]
		} else {
			browser = "Firefox"
		}
	} else if strings.Contains(userAgent, "Safari") && !strings.Contains(userAgent, "Chrome") { // Safari (ensure it's not Chrome)
		re := regexp.MustCompile(`Version/([0-9.]+) Safari/`)
		matches := re.FindStringSubmatch(userAgent)
		if len(matches) > 1 {
			browser = "Safari " + matches[1]
		} else {
			browser = "Safari"
		}
	} else if strings.Contains(userAgent, "MSIE") || strings.Contains(userAgent, "Trident/") { // Internet Explorer
		browser = "IE"
	} else if strings.Contains(userAgent, "Opera") || strings.Contains(userAgent, "OPR/") { // Opera
		re := regexp.MustCompile(`(Opera|OPR)/([0-9.]+)`)
		matches := re.FindStringSubmatch(userAgent)
		if len(matches) > 2 {
			browser = "Opera " + matches[2]
		} else {
			browser = "Opera"
		}
	}

	if browser != "Unknown Browser" && os != "Unknown OS" {
		return fmt.Sprintf("%s on %s", strings.TrimSpace(strings.Split(browser, " ")[0]), os)
	} else if browser != "Unknown Browser" {
		return strings.TrimSpace(strings.Split(browser, " ")[0])
	} else if os != "Unknown OS" {
		return os
	}
	return "Unknown Device" // Fallback if nothing specific is found
}
