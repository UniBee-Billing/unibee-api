package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
	"unibee-api/internal/query"
	_test "unibee-api/test"
)

func TestGetSubscription(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		one := query.GetSubscriptionBySubscriptionId(context.Background(), "sub202402045lnIGlOvznJmSSI")
		_test.AssertNotNil(one)
		t.Assert(one.SubscriptionId, "sub202402045lnIGlOvznJmSSI")
	})
}

//func TestCancelInvoice(t *testing.T) {
//	gtest.C(t, func(t *gtest.T) {
//		err := CancelProcessingInvoice(context.Background(), "ddddd")
//		_test.AssertNotNil(err)
//	})
//}

func setup() {
	fmt.Println("Before all tests")
}

func teardown() {
	fmt.Println("After all tests")
}
