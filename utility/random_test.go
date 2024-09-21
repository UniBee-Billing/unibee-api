package utility

import (
	"fmt"
	"testing"
)

func TestGenerateRandomAlphanumeric(t *testing.T) {
	fmt.Println(GenerateRandomAlphanumeric(32))
	fmt.Println(CreateInvoiceId())
}
