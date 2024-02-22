package middleware

import (
	"fmt"
	"net/url"
	"testing"
)

func Test_Domain(t *testing.T) {
	origin := "https://www.example.com/path/to/resource"

	parsedURL, err := url.Parse(origin)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// Extract the host (domain) from the parsed URL
	domain := parsedURL.Hostname()

	fmt.Println("Domain:", domain)
}
