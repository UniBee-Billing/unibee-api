package merchant

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) Update(ctx context.Context, req *member.UpdateReq) (res *member.UpdateRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "Merchant Member Not Found")
	one := query.GetMerchantMemberById(ctx, _interface.Context().Get(ctx).MerchantMember.Id)
	utility.Assert(one != nil, "Merchant Member Not Found")
	_, err = dao.MerchantMember.Ctx(ctx).Data(g.Map{
		dao.MerchantMember.Columns().FirstName: req.FirstName,
		dao.MerchantMember.Columns().LastName:  req.LastName,
		dao.MerchantMember.Columns().Mobile:    req.Mobile,
	}).Where(dao.MerchantMember.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	if req.Metadata != nil {
		var metadata = make(map[string]interface{})
		if len(one.MetaData) > 0 {
			_ = gjson.Unmarshal([]byte(one.MetaData), &metadata)
		}
		for k, v := range *req.Metadata {
			metadata[k] = v
		}
		_, _ = dao.MerchantMember.Ctx(ctx).Data(g.Map{
			dao.MerchantMember.Columns().MetaData: utility.MarshalToJsonString(metadata),
		}).Where(dao.MerchantMember.Columns().Id, one.Id).OmitNil().Update()
	}
	return &member.UpdateRes{MerchantMember: detail.ConvertMemberToDetail(ctx, one)}, nil
}
