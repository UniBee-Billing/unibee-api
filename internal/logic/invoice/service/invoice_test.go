package service

import (
	"context"
	"testing"
)

func TestCancelInvoice(t *testing.T) {
	err := CancelProcessingInvoice(context.Background(), "ddddd")
	if err != nil {
		return
	}
}
