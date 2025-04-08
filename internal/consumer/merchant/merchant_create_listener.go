package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	"strconv"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/logic/member"
	"unibee/internal/logic/merchant"
	"unibee/internal/logic/middleware"
	"unibee/internal/logic/vat_gateway/setup"
	"unibee/internal/query"
	"unibee/utility"
)

type MerchantCreateListener struct {
}

func (t MerchantCreateListener) GetTopic() string {
	return redismq2.TopicMerchantCreatedWebhook.Topic
}

func (t MerchantCreateListener) GetTag() string {
	return redismq2.TopicMerchantCreatedWebhook.Tag
}

func (t MerchantCreateListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Infof(ctx, "MerchantCreateListener Receive Message:%s", utility.MarshalToJsonString(message))
	if len(message.Body) > 0 {
		merchantId, _ := strconv.ParseUint(message.Body, 10, 64)
		err := setup.InitMerchantDefaultVatGateway(ctx, merchantId)
		if err != nil {
			g.Log().Errorf(ctx, "MerchantCreateListener InitMerchantDefaultVatGateway err:%s", err.Error())
			return redismq.ReconsumeLater
		}
		merchant.ReloadAllMerchantsCacheForSDKAuthBackground()
		owner := query.GetMerchantOwnerMember(ctx, merchantId)
		if owner != nil {
			member.ReloadMemberCacheForSdkAuthBackground(owner.Id)
		}
		err = merchant.SetupForCloudMode(ctx, merchantId)
		if err != nil {
			g.Log().Errorf(ctx, "MerchantCreateListener SetupForCloudMode err:%s", err.Error())
			return redismq.ReconsumeLater
		}
		middleware.GetMerchantLicense(ctx, merchantId)
	}

	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewUserAccountCreateListener())
	fmt.Println("NewMerchantCreateListener RegisterListener")
}

func NewUserAccountCreateListener() *MerchantCreateListener {
	return &MerchantCreateListener{}
}
