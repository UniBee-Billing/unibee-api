package test

import (
	"context"
	"fmt"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	_ "github.com/gogf/gf/v2/test/gtest"
	"unibee/internal/cmd/config"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/test/prepare"
	"unibee/utility"
)

var TestMerchant *entity.Merchant
var TestMerchantMember *entity.MerchantMember
var TestUser *entity.UserAccount
var TestPlan *entity.Plan
var TestAddon *entity.Plan

func init() {
	ctx := context.Background()
	err := g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetPath("/test")
	if err != nil {
		return
	}
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName("test_config.yaml")

	config.SetupDefaultConfigs(ctx)
	TestMerchantMember = query.GetMerchantMemberByEmail(ctx, "test@wowow.io")
	if TestMerchantMember == nil {
		TestMerchant, TestMerchantMember, err = prepare.CreateTestMerchantAccount(ctx)
		utility.AssertError(err, "CreateTestMerchantAccount err")
	}
	utility.Assert(TestMerchantMember.MerchantId > 0, "TestMerchantMember err")
	TestMerchant = query.GetMerchantById(ctx, TestMerchantMember.MerchantId)
	utility.Assert(TestMerchant != nil, "TestMerchant failure")
	utility.Assert(TestMerchantMember != nil, "TestMerchantMember failure")
	TestUser = query.GetUserAccountByEmail(ctx, TestMerchant.Id, "testuser@wowow.io")
	if TestUser == nil {
		TestUser, err = prepare.CreateTestUser(ctx, TestMerchant.Id)
		utility.AssertError(err, "CreateTestUser err")
	}
	utility.Assert(TestUser != nil, "TestUser err")
	utility.Assert(TestUser.MerchantId > 0, "TestUser err")
	TestPlan = prepare.GetPlanByName(ctx, "autotest_x")
	if TestPlan == nil {
		TestPlan, err = prepare.CreateTestPlan(ctx, TestMerchant.Id)
		utility.AssertError(err, "CreateTestPlan err")
	}
	utility.Assert(TestPlan != nil, "TestPlan err")
	utility.Assert(TestPlan.MerchantId > 0, "TestPlan err")
	TestAddon = prepare.GetPlanByName(ctx, "autotest_addon_x")
	if TestAddon == nil {
		TestAddon, err = prepare.CreateTestAddon(ctx, TestMerchant.Id)
		utility.AssertError(err, "CreateTestAddon err")
	}
	utility.Assert(TestAddon != nil, "TestAddon err")
	utility.Assert(TestAddon.MerchantId > 0, "TestAddon err")
}

func AssertNotNil(value interface{}) {
	if utility.IsNil(value) {
		panic(fmt.Sprintf(`[ASSERT] EXPECT Value != nil`))
	}
}
