package test

import (
	"fmt"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	_ "github.com/gogf/gf/v2/test/gtest"
	"unibee-api/utility"
)

func init() {
	err := g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetPath("/test")
	if err != nil {
		return
	}
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName("test_config.yaml")
}

func AssertNotNil(value interface{}) {
	if utility.IsNil(value) {
		panic(fmt.Sprintf(`[ASSERT] EXPECT Value != nil`))
	}
}
