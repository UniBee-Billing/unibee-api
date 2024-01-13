package vatsense

import (
	"fmt"
	"testing"
)

var apiKey = "d9d57f2212cba7e286b3fb9cbb2ad419"

func Test(t *testing.T) {
	fmt.Println(ListAllCountries(apiKey))
}
