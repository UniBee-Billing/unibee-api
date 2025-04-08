package middleware

import (
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"testing"
)

func TestCheckQPSLimit(t *testing.T) {
	//gredis.SetConfig(&gredis.Config{
	//	Address: "127.0.0.1:6379",
	//	Db:      1,
	//	Pass:    "changeme",
	//})
	//ctx := context.Background()
	//key := "qps:test:user1000"
	//maxQPS := 5
	//expireMs := 1000

	//for i := 0; i < 15; i++ {
	//	if CheckQPSLimit(ctx, key, maxQPS, expireMs) {
	//		fmt.Println("Pass ✅")
	//	} else {
	//		fmt.Println("Reject ❌")
	//	}
	//}
	//time.Sleep(1 * time.Second)
	//for i := 0; i < 15; i++ {
	//	if CheckQPSLimit(ctx, key, maxQPS, expireMs) {
	//		fmt.Println("Pass ✅")
	//	} else {
	//		fmt.Println("Reject ❌")
	//	}
	//}
}
