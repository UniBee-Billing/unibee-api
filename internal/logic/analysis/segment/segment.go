package segment

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/segmentio/analytics-go/v3"
	"strconv"
	"unibee/internal/consumer/webhook/log"
	"unibee/internal/logic/merchant_config"
	entity "unibee/internal/model/entity/oversea_pay"
)

const (
	KeyMerchantSegmentServer     = "KEY_MERCHANT_SEGMENT_SERVER"
	KeyMerchantSegmentUserPortal = "KEY_MERCHANT_SEGMENT_USER_PORTAL"
)

func GetMerchantSegmentUserPortalConfig(ctx context.Context, merchantId uint64) (value string) {
	keyConfig := merchant_config.GetMerchantConfig(ctx, merchantId, KeyMerchantSegmentUserPortal)
	if keyConfig != nil {
		value = keyConfig.ConfigValue
	}
	return
}

func GetMerchantSegmentServerSideConfig(ctx context.Context, merchantId uint64) (value string) {
	keyConfig := merchant_config.GetMerchantConfig(ctx, merchantId, KeyMerchantSegmentServer)
	if keyConfig != nil {
		value = keyConfig.ConfigValue
	}
	return
}

func RegisterSegmentUserBackground(superCtx context.Context, merchantId uint64, user *entity.UserAccount) {
	if merchantId <= 0 || user == nil {
		return
	}
	secret := GetMerchantSegmentServerSideConfig(superCtx, merchantId)
	if len(secret) <= 0 {
		return
	}
	go func() {
		ctx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				log.PrintPanic(ctx, err)
				return
			}
		}()

		client := analytics.New(secret)
		defer func(client analytics.Client) {
			err = client.Close()
			if err != nil {
				g.Log().Errorf(ctx, "RegisterSegmentUserBackground error:%s", err.Error())
			}
		}(client)
		err = client.Enqueue(analytics.Identify{
			UserId: strconv.FormatUint(user.Id, 10),
			Traits: analytics.NewTraits().
				SetName(fmt.Sprintf("%s %s", user.FirstName, user.LastName)).
				SetEmail(user.Email).
				Set("externalUserId", user.ExternalUserId),
		})
		if err != nil {
			g.Log().Errorf(ctx, "RegisterSegmentUserBackground error:%s", err.Error())
			return
		}
	}()
}

func TrackSegmentEventBackground(superCtx context.Context, merchantId uint64, user *entity.UserAccount, event string, data map[string]interface{}) {
	if merchantId <= 0 || user == nil {
		return
	}
	secret := GetMerchantSegmentServerSideConfig(superCtx, merchantId)
	if len(secret) <= 0 {
		return
	}
	go func() {
		ctx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				log.PrintPanic(ctx, err)
				return
			}
		}()

		client := analytics.New(secret)
		defer func(client analytics.Client) {
			err = client.Close()
			if err != nil {
				g.Log().Errorf(ctx, "TrackSegmentEvent error:%s", err.Error())
			}
		}(client)

		err = client.Enqueue(analytics.Identify{
			UserId: strconv.FormatUint(user.Id, 10),
			Traits: analytics.NewTraits().
				SetName(fmt.Sprintf("%s %s", user.FirstName, user.LastName)).
				SetEmail(user.Email).
				Set("externalUserId", user.ExternalUserId),
		})
		if err != nil {
			g.Log().Errorf(ctx, "TrackSegmentEvent Register error:%s", err.Error())
		}

		properties := analytics.NewProperties()
		for key, value := range data {
			properties.Set(key, value)
		}
		err = client.Enqueue(analytics.Track{
			Event:      event,
			UserId:     strconv.FormatUint(user.Id, 10),
			Properties: properties,
		})

		if err != nil {
			g.Log().Errorf(ctx, "TrackSegmentEvent error:%s", err.Error())
			return
		}
	}()
}
