package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	"unibee/internal/logic/member"
	"unibee/internal/logic/merchant"
	"unibee/internal/logic/vat_gateway/setup"
	"unibee/internal/query"
	"unibee/utility"
)

func init() {
	redismq.RegisterInvoke("GetMerchantByOwnerEmail", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		g.Log().Infof(ctx, "GetMerchantByOwnerEmail:%s", request)
		if request == nil || len(fmt.Sprintf("%s", request)) == 0 {
			return nil, gerror.New("invalid email")
		}
		one := query.GetMerchantByOwnerEmail(ctx, fmt.Sprintf("%s", request))
		if one == nil {
			return nil, gerror.New("not found")
		}
		return one, nil
	})
	redismq.RegisterInvoke("GetMerchantById", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		g.Log().Infof(ctx, "GetMerchantById:%s", request)
		if merchantId, ok := request.(float64); ok {
			one := query.GetMerchantById(ctx, uint64(merchantId))
			if one == nil {
				return nil, gerror.New("not found")
			}
			return one, nil
		} else {
			return nil, gerror.New("invalid request")
		}
	})
	redismq.RegisterInvoke("GetMerchantMemberByEmail", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		g.Log().Infof(ctx, "GetMerchantMemberByEmail:%s", request)
		if request == nil || len(fmt.Sprintf("%s", request)) == 0 {
			return nil, gerror.New("invalid email")
		}
		member := query.GetMerchantMemberByEmail(ctx, fmt.Sprintf("%s", request))
		if member == nil {
			return nil, gerror.New("not found")
		}
		return member, nil
	})
	redismq.RegisterInvoke("GetMerchantOwnerMember", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		g.Log().Infof(ctx, "GetMerchantOwnerMember:%s", request)
		if merchantId, ok := request.(float64); ok {
			one := query.GetMerchantOwnerMember(ctx, uint64(merchantId))
			if one == nil {
				return nil, gerror.New("not found")
			}
			return one, nil
		} else {
			return nil, gerror.New("invalid request")
		}
	})
	redismq.RegisterInvoke("InitMerchantDefaultVatGateway", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		g.Log().Infof(ctx, "InitMerchantDefaultVatGateway:%s", request)
		if merchantId, ok := request.(float64); ok {
			err = setup.InitMerchantDefaultVatGateway(ctx, uint64(merchantId))
			if err != nil {
				return nil, err
			}
			return nil, nil
		} else {
			return nil, gerror.New("invalid request")
		}
	})
	redismq.RegisterInvoke("QueryOrCreateMerchant", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		g.Log().Infof(ctx, "QueryOrCreateMerchant:%s", request)
		if len(fmt.Sprintf("%s", request)) == 0 {
			return nil, gerror.New("invalid request")
		}
		var createMerchantReq *merchant.CreateMerchantInternalReq
		err = utility.UnmarshalFromJsonString(fmt.Sprintf("%s", request), &createMerchantReq)
		if err != nil {
			return nil, err
		}
		if createMerchantReq != nil {
			mer, targetMember, err := merchant.QueryOrCreateMerchant(ctx, createMerchantReq)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{"merchant": mer, "member": targetMember}, err
		} else {
			return nil, gerror.New("UnmarshalFromJsonString request error")
		}
	})
	redismq.RegisterInvoke("NewMemberSessionByEmail", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		g.Log().Infof(ctx, "NewMemberSessionByEmail:%s", request)
		targetMember := query.GetMerchantMemberByEmail(ctx, fmt.Sprintf("%s", request))
		if targetMember != nil {
			session, err := member.NewSession(ctx, int64(targetMember.Id), "")
			if err != nil {
				return nil, err
			}
			return session, nil
		} else {
			return nil, gerror.New("member not found")
		}
	})
}
