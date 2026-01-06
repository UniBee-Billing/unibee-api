package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"unibee/internal/cmd/config"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/member"
	"unibee/internal/logic/merchant"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/vat_gateway/setup"
	"unibee/internal/query"
	"unibee/utility"
)

type UpdateMerchantInternalReq struct {
	MerchantId          int64  `json:"merchantId" dc:"Id"`
	CountryCode         string `json:"countryCode" dc:"Country Code"`
	CountryName         string `json:"countryName" dc:"Country Name"`
	CompanyVatNumber    string `json:"companyVatNumber" dc:"Company VAT Number"`
	CompanyRegistryCode string `json:"companyRegistryCode" dc:"Company Registry Code"`
}

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
		one := query.GetMerchantMemberByEmail(ctx, fmt.Sprintf("%s", request))
		if one == nil {
			return nil, gerror.New("not found")
		}
		return one, nil
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
	redismq.RegisterInvoke("UpdateMerchantById", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		g.Log().Infof(ctx, "UpdateMerchant:%s", request)
		if len(fmt.Sprintf("%s", request)) == 0 {
			return nil, gerror.New("invalid request")
		}
		var updateMerchantReq *UpdateMerchantInternalReq
		err = utility.UnmarshalFromJsonString(fmt.Sprintf("%s", request), &updateMerchantReq)
		if err != nil {
			return nil, err
		}
		if updateMerchantReq != nil && updateMerchantReq.MerchantId > 0 {
			one := query.GetMerchantById(ctx, uint64(updateMerchantReq.MerchantId))
			if one == nil {
				return nil, gerror.New("merchant not found")
			}
			_, err = dao.Merchant.Ctx(ctx).Data(g.Map{
				dao.Merchant.Columns().CountryCode: updateMerchantReq.CountryCode,
				dao.Merchant.Columns().CountryName: updateMerchantReq.CountryName,
				dao.Merchant.Columns().BusinessNum: updateMerchantReq.CompanyVatNumber,
				dao.Merchant.Columns().Idcard:      updateMerchantReq.CompanyRegistryCode,
				dao.Merchant.Columns().GmtModify:   gtime.Now(),
			}).Where(dao.Merchant.Columns().Id, updateMerchantReq.MerchantId).OmitEmpty().Update()
			if err != nil {
				return nil, err
			}
			operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
				MerchantId:     uint64(updateMerchantReq.MerchantId),
				Target:         fmt.Sprintf("UpdateMerchantByInternalService"),
				Content:        "Update",
				UserId:         0,
				SubscriptionId: "",
				InvoiceId:      "",
				PlanId:         0,
				DiscountCode:   "",
			}, err)
			return nil, err
		} else {
			return nil, gerror.New("UnmarshalFromJsonString request error")
		}
	})
	redismq.RegisterInvoke("DeleteMemberByEmail", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if config.GetConfigInstance().IsProd() || config.GetConfigInstance().Mode != "cloud" {
			return nil, gerror.New("not support env")
		}
		g.Log().Infof(ctx, "DeleteMemberByEmail:%s", request)
		targetMember := query.GetMerchantMemberByEmail(ctx, fmt.Sprintf("%s", request))
		if targetMember != nil && targetMember.Role != "Owner" {
			_, err = dao.MerchantMember.Ctx(ctx).Where(dao.MerchantMember.Columns().Id, targetMember.Id).Delete()
			if err == nil {
				operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
					MerchantId:     targetMember.MerchantId,
					Target:         fmt.Sprintf("Member(%v)", targetMember.Email),
					Content:        "DeleteViaProdMerchantCreation",
					UserId:         0,
					SubscriptionId: "",
					InvoiceId:      "",
					PlanId:         0,
					DiscountCode:   "",
				}, err)
			}
			return nil, err
		} else {
			return nil, gerror.New("member not found or is owner")
		}
	})
	redismq.RegisterInvoke("NewMemberSessionByEmail", func(ctx context.Context, request interface{}) (response interface{}, err error) {
		g.Log().Infof(ctx, "NewMemberSessionByEmail:%s", request)
		targetMember := query.GetMerchantMemberByEmail(ctx, fmt.Sprintf("%s", request))
		if targetMember != nil {
			session, err := member.NewMemberSession(ctx, int64(targetMember.Id), "")
			if err != nil {
				return nil, err
			}
			return session, nil
		} else {
			return nil, gerror.New("member not found")
		}
	})
}
